# MongoDB

## 💡 Ключевые идеи

1. **NoSQL** — документо-ориентированная база данных
2. **BSON** — бинарный JSON для хранения данных
3. **Коллекции** — аналог таблиц в SQL
4. **Документы** — записи в формате JSON/BSON
5. **Официальный драйвер** — mongo-go-driver

---

## 📋 Синтаксис

### Установка драйвера

```bash
go get go.mongodb.org/mongo-driver/mongo
```

### Подключение

```go
client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
defer client.Disconnect(ctx)

// Получение коллекции
collection := client.Database("dbname").Collection("collname")
```

### Основные операции

```go
// Вставка
InsertOne(ctx, document)
InsertMany(ctx, documents)

// Поиск
FindOne(ctx, filter)
Find(ctx, filter)

// Обновление
UpdateOne(ctx, filter, update)
UpdateMany(ctx, filter, update)

// Удаление
DeleteOne(ctx, filter)
DeleteMany(ctx, filter)
```

---

## 💻 Примеры кода

### Подключение к MongoDB

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    // Контекст с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Подключение
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    
    // Проверка подключения
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Connected to MongoDB!")
}
```

### Определение модели данных

```go
package main

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    Model   string             `bson:"model"`
    Company string             `bson:"company"`
    Price   int                `bson:"price"`
}
```

### INSERT — добавление документа

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    Model   string             `bson:"model"`
    Company string             `bson:"company"`
    Price   int                `bson:"price"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Один документ
    product := Product{
        Model:   "iPhone 15",
        Company: "Apple",
        Price:   99000,
    }
    
    result, err := collection.InsertOne(ctx, product)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Inserted ID:", result.InsertedID)
}
```

### Множественная вставка

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
    Model   string `bson:"model"`
    Company string `bson:"company"`
    Price   int    `bson:"price"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    products := []interface{}{
        Product{Model: "Pixel 8", Company: "Google", Price: 75000},
        Product{Model: "Galaxy S24", Company: "Samsung", Price: 85000},
        Product{Model: "Mi 14", Company: "Xiaomi", Price: 55000},
    }
    
    result, err := collection.InsertMany(ctx, products)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Inserted %d documents\n", len(result.InsertedIDs))
}
```

### FIND — поиск документов

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    Model   string             `bson:"model"`
    Company string             `bson:"company"`
    Price   int                `bson:"price"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Поиск всех документов
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(ctx)
    
    var products []Product
    
    // Декодируем все результаты
    if err = cursor.All(ctx, &products); err != nil {
        log.Fatal(err)
    }
    
    for _, p := range products {
        fmt.Printf("%s (%s) - %d руб.\n", p.Model, p.Company, p.Price)
    }
}
```

### Поиск с фильтром

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    Model   string             `bson:"model"`
    Company string             `bson:"company"`
    Price   int                `bson:"price"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Фильтр: цена > 60000
    filter := bson.M{
        "price": bson.M{"$gt": 60000},
    }
    
    cursor, _ := collection.Find(ctx, filter)
    defer cursor.Close(ctx)
    
    fmt.Println("Products over 60000:")
    for cursor.Next(ctx) {
        var p Product
        cursor.Decode(&p)
        fmt.Printf("  %s - %d\n", p.Model, p.Price)
    }
}
```

### FindOne — один документ

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    Model   string             `bson:"model"`
    Company string             `bson:"company"`
    Price   int                `bson:"price"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    var product Product
    
    filter := bson.M{"model": "iPhone 15"}
    err := collection.FindOne(ctx, filter).Decode(&product)
    
    if err == mongo.ErrNoDocuments {
        fmt.Println("Product not found")
        return
    }
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found: %+v\n", product)
}
```

### Поиск по ID

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    Model   string             `bson:"model"`
    Company string             `bson:"company"`
    Price   int                `bson:"price"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Преобразуем строку в ObjectID
    objID, err := primitive.ObjectIDFromHex("65abc123def456789012345")
    if err != nil {
        log.Fatal("Invalid ObjectID:", err)
    }
    
    var product Product
    filter := bson.M{"_id": objID}
    err = collection.FindOne(ctx, filter).Decode(&product)
    
    if err == mongo.ErrNoDocuments {
        fmt.Println("Product not found")
        return
    }
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found: %+v\n", product)
}
```

### UPDATE — обновление

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Обновление одного документа
    filter := bson.M{"model": "iPhone 15"}
    update := bson.M{
        "$set": bson.M{"price": 95000},
    }
    
    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Matched: %d, Modified: %d\n", 
        result.MatchedCount, result.ModifiedCount)
}
```

### Обновление нескольких документов

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Увеличить цену на 10% для всех Apple
    filter := bson.M{"company": "Apple"}
    update := bson.M{
        "$mul": bson.M{"price": 1.1},
    }
    
    result, err := collection.UpdateMany(ctx, filter, update)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Updated %d documents\n", result.ModifiedCount)
}
```

### DELETE — удаление

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Удаление одного документа
    filter := bson.M{"model": "Mi 14"}
    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Deleted %d document(s)\n", result.DeletedCount)
    
    // Удаление нескольких
    filter = bson.M{"price": bson.M{"$lt": 50000}}
    result, _ = collection.DeleteMany(ctx, filter)
    fmt.Printf("Deleted %d cheap products\n", result.DeletedCount)
}
```

### Агрегация

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Агрегация: средняя цена по компаниям
    pipeline := []bson.M{
        {"$group": bson.M{
            "_id":       "$company",
            "avgPrice":  bson.M{"$avg": "$price"},
            "count":     bson.M{"$sum": 1},
        }},
        {"$sort": bson.M{"avgPrice": -1}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(ctx)
    
    fmt.Println("Average price by company:")
    for cursor.Next(ctx) {
        var result bson.M
        cursor.Decode(&result)
        fmt.Printf("  %s: %.0f руб. (%v products)\n", 
            result["_id"], result["avgPrice"], result["count"])
    }
}
```

### Индексы

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    defer client.Disconnect(ctx)
    
    collection := client.Database("productdb").Collection("products")
    
    // Создание индекса
    indexModel := mongo.IndexModel{
        Keys:    bson.D{{Key: "model", Value: 1}},
        Options: options.Index().SetUnique(true),
    }
    
    name, err := collection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        fmt.Println("Index error:", err)
        return
    }
    
    fmt.Println("Created index:", name)
}
```

---

## ⚠️ Частые ошибки

### 1. Забыли контекст

```go
// ❌ Операции могут зависнуть
collection.Find(nil, filter)

// ✅ Используйте контекст с таймаутом
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
collection.Find(ctx, filter)
```

### 2. Не закрывают cursor

```go
// ❌ Утечка ресурсов
cursor, _ := collection.Find(ctx, filter)
// cursor не закрыт!

// ✅ Всегда закрывайте cursor
defer cursor.Close(ctx)
```

### 3. Неверные BSON теги

```go
// ❌ Теги игнорируются
type Product struct {
    ID    string `json:"_id"`  // json, не bson!
    Model string `json:"model"`
}

// ✅ Используйте bson теги
type Product struct {
    ID    primitive.ObjectID `bson:"_id,omitempty"`
    Model string             `bson:"model"`
}
```

---

## 📝 Практика

### Задача 1: Product catalog
CRUD для каталога товаров.

### Задача 2: User profiles
Хранение профилей с вложенными документами.

### Задача 3: Blog posts
Блог с тегами и комментариями.

### Задача 4: Analytics
Сбор и анализ событий.

### Задача 5: Search
Текстовый поиск с индексами.

### Задача 6: Geolocation
Геопространственные запросы.

### Задача 7: Time series
Временные ряды данных.

### Задача 8: Transactions
Многодокументные транзакции.
