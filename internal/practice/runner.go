package practice

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	// MaxCodeSize — максимальный размер кода (100KB).
	MaxCodeSize = 100 * 1024
	// RunTimeout — таймаут выполнения (15 секунд).
	RunTimeout = 15 * time.Second
)

// RunResult — результат выполнения кода.
type RunResult struct {
	Success bool
	Stdout  string
	Stderr  string
	Error   string
}

// Runner — интерфейс для выполнения Go-кода.
type Runner interface {
	Run(ctx context.Context, code string) (*RunResult, error)
	Check(ctx context.Context, code string, testsGo string) (*RunResult, error)
}

// LocalRunner — локальный runner (выполняет код через go run/test).
type LocalRunner struct {
	tempDir string
}

// NewLocalRunner создаёт новый локальный runner.
func NewLocalRunner() *LocalRunner {
	return &LocalRunner{}
}

// Run выполняет Go-код и возвращает результат.
func (r *LocalRunner) Run(ctx context.Context, code string) (*RunResult, error) {
	// Проверяем размер кода
	if len(code) > MaxCodeSize {
		return &RunResult{
			Success: false,
			Error:   fmt.Sprintf("Код слишком большой: %d байт (максимум %d)", len(code), MaxCodeSize),
		}, nil
	}

	// Создаём временную директорию
	tempDir, err := os.MkdirTemp("", "gorun-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Записываем код в файл
	mainFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("write main.go: %w", err)
	}

	// Создаём go.mod
	goMod := "module runner\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return nil, fmt.Errorf("write go.mod: %w", err)
	}

	// Устанавливаем таймаут
	ctx, cancel := context.WithTimeout(ctx, RunTimeout)
	defer cancel()

	// Запускаем go run
	cmd := exec.CommandContext(ctx, "go", "run", "main.go")
	cmd.Dir = tempDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	result := &RunResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Success = false
		result.Error = fmt.Sprintf("Превышено время выполнения (%v)", RunTimeout)
		return result, nil
	}

	if err != nil {
		result.Success = false
		if result.Stderr != "" {
			result.Error = result.Stderr
		} else {
			result.Error = err.Error()
		}
		return result, nil
	}

	result.Success = true
	return result, nil
}

// Check проверяет код с помощью тестов.
func (r *LocalRunner) Check(ctx context.Context, code string, testsGo string) (*RunResult, error) {
	// Проверяем размер кода
	if len(code) > MaxCodeSize {
		return &RunResult{
			Success: false,
			Error:   fmt.Sprintf("Код слишком большой: %d байт (максимум %d)", len(code), MaxCodeSize),
		}, nil
	}

	// Создаём временную директорию
	tempDir, err := os.MkdirTemp("", "gocheck-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Записываем код пользователя
	mainFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("write main.go: %w", err)
	}

	// Записываем тесты
	testFile := filepath.Join(tempDir, "main_test.go")
	if err := os.WriteFile(testFile, []byte(testsGo), 0644); err != nil {
		return nil, fmt.Errorf("write main_test.go: %w", err)
	}

	// Создаём go.mod
	goMod := "module runner\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return nil, fmt.Errorf("write go.mod: %w", err)
	}

	// Устанавливаем таймаут
	ctx, cancel := context.WithTimeout(ctx, RunTimeout)
	defer cancel()

	// Запускаем go test
	cmd := exec.CommandContext(ctx, "go", "test", "-v", ".")
	cmd.Dir = tempDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	result := &RunResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if ctx.Err() == context.DeadlineExceeded {
		result.Success = false
		result.Error = fmt.Sprintf("Превышено время выполнения (%v)", RunTimeout)
		return result, nil
	}

	if err != nil {
		result.Success = false
		// Для тестов ошибка обычно означает, что тесты не прошли
		if result.Stdout != "" {
			result.Error = result.Stdout
		} else if result.Stderr != "" {
			result.Error = result.Stderr
		} else {
			result.Error = err.Error()
		}
		return result, nil
	}

	result.Success = true
	return result, nil
}
