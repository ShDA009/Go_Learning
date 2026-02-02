package ingest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// TOCEntry — запись в оглавлении (урок или раздел).
type TOCEntry struct {
	Title      string
	URL        string
	IsModule   bool
	ModuleSlug string
	OrderIndex int
}

// Crawler скачивает страницы с сайта.
type Crawler struct {
	client  *http.Client
	baseURL string
}

// NewCrawler создаёт новый crawler.
func NewCrawler(baseURL string) *Crawler {
	return &Crawler{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

// FetchTOC загружает оглавление курса.
func (c *Crawler) FetchTOC(ctx context.Context) ([]TOCEntry, error) {
	url := c.baseURL + "/"

	body, err := c.fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetch TOC: %w", err)
	}

	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parse TOC HTML: %w", err)
	}

	return c.parseTOC(doc), nil
}

// FetchPage загружает страницу урока.
func (c *Crawler) FetchPage(ctx context.Context, path string) (string, error) {
	url := path
	if !strings.HasPrefix(path, "http") {
		url = c.baseURL + "/" + strings.TrimPrefix(path, "/")
	}

	return c.fetch(ctx, url)
}

func (c *Crawler) fetch(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "GoLearning/1.0 (educational crawler)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// parseTOC извлекает ссылки на уроки из HTML оглавления.
func (c *Crawler) parseTOC(doc *html.Node) []TOCEntry {
	var entries []TOCEntry
	var currentModule string
	orderIndex := 0

	// Ищем div.nav или подобный контейнер с навигацией
	var findNav func(*html.Node) *html.Node
	findNav = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode {
			// Ищем nav или div с классом nav/menu/sidebar
			if n.Data == "nav" {
				return n
			}
			for _, attr := range n.Attr {
				if attr.Key == "class" && (strings.Contains(attr.Val, "nav") ||
					strings.Contains(attr.Val, "menu") ||
					strings.Contains(attr.Val, "sidebar")) {
					return n
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if found := findNav(child); found != nil {
				return found
			}
		}
		return nil
	}

	// Извлекаем все ссылки
	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Проверяем заголовки разделов (h2, h3, strong, b)
			if n.Data == "h2" || n.Data == "h3" || n.Data == "strong" || n.Data == "b" {
				text := getTextContent(n)
				if text != "" && !strings.Contains(strings.ToLower(text), "metanit") {
					currentModule = text
				}
			}

			// Ищем ссылки на уроки
			if n.Data == "a" {
				href := getAttr(n, "href")
				title := strings.TrimSpace(getTextContent(n))

				// Фильтруем: только ссылки на уроки Go (например, /go/tutorial/1.1.php)
				if strings.Contains(href, "/go/tutorial/") && strings.HasSuffix(href, ".php") {
					// Исключаем навигационные ссылки
					if title != "" && !strings.Contains(strings.ToLower(title), "metanit") &&
						!strings.Contains(strings.ToLower(title), "предыдущ") &&
						!strings.Contains(strings.ToLower(title), "следующ") {
						orderIndex++
						entries = append(entries, TOCEntry{
							Title:      title,
							URL:        href,
							ModuleSlug: slugify(currentModule),
							OrderIndex: orderIndex,
						})
					}
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			extractLinks(child)
		}
	}

	// Сначала пробуем найти навигацию
	nav := findNav(doc)
	if nav != nil {
		extractLinks(nav)
	}

	// Если ничего не нашли, ищем по всему документу
	if len(entries) == 0 {
		extractLinks(doc)
	}

	return entries
}

// getTextContent возвращает текстовое содержимое узла.
func getTextContent(n *html.Node) string {
	var sb strings.Builder
	var extract func(*html.Node)
	extract = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			extract(child)
		}
	}
	extract(n)
	return strings.TrimSpace(sb.String())
}

// getAttr возвращает значение атрибута.
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
