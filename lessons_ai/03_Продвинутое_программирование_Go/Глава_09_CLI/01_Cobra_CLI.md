# Cobra CLI Framework

---

## 💡 Ключевые идеи

1. **Cobra** — самый популярный фреймворк для CLI в Go
2. **Commands** — подкоманды (git commit, docker run)
3. **Flags** — параметры (--verbose, -v)
4. **Arguments** — позиционные аргументы
5. **Autocomplete** — автодополнение для bash/zsh

---

## 📖 Теория

### Что такое CLI?

**CLI (Command Line Interface)** — интерфейс командной строки. Это программы, которые вы запускаете в терминале:

```bash
git commit -m "message"
docker run -p 8080:80 nginx
kubectl get pods --namespace default
```

### Зачем писать CLI на Go?

1. **Один бинарник** — не нужно устанавливать runtime (как Python, Node.js)
2. **Кросс-платформенность** — компилируется под Linux, macOS, Windows
3. **Скорость** — запускается мгновенно
4. **Экосистема** — Cobra + Viper = мощные CLI

Примеры известных CLI на Go: Docker, Kubernetes (kubectl), Terraform, Hugo.

### Что такое Cobra?

**Cobra** — это библиотека для создания CLI-приложений. Её используют в Docker, Kubernetes, GitHub CLI, и многих других проектах.

Cobra предоставляет:
- Команды и подкоманды (`app serve`, `app db migrate`)
- Флаги (`--verbose`, `-v`, `--config=config.yaml`)
- Аргументы (`app add user john`)
- Автодополнение (tab completion)
- Генерация документации

### Анатомия CLI-команды

```bash
myapp server start --port 8080 --verbose
│      │      │     │          │
│      │      │     │          └── флаг (boolean)
│      │      │     └── флаг со значением
│      │      └── подкоманда (subcommand)
│      └── команда
└── приложение
```

### Cobra: структура проекта

```
myapp/
├── cmd/
│   ├── root.go      # корневая команда
│   ├── serve.go     # myapp serve
│   └── db/
│       ├── db.go    # myapp db
│       └── migrate.go # myapp db migrate
├── main.go
└── go.mod
```

### Команды в Cobra

```go
var serveCmd = &cobra.Command{
    Use:   "serve",                        // как вызывать
    Short: "Start the server",             // короткое описание
    Long:  `Long description...`,          // подробное описание
    Run: func(cmd *cobra.Command, args []string) {
        // код выполнения
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)  // регистрация
}
```

### Флаги (Flags)

**Persistent Flags** — доступны для команды и всех подкоманд:
```go
rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
```

**Local Flags** — только для конкретной команды:
```go
serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "server port")
```

### Viper — конфигурация

**Viper** — библиотека для работы с конфигурацией. Интегрируется с Cobra:

```go
viper.SetConfigName("config")
viper.AddConfigPath(".")
viper.ReadInConfig()

// Приоритет:
// 1. Флаги командной строки (--port=8080)
// 2. Переменные окружения (PORT=8080)
// 3. Файл конфигурации (config.yaml)
// 4. Значения по умолчанию
```

### Аргументы

```go
var addCmd = &cobra.Command{
    Use:   "add [name]",
    Args:  cobra.ExactArgs(1),  // ровно 1 аргумент
    Run: func(cmd *cobra.Command, args []string) {
        name := args[0]
        fmt.Println("Adding:", name)
    },
}
```

Валидаторы аргументов:
- `cobra.NoArgs` — запрещены
- `cobra.ExactArgs(n)` — ровно n
- `cobra.MinimumNArgs(n)` — минимум n
- `cobra.MaximumNArgs(n)` — максимум n

### Best Practices

1. **Логичная структура команд** — группируйте связанные команды
2. **Понятные описания** — пишите Short и Long
3. **Разумные умолчания** — флаги должны иметь sensible defaults
4. **Validation** — проверяйте входные данные
5. **Exit codes** — 0 = успех, не-0 = ошибка

---

## 📋 Установка

```bash
go get -u github.com/spf13/cobra/cobra
go install github.com/spf13/cobra-cli@latest
```

---

## 💻 Примеры кода

### Пример 1: Базовая структура

```go
// cmd/root.go
package cmd

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgFile string
    verbose bool
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "MyApp - example CLI application",
    Long: `MyApp is a CLI application that demonstrates
the power of Cobra framework for building
command-line tools in Go.`,
    Version: "1.0.0",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    
    // Persistent flags (доступны для всех подкоманд)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myapp.yaml)")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
    
    // Привязка к viper
    viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, _ := os.UserHomeDir()
        viper.AddConfigPath(home)
        viper.SetConfigName(".myapp")
    }
    
    viper.AutomaticEnv()
    viper.ReadInConfig()
}
```

### Пример 2: Подкоманды

```go
// cmd/serve.go
package cmd

import (
    "fmt"
    "net/http"
    
    "github.com/spf13/cobra"
)

var (
    port int
    host string
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the HTTP server",
    Long:  `Start the HTTP server with specified host and port.`,
    Example: `  myapp serve
  myapp serve --port 3000
  myapp serve -p 3000 -H 0.0.0.0`,
    RunE: func(cmd *cobra.Command, args []string) error {
        addr := fmt.Sprintf("%s:%d", host, port)
        fmt.Printf("Starting server on %s\n", addr)
        
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Hello from MyApp!")
        })
        
        return http.ListenAndServe(addr, nil)
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)
    
    serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to listen on")
    serveCmd.Flags().StringVarP(&host, "host", "H", "localhost", "host to bind to")
}
```

```go
// cmd/user.go
package cmd

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
    Use:   "user",
    Short: "Manage users",
    Long:  `Commands for managing users in the system.`,
}

var userListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all users",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Listing users...")
    },
}

var userCreateCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new user",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        email, _ := cmd.Flags().GetString("email")
        fmt.Printf("Creating user: %s (%s)\n", args[0], email)
    },
}

var userDeleteCmd = &cobra.Command{
    Use:   "delete [id]",
    Short: "Delete a user",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        force, _ := cmd.Flags().GetBool("force")
        if !force {
            fmt.Println("Are you sure? Use --force to confirm")
            return
        }
        fmt.Printf("Deleting user: %s\n", args[0])
    },
}

func init() {
    rootCmd.AddCommand(userCmd)
    
    userCmd.AddCommand(userListCmd)
    userCmd.AddCommand(userCreateCmd)
    userCmd.AddCommand(userDeleteCmd)
    
    userCreateCmd.Flags().StringP("email", "e", "", "user email (required)")
    userCreateCmd.MarkFlagRequired("email")
    
    userDeleteCmd.Flags().BoolP("force", "f", false, "force deletion")
}
```

### Пример 3: Интерактивный ввод

```go
// cmd/init.go
package cmd

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    
    "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize a new project",
    Run: func(cmd *cobra.Command, args []string) {
        reader := bufio.NewReader(os.Stdin)
        
        fmt.Print("Project name: ")
        name, _ := reader.ReadString('\n')
        name = strings.TrimSpace(name)
        
        fmt.Print("Description: ")
        desc, _ := reader.ReadString('\n')
        desc = strings.TrimSpace(desc)
        
        fmt.Printf("\nCreating project: %s\n", name)
        fmt.Printf("Description: %s\n", desc)
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
```

### Пример 4: Валидация аргументов

```go
// cmd/deploy.go
package cmd

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
    Use:   "deploy [environment]",
    Short: "Deploy to environment",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return fmt.Errorf("requires exactly one argument: environment")
        }
        
        validEnvs := []string{"dev", "staging", "production"}
        for _, env := range validEnvs {
            if args[0] == env {
                return nil
            }
        }
        
        return fmt.Errorf("invalid environment: %s (valid: %v)", args[0], validEnvs)
    },
    Run: func(cmd *cobra.Command, args []string) {
        env := args[0]
        dry, _ := cmd.Flags().GetBool("dry-run")
        
        if dry {
            fmt.Printf("[DRY RUN] Would deploy to: %s\n", env)
            return
        }
        
        fmt.Printf("Deploying to: %s\n", env)
    },
}

func init() {
    rootCmd.AddCommand(deployCmd)
    deployCmd.Flags().Bool("dry-run", false, "perform a dry run")
}
```

### Пример 5: Полное приложение

```go
// main.go
package main

import "myapp/cmd"

func main() {
    cmd.Execute()
}
```

```go
// cmd/version.go
package cmd

import (
    "fmt"
    "runtime"
    
    "github.com/spf13/cobra"
)

var (
    Version   = "dev"
    Commit    = "none"
    BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Run: func(cmd *cobra.Command, args []string) {
        full, _ := cmd.Flags().GetBool("full")
        
        if full {
            fmt.Printf("Version:    %s\n", Version)
            fmt.Printf("Commit:     %s\n", Commit)
            fmt.Printf("Built:      %s\n", BuildDate)
            fmt.Printf("Go version: %s\n", runtime.Version())
            fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
        } else {
            fmt.Printf("myapp version %s\n", Version)
        }
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
    versionCmd.Flags().Bool("full", false, "show full version info")
}
```

### Пример 6: Автодополнение

```go
// cmd/completion.go
package cmd

import (
    "os"
    
    "github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate completion script",
    Long: `To load completions:

Bash:
  $ source <(myapp completion bash)
  # To load completions for each session:
  # Linux:
  $ myapp completion bash > /etc/bash_completion.d/myapp
  # macOS:
  $ myapp completion bash > $(brew --prefix)/etc/bash_completion.d/myapp

Zsh:
  $ myapp completion zsh > "${fpath[1]}/_myapp"

Fish:
  $ myapp completion fish | source
`,
    Args: cobra.ExactValidArgs(1),
    ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
    Run: func(cmd *cobra.Command, args []string) {
        switch args[0] {
        case "bash":
            rootCmd.GenBashCompletion(os.Stdout)
        case "zsh":
            rootCmd.GenZshCompletion(os.Stdout)
        case "fish":
            rootCmd.GenFishCompletion(os.Stdout, true)
        case "powershell":
            rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
        }
    },
}

func init() {
    rootCmd.AddCommand(completionCmd)
}
```

---

## 🏋️ Практические задания

### Задание 1: Базовая структура CLI

Создайте CLI с корневой командой и подкомандами.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```go
// main.go
package main

import "yourapp/cmd"

func main() {
    cmd.Execute()
}
```

**Баллы:** 10

### Задание 2: Todo CLI

Создайте команды для управления задачами.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Использование:**
```bash
myapp add "Buy groceries" -p high
myapp list --all
myapp complete 1
myapp delete 1
```

**Баллы:** 15

### Задание 3: Конфигурация с Viper

Интегрируйте Viper для конфигурации.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```yaml
# config.yaml
database:
  host: localhost
  port: 5432
  name: myapp

server:
  port: 8080
```

**Баллы:** 15

### Задание 4: Интерактивный режим

Добавьте интерактивный режим с подтверждениями.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 10

### Задание 5: Вывод в разных форматах

Добавьте поддержку JSON/YAML/Table вывода.

**Начальный код:**
```go
package main

import (
    "encoding/json"
    "fmt"
)

// TODO: Работа с JSON

func main() {
    // Ваш код здесь
    
}
```

**Использование:**
```bash
myapp list
myapp list -o json
myapp list -o yaml
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Cobra Documentation](https://cobra.dev/)
- [Cobra GitHub](https://github.com/spf13/cobra)
- [Viper Configuration](https://github.com/spf13/viper)
