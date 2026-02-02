package ingest

import (
	"strings"

	"golang.org/x/net/html"
)

// ParsedContent — распарсенный контент страницы.
type ParsedContent struct {
	Title       string
	Paragraphs  []string
	CodeBlocks  []CodeBlock
	Lists       []string
	RawHTML     string
}

// CodeBlock — блок кода из страницы.
type CodeBlock struct {
	Code     string
	Language string
	Caption  string
}

// Parser парсит HTML страницы урока.
type Parser struct{}

// NewParser создаёт новый парсер.
func NewParser() *Parser {
	return &Parser{}
}

// Parse парсит HTML и извлекает структурированный контент.
func (p *Parser) Parse(htmlContent string) (*ParsedContent, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	content := &ParsedContent{
		RawHTML: htmlContent,
	}

	// Находим основной контент (обычно в article, main, или div с определённым классом)
	mainContent := p.findMainContent(doc)
	if mainContent == nil {
		mainContent = doc
	}

	// Извлекаем заголовок
	content.Title = p.extractTitle(doc)

	// Извлекаем параграфы и блоки кода
	p.extractContent(mainContent, content)

	return content, nil
}

// findMainContent ищет основной контент страницы.
func (p *Parser) findMainContent(doc *html.Node) *html.Node {
	var find func(*html.Node) *html.Node
	find = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode {
			// Ищем article, main или div с классом content/article/main
			if n.Data == "article" || n.Data == "main" {
				return n
			}

			// Проверяем классы
			class := getAttr(n, "class")
			if n.Data == "div" && (strings.Contains(class, "content") ||
				strings.Contains(class, "article") ||
				strings.Contains(class, "main") ||
				strings.Contains(class, "center") ||
				strings.Contains(class, "tutorial")) {
				return n
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if found := find(child); found != nil {
				return found
			}
		}
		return nil
	}

	return find(doc)
}

// extractTitle извлекает заголовок страницы.
func (p *Parser) extractTitle(doc *html.Node) string {
	var find func(*html.Node) string
	find = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "h1" {
			return getTextContent(n)
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if title := find(child); title != "" {
				return title
			}
		}
		return ""
	}

	title := find(doc)
	if title == "" {
		// Пробуем найти в <title>
		var findTitle func(*html.Node) string
		findTitle = func(n *html.Node) string {
			if n.Type == html.ElementNode && n.Data == "title" {
				return getTextContent(n)
			}
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				if t := findTitle(child); t != "" {
					return t
				}
			}
			return ""
		}
		title = findTitle(doc)
		// Убираем суффикс сайта
		if idx := strings.Index(title, "|"); idx > 0 {
			title = strings.TrimSpace(title[:idx])
		}
		if idx := strings.Index(title, "-"); idx > 0 {
			title = strings.TrimSpace(title[:idx])
		}
	}

	return title
}

// extractContent извлекает параграфы, списки и блоки кода.
func (p *Parser) extractContent(n *html.Node, content *ParsedContent) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "p":
			text := getTextContent(n)
			if text != "" && len(text) > 10 {
				// Проверяем, не рекламный ли это блок
				if !isAdvertisement(text) {
					content.Paragraphs = append(content.Paragraphs, text)
				}
			}

		case "pre", "code":
			// Извлекаем блок кода
			code := getTextContent(n)
			if code != "" && len(code) > 5 {
				lang := detectLanguage(code, getAttr(n, "class"))
				content.CodeBlocks = append(content.CodeBlocks, CodeBlock{
					Code:     code,
					Language: lang,
				})
				return // Не углубляемся в code
			}

		case "ul", "ol":
			// Извлекаем список
			list := p.extractList(n)
			if list != "" {
				content.Lists = append(content.Lists, list)
			}
			return // Не углубляемся в списки
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		p.extractContent(child, content)
	}
}

// extractList извлекает текст списка.
func (p *Parser) extractList(n *html.Node) string {
	var items []string
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "li" {
			text := getTextContent(child)
			if text != "" {
				items = append(items, "- "+text)
			}
		}
	}
	return strings.Join(items, "\n")
}

// isAdvertisement проверяет, является ли текст рекламой.
func isAdvertisement(text string) bool {
	lower := strings.ToLower(text)
	adKeywords := []string{
		"реклама", "advertisement", "sponsor",
		"яндекс", "google ads", "click here",
		"партнёр", "partner", "cookies",
	}
	for _, kw := range adKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// detectLanguage определяет язык блока кода.
func detectLanguage(code, class string) string {
	// По классу
	if strings.Contains(class, "go") || strings.Contains(class, "golang") {
		return "go"
	}
	if strings.Contains(class, "bash") || strings.Contains(class, "shell") {
		return "bash"
	}

	// По содержимому
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "package ") ||
		strings.Contains(code, "func ") ||
		strings.Contains(code, "import (") ||
		strings.Contains(code, "fmt.") {
		return "go"
	}

	if strings.HasPrefix(code, "$") || strings.HasPrefix(code, "go ") {
		return "bash"
	}

	return "go" // По умолчанию Go
}
