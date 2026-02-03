# Введение в gRPC

---

## 💡 Ключевые идеи

1. **gRPC** — высокопроизводительный RPC framework от Google
2. **Protocol Buffers** — язык описания интерфейсов и сериализации
3. **HTTP/2** — мультиплексирование, streaming, сжатие
4. **Типы вызовов** — Unary, Server/Client/Bidirectional Streaming
5. **Code Generation** — автогенерация кода из .proto файлов

### gRPC vs REST

| Аспект | gRPC | REST |
|--------|------|------|
| Протокол | HTTP/2 | HTTP/1.1 |
| Формат | Protobuf (binary) | JSON (text) |
| Контракт | Строгий (.proto) | OpenAPI (опционально) |
| Streaming | Встроенный | WebSocket |
| Скорость | Быстрее | Медленнее |

---

## 📖 Теория

### Что такое gRPC?

**gRPC** — это современный RPC-фреймворк с открытым исходным кодом от Google. RPC (Remote Procedure Call) позволяет вызывать функции на удалённом сервере так же просто, как локальные.

```go
// Выглядит как локальный вызов
user, err := userClient.GetUser(ctx, &GetUserRequest{Id: 1})
// Но на самом деле это сетевой запрос к другому сервису
```

### Почему gRPC популярен?

1. **Скорость** — бинарный формат Protobuf в 3-10 раз быстрее JSON
2. **Строгая типизация** — контракт определён в .proto файле
3. **Code Generation** — клиент и сервер генерируются автоматически
4. **HTTP/2** — мультиплексирование, сжатие заголовков
5. **Streaming** — потоковая передача данных в обе стороны
6. **Cross-language** — Go, Java, Python, C++, и другие

### Protocol Buffers (Protobuf)

**Protobuf** — это язык описания интерфейсов и формат сериализации от Google.

```protobuf
syntax = "proto3";

message User {
  int64 id = 1;      // номер поля (не значение!)
  string name = 2;
  string email = 3;
}
```

Числа `1, 2, 3` — это **номера полей**, не значения. Они используются для компактной бинарной сериализации.

### Почему Protobuf быстрее JSON?

**JSON:**
```json
{"id":12345,"name":"John","email":"john@example.com"}
```
- Текстовый формат
- Имена полей повторяются
- Числа как строки

**Protobuf (бинарный):**
```
08 b9 60 12 04 John 1a 10 john@example.com
```
- Бинарный формат
- Только номера полей и значения
- Числа в компактном виде

### Типы вызовов в gRPC

**1. Unary RPC** — один запрос, один ответ:
```protobuf
rpc GetUser(GetUserRequest) returns (User);
```
Как обычный HTTP запрос.

**2. Server Streaming** — один запрос, поток ответов:
```protobuf
rpc ListUsers(ListUsersRequest) returns (stream User);
```
Сервер отправляет много сообщений (например, результаты поиска).

**3. Client Streaming** — поток запросов, один ответ:
```protobuf
rpc UploadFile(stream Chunk) returns (UploadStatus);
```
Клиент отправляет много сообщений (загрузка файла по частям).

**4. Bidirectional Streaming** — потоки в обе стороны:
```protobuf
rpc Chat(stream Message) returns (stream Message);
```
Чат, игры в реальном времени.

### Когда использовать gRPC?

✅ **Используйте когда:**
- Микросервисы общаются друг с другом
- Важна производительность
- Нужен streaming
- Строгие контракты между командами

❌ **Используйте REST когда:**
- Публичный API (браузеры не поддерживают gRPC напрямую)
- Простые CRUD операции
- Нужна читаемость запросов для отладки

### gRPC-Web и gRPC-Gateway

**Проблема:** Браузеры не поддерживают HTTP/2 для gRPC напрямую.

**Решения:**
- **gRPC-Web** — прокси, переводящий gRPC в формат, понятный браузерам
- **gRPC-Gateway** — генерирует REST API из .proto файлов

```
Browser → REST → grpc-gateway → gRPC → Backend
```

### Workflow разработки с gRPC

1. Определить сервис в `.proto` файле
2. Сгенерировать код: `protoc --go_out=. --go-grpc_out=. user.proto`
3. Реализовать сервер (имплементировать интерфейс)
4. Создать клиент (использовать сгенерированный stub)

---

## 📋 Установка

```bash
# Protobuf compiler
# macOS
brew install protobuf

# Ubuntu
apt install -y protobuf-compiler

# Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

---

## 💻 Примеры кода

### Пример 1: Определение сервиса (.proto)

```protobuf
// proto/user/user.proto
syntax = "proto3";

package user;

option go_package = "myapp/gen/user";

// Сообщения
message User {
  int64 id = 1;
  string email = 2;
  string name = 3;
  int32 age = 4;
  repeated string roles = 5;
  UserStatus status = 6;
}

enum UserStatus {
  USER_STATUS_UNSPECIFIED = 0;
  USER_STATUS_ACTIVE = 1;
  USER_STATUS_INACTIVE = 2;
  USER_STATUS_BANNED = 3;
}

// Запросы и ответы
message CreateUserRequest {
  string email = 1;
  string name = 2;
  string password = 3;
}

message CreateUserResponse {
  User user = 1;
}

message GetUserRequest {
  int64 id = 1;
}

message GetUserResponse {
  User user = 1;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}

message UpdateUserRequest {
  int64 id = 1;
  optional string name = 2;
  optional int32 age = 3;
}

message UpdateUserResponse {
  User user = 1;
}

message DeleteUserRequest {
  int64 id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}

// Сервис
service UserService {
  // Unary RPC
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  
  // Server streaming
  rpc ListUsers(ListUsersRequest) returns (stream User);
  
  // Client streaming
  rpc BatchCreateUsers(stream CreateUserRequest) returns (ListUsersResponse);
  
  // Bidirectional streaming
  rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}

message ChatMessage {
  int64 user_id = 1;
  string message = 2;
  int64 timestamp = 3;
}
```

### Пример 2: Генерация кода

```bash
# Makefile
.PHONY: proto

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/user.proto
```

### Пример 3: Реализация сервера

```go
// internal/server/user_server.go
package server

import (
    "context"
    "io"
    "sync"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    pb "myapp/gen/user"
)

type UserServer struct {
    pb.UnimplementedUserServiceServer
    
    mu    sync.RWMutex
    users map[int64]*pb.User
    nextID int64
}

func NewUserServer() *UserServer {
    return &UserServer{
        users:  make(map[int64]*pb.User),
        nextID: 1,
    }
}

// Unary RPC
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Валидация
    if req.Email == "" {
        return nil, status.Error(codes.InvalidArgument, "email is required")
    }
    
    user := &pb.User{
        Id:     s.nextID,
        Email:  req.Email,
        Name:   req.Name,
        Status: pb.UserStatus_USER_STATUS_ACTIVE,
    }
    
    s.users[user.Id] = user
    s.nextID++
    
    return &pb.CreateUserResponse{User: user}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    user, ok := s.users[req.Id]
    if !ok {
        return nil, status.Error(codes.NotFound, "user not found")
    }
    
    return &pb.GetUserResponse{User: user}, nil
}

// Server streaming
func (s *UserServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    for _, user := range s.users {
        if err := stream.Send(user); err != nil {
            return err
        }
    }
    
    return nil
}

// Client streaming
func (s *UserServer) BatchCreateUsers(stream pb.UserService_BatchCreateUsersServer) error {
    var users []*pb.User
    
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            // Клиент завершил отправку
            return stream.SendAndClose(&pb.ListUsersResponse{
                Users: users,
                Total: int32(len(users)),
            })
        }
        if err != nil {
            return err
        }
        
        // Создаём пользователя
        s.mu.Lock()
        user := &pb.User{
            Id:     s.nextID,
            Email:  req.Email,
            Name:   req.Name,
            Status: pb.UserStatus_USER_STATUS_ACTIVE,
        }
        s.users[user.Id] = user
        s.nextID++
        s.mu.Unlock()
        
        users = append(users, user)
    }
}

// Bidirectional streaming
func (s *UserServer) Chat(stream pb.UserService_ChatServer) error {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        
        // Эхо-ответ
        response := &pb.ChatMessage{
            UserId:    msg.UserId,
            Message:   "Server received: " + msg.Message,
            Timestamp: time.Now().Unix(),
        }
        
        if err := stream.Send(response); err != nil {
            return err
        }
    }
}
```

### Пример 4: Запуск сервера

```go
// cmd/server/main.go
package main

import (
    "log"
    "net"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    
    pb "myapp/gen/user"
    "myapp/internal/server"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    // Создаём gRPC сервер
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),
    )
    
    // Регистрируем сервис
    userServer := server.NewUserServer()
    pb.RegisterUserServiceServer(grpcServer, userServer)
    
    // Включаем reflection для отладки
    reflection.Register(grpcServer)
    
    log.Println("gRPC server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

// Interceptor для логирования
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    
    resp, err := handler(ctx, req)
    
    log.Printf("Method: %s, Duration: %v, Error: %v",
        info.FullMethod, time.Since(start), err)
    
    return resp, err
}
```

### Пример 5: Клиент

```go
// cmd/client/main.go
package main

import (
    "context"
    "io"
    "log"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    pb "myapp/gen/user"
)

func main() {
    // Подключение
    conn, err := grpc.Dial(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }
    defer conn.Close()
    
    client := pb.NewUserServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Unary call
    createResp, err := client.CreateUser(ctx, &pb.CreateUserRequest{
        Email: "john@example.com",
        Name:  "John Doe",
    })
    if err != nil {
        log.Fatalf("CreateUser failed: %v", err)
    }
    log.Printf("Created user: %v", createResp.User)
    
    // Server streaming
    stream, err := client.ListUsers(ctx, &pb.ListUsersRequest{})
    if err != nil {
        log.Fatalf("ListUsers failed: %v", err)
    }
    
    for {
        user, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("stream error: %v", err)
        }
        log.Printf("User: %v", user)
    }
}
```

### Пример 6: Interceptors

```go
// internal/interceptors/interceptors.go
package interceptors

import (
    "context"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
)

// Auth interceptor
func AuthInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // Получаем metadata
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "missing metadata")
    }
    
    // Проверяем токен
    tokens := md.Get("authorization")
    if len(tokens) == 0 {
        return nil, status.Error(codes.Unauthenticated, "missing token")
    }
    
    // Валидация токена
    if tokens[0] != "Bearer valid-token" {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }
    
    return handler(ctx, req)
}

// Recovery interceptor
func RecoveryInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (resp interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = status.Errorf(codes.Internal, "panic: %v", r)
        }
    }()
    
    return handler(ctx, req)
}

// Timeout interceptor
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        ctx, cancel := context.WithTimeout(ctx, timeout)
        defer cancel()
        
        return handler(ctx, req)
    }
}

// Chain interceptors
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        chain := handler
        for i := len(interceptors) - 1; i >= 0; i-- {
            interceptor := interceptors[i]
            next := chain
            chain = func(ctx context.Context, req interface{}) (interface{}, error) {
                return interceptor(ctx, req, info, next)
            }
        }
        return chain(ctx, req)
    }
}
```

---

## 🏋️ Практические задания

### Задание 1: Proto файл и генерация

Создайте .proto файл и сгенерируйте код.

**Начальный код:**
```protobuf
// proto/user.proto
syntax = "proto3";

package user;
option go_package = "github.com/yourname/app/proto/user";

message User {
    int64 id = 1;
    string name = 2;
    string email = 3;
}

message GetUserRequest {
    int64 id = 1;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
}

message ListUsersRequest {
    int32 page = 1;
    int32 page_size = 2;
}

message ListUsersResponse {
    repeated User users = 1;
    int32 total = 2;
}

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}
```

**Команда генерации:**
```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

**Баллы:** 10

### Задание 2: gRPC Server

Реализуйте gRPC сервер.

**Начальный код:**
```go
package main

import (
    "errors"
    "fmt"
    "net/http"
    "os"
)

// TODO: Реализуйте HTTP обработчик

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 15

### Задание 3: gRPC Client

Создайте клиент для вызова сервиса.

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

### Задание 4: Interceptor для логирования

Добавьте interceptor для логирования запросов.

**Начальный код:**
```go
package main

import (
    "fmt"
    "net/http"
)

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 10

### Задание 5: Server Streaming

Реализуйте server-side streaming.

**Начальный код:**
```protobuf
service NotificationService {
    rpc Subscribe(SubscribeRequest) returns (stream Notification);
}
```

```go
func (s *notificationServer) Subscribe(
    req *pb.SubscribeRequest,
    stream pb.NotificationService_SubscribeServer,
) error {
    for i := 0; i < 5; i++ {
        notification := &pb.Notification{
            Id:      int64(i),
            Message: fmt.Sprintf("Notification %d", i),
        }
        
        if err := stream.Send(notification); err != nil {
            return err
        }
        
        time.Sleep(time.Second)
    }
    return nil
}
```

**Баллы:** 15

---

## 🔗 Полезные ссылки

- [gRPC Documentation](https://grpc.io/docs/languages/go/)
- [Protocol Buffers](https://protobuf.dev/)
- [grpc-go](https://github.com/grpc/grpc-go)
