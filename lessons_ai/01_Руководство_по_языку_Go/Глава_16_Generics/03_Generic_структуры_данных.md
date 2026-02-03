# Generic структуры данных

---

## 💡 Ключевые идеи

1. **Generic-коллекции** — универсальные структуры данных
2. **Type-safe** — безопасность типов на этапе компиляции
3. **Повторное использование** — один код для разных типов
4. **Стандартные структуры** — Stack, Queue, Set, LinkedList, Tree

### Преимущества обобщённых структур

| До Generics | С Generics |
|-------------|------------|
| `interface{}` везде | Типизированные значения |
| Runtime проверки | Compile-time проверки |
| Type assertions | Прямой доступ |
| Возможные panic | Безопасный код |

---

## 📖 Теория

### Зачем нужны Generic-структуры данных?

Представьте, что вам нужен стек (Stack). До Generics было два варианта:

**Вариант 1: Стек для конкретного типа**
```go
type IntStack struct {
    items []int
}

type StringStack struct {
    items []string
}

type UserStack struct {
    items []User
}
// ... и так для каждого типа
```

**Вариант 2: Стек с interface{}**
```go
type Stack struct {
    items []interface{}
}

func (s *Stack) Push(item interface{}) {
    s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
    // возвращает interface{}, нужно приводить тип
}
```

Проблемы второго подхода:
- Нет проверки типов — можно положить int и string в один стек
- Нужно приводить типы при получении: `stack.Pop().(int)`
- Ошибки обнаруживаются только в runtime

### Generic-стек решает все проблемы

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
    // возвращает T — конкретный тип!
}
```

Использование:
```go
intStack := &Stack[int]{}
intStack.Push(1)
intStack.Push(2)
value := intStack.Pop()  // int, не interface{}!

stringStack := &Stack[string]{}
stringStack.Push("hello")
// stringStack.Push(123)  // ОШИБКА КОМПИЛЯЦИИ!
```

### Типобезопасность на этапе компиляции

Главное преимущество Generic-структур — **ошибки обнаруживаются при компиляции**, а не при выполнении:

```go
// С interface{} — скомпилируется, упадёт в runtime
stack := &OldStack{}
stack.Push(123)
name := stack.Pop().(string)  // panic!

// С Generics — не скомпилируется
stack := &Stack[int]{}
stack.Push(123)
name := stack.Pop()  // name имеет тип int
// name + " world"   // ОШИБКА КОМПИЛЯЦИИ: int + string
```

### Какие структуры данных реализовать с Generics?

**Коллекции:**
- **Stack** — последний вошёл, первый вышел (LIFO)
- **Queue** — первый вошёл, первый вышел (FIFO)
- **Set** — множество уникальных элементов
- **LinkedList** — связный список

**Деревья:**
- **BinaryTree** — бинарное дерево
- **BST** — бинарное дерево поиска
- **Heap** — куча (приоритетная очередь)

**Ассоциативные:**
- **OrderedMap** — map с сохранением порядка
- **MultiMap** — map с несколькими значениями для ключа

### Выбор constraint для структуры

| Структура | Constraint | Почему |
|-----------|------------|--------|
| Stack | `any` | Нет операций над элементами |
| Queue | `any` | Нет операций над элементами |
| Set | `comparable` | Нужно проверять уникальность (==) |
| BST | `cmp.Ordered` | Нужно сравнивать (<, >) |
| Heap | `cmp.Ordered` | Нужно сравнивать (<, >) |

### Пустой элемент (zero value)

При возврате элемента из пустой структуры возникает вопрос: что возвращать?

```go
func (s *Stack[T]) Pop() T {
    if len(s.items) == 0 {
        var zero T  // zero value для типа T
        return zero
    }
    // ...
}
```

Или лучше возвращать флаг успеха:
```go
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    // ...
    return item, true
}
```

### Производительность

Generic-структуры в Go так же быстры, как специализированные версии. Go генерирует оптимизированный код для каждого используемого типа на этапе компиляции (мономорфизация).

---

## 💻 Примеры кода

### Пример 1: Set (множество)

```go
package main

import "fmt"

// Обобщённое множество
type Set[T comparable] struct {
    items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
    return &Set[T]{
        items: make(map[T]struct{}),
    }
}

func (s *Set[T]) Add(item T) {
    s.items[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
    delete(s.items, item)
}

func (s *Set[T]) Contains(item T) bool {
    _, ok := s.items[item]
    return ok
}

func (s *Set[T]) Size() int {
    return len(s.items)
}

func (s *Set[T]) ToSlice() []T {
    result := make([]T, 0, len(s.items))
    for item := range s.items {
        result = append(result, item)
    }
    return result
}

// Операции над множествами
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
    result := NewSet[T]()
    for item := range s.items {
        result.Add(item)
    }
    for item := range other.items {
        result.Add(item)
    }
    return result
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
    result := NewSet[T]()
    for item := range s.items {
        if other.Contains(item) {
            result.Add(item)
        }
    }
    return result
}

func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
    result := NewSet[T]()
    for item := range s.items {
        if !other.Contains(item) {
            result.Add(item)
        }
    }
    return result
}

func main() {
    // Set строк
    fruits := NewSet[string]()
    fruits.Add("apple")
    fruits.Add("banana")
    fruits.Add("cherry")
    fruits.Add("apple")  // дубликат игнорируется
    
    fmt.Println("Size:", fruits.Size())        // 3
    fmt.Println("Has apple:", fruits.Contains("apple"))  // true
    fmt.Println("Fruits:", fruits.ToSlice())   // [apple banana cherry]
    
    // Операции
    citrus := NewSet[string]()
    citrus.Add("orange")
    citrus.Add("lemon")
    citrus.Add("apple")  // пересечение
    
    union := fruits.Union(citrus)
    fmt.Println("Union:", union.ToSlice())
    
    intersection := fruits.Intersection(citrus)
    fmt.Println("Intersection:", intersection.ToSlice())  // [apple]
    
    // Set чисел
    numbers := NewSet[int]()
    numbers.Add(1)
    numbers.Add(2)
    numbers.Add(3)
    fmt.Println("Numbers:", numbers.ToSlice())
}
```

### Пример 2: Queue (очередь)

```go
package main

import (
    "errors"
    "fmt"
)

var ErrQueueEmpty = errors.New("queue is empty")

// Обобщённая очередь
type Queue[T any] struct {
    items []T
}

func NewQueue[T any]() *Queue[T] {
    return &Queue[T]{items: make([]T, 0)}
}

func (q *Queue[T]) Enqueue(item T) {
    q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, error) {
    var zero T
    if len(q.items) == 0 {
        return zero, ErrQueueEmpty
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item, nil
}

func (q *Queue[T]) Peek() (T, error) {
    var zero T
    if len(q.items) == 0 {
        return zero, ErrQueueEmpty
    }
    return q.items[0], nil
}

func (q *Queue[T]) IsEmpty() bool {
    return len(q.items) == 0
}

func (q *Queue[T]) Size() int {
    return len(q.items)
}

// Приоритетная очередь
type PriorityQueue[T any] struct {
    items    []T
    lessFunc func(a, b T) bool
}

func NewPriorityQueue[T any](less func(a, b T) bool) *PriorityQueue[T] {
    return &PriorityQueue[T]{
        items:    make([]T, 0),
        lessFunc: less,
    }
}

func (pq *PriorityQueue[T]) Push(item T) {
    pq.items = append(pq.items, item)
    pq.up(len(pq.items) - 1)
}

func (pq *PriorityQueue[T]) Pop() (T, error) {
    var zero T
    if len(pq.items) == 0 {
        return zero, errors.New("queue is empty")
    }
    
    item := pq.items[0]
    last := len(pq.items) - 1
    pq.items[0] = pq.items[last]
    pq.items = pq.items[:last]
    
    if len(pq.items) > 0 {
        pq.down(0)
    }
    
    return item, nil
}

func (pq *PriorityQueue[T]) up(i int) {
    for i > 0 {
        parent := (i - 1) / 2
        if !pq.lessFunc(pq.items[i], pq.items[parent]) {
            break
        }
        pq.items[i], pq.items[parent] = pq.items[parent], pq.items[i]
        i = parent
    }
}

func (pq *PriorityQueue[T]) down(i int) {
    for {
        left := 2*i + 1
        if left >= len(pq.items) {
            break
        }
        
        j := left
        if right := left + 1; right < len(pq.items) && pq.lessFunc(pq.items[right], pq.items[left]) {
            j = right
        }
        
        if !pq.lessFunc(pq.items[j], pq.items[i]) {
            break
        }
        
        pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
        i = j
    }
}

func main() {
    // Простая очередь
    queue := NewQueue[string]()
    queue.Enqueue("first")
    queue.Enqueue("second")
    queue.Enqueue("third")
    
    for !queue.IsEmpty() {
        item, _ := queue.Dequeue()
        fmt.Println(item)  // first, second, third (FIFO)
    }
    
    // Приоритетная очередь (минимум первым)
    pq := NewPriorityQueue[int](func(a, b int) bool {
        return a < b
    })
    
    pq.Push(5)
    pq.Push(2)
    pq.Push(8)
    pq.Push(1)
    
    for len(pq.items) > 0 {
        item, _ := pq.Pop()
        fmt.Print(item, " ")  // 1 2 5 8
    }
    fmt.Println()
}
```

### Пример 3: LinkedList (связный список)

```go
package main

import (
    "errors"
    "fmt"
)

// Узел списка
type Node[T any] struct {
    Value T
    Next  *Node[T]
    Prev  *Node[T]
}

// Двусвязный список
type LinkedList[T any] struct {
    head *Node[T]
    tail *Node[T]
    size int
}

func NewLinkedList[T any]() *LinkedList[T] {
    return &LinkedList[T]{}
}

func (l *LinkedList[T]) PushFront(value T) {
    node := &Node[T]{Value: value}
    
    if l.head == nil {
        l.head = node
        l.tail = node
    } else {
        node.Next = l.head
        l.head.Prev = node
        l.head = node
    }
    l.size++
}

func (l *LinkedList[T]) PushBack(value T) {
    node := &Node[T]{Value: value}
    
    if l.tail == nil {
        l.head = node
        l.tail = node
    } else {
        node.Prev = l.tail
        l.tail.Next = node
        l.tail = node
    }
    l.size++
}

func (l *LinkedList[T]) PopFront() (T, error) {
    var zero T
    if l.head == nil {
        return zero, errors.New("list is empty")
    }
    
    value := l.head.Value
    l.head = l.head.Next
    
    if l.head == nil {
        l.tail = nil
    } else {
        l.head.Prev = nil
    }
    
    l.size--
    return value, nil
}

func (l *LinkedList[T]) PopBack() (T, error) {
    var zero T
    if l.tail == nil {
        return zero, errors.New("list is empty")
    }
    
    value := l.tail.Value
    l.tail = l.tail.Prev
    
    if l.tail == nil {
        l.head = nil
    } else {
        l.tail.Next = nil
    }
    
    l.size--
    return value, nil
}

func (l *LinkedList[T]) Size() int {
    return l.size
}

func (l *LinkedList[T]) IsEmpty() bool {
    return l.size == 0
}

func (l *LinkedList[T]) ToSlice() []T {
    result := make([]T, 0, l.size)
    for node := l.head; node != nil; node = node.Next {
        result = append(result, node.Value)
    }
    return result
}

func (l *LinkedList[T]) ForEach(fn func(T)) {
    for node := l.head; node != nil; node = node.Next {
        fn(node.Value)
    }
}

func main() {
    list := NewLinkedList[int]()
    
    list.PushBack(1)
    list.PushBack(2)
    list.PushBack(3)
    list.PushFront(0)
    
    fmt.Println("List:", list.ToSlice())  // [0 1 2 3]
    
    // Итерация
    list.ForEach(func(v int) {
        fmt.Printf("%d ", v)
    })
    fmt.Println()
    
    // Pop
    front, _ := list.PopFront()
    back, _ := list.PopBack()
    fmt.Printf("Front: %d, Back: %d\n", front, back)  // 0, 3
    fmt.Println("List:", list.ToSlice())  // [1 2]
}
```

### Пример 4: Optional (опциональное значение)

```go
package main

import "fmt"

// Optional — значение, которое может отсутствовать
type Optional[T any] struct {
    value   T
    present bool
}

func Some[T any](value T) Optional[T] {
    return Optional[T]{value: value, present: true}
}

func None[T any]() Optional[T] {
    return Optional[T]{present: false}
}

func (o Optional[T]) IsPresent() bool {
    return o.present
}

func (o Optional[T]) Get() (T, bool) {
    return o.value, o.present
}

func (o Optional[T]) GetOrDefault(defaultValue T) T {
    if o.present {
        return o.value
    }
    return defaultValue
}

func (o Optional[T]) GetOrElse(fn func() T) T {
    if o.present {
        return o.value
    }
    return fn()
}

func (o Optional[T]) Map(fn func(T) T) Optional[T] {
    if o.present {
        return Some(fn(o.value))
    }
    return None[T]()
}

func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
    if o.present && predicate(o.value) {
        return o
    }
    return None[T]()
}

// Функция, возвращающая Optional
func FindUser(id int) Optional[string] {
    users := map[int]string{
        1: "Alice",
        2: "Bob",
    }
    
    if name, ok := users[id]; ok {
        return Some(name)
    }
    return None[string]()
}

func main() {
    // Использование
    user1 := FindUser(1)
    user3 := FindUser(3)
    
    // Проверка наличия
    if user1.IsPresent() {
        name, _ := user1.Get()
        fmt.Println("Found:", name)
    }
    
    // GetOrDefault
    fmt.Println("User 1:", user1.GetOrDefault("Unknown"))  // Alice
    fmt.Println("User 3:", user3.GetOrDefault("Unknown"))  // Unknown
    
    // Map
    upperUser := user1.Map(func(s string) string {
        return "Mr. " + s
    })
    fmt.Println("Upper:", upperUser.GetOrDefault(""))  // Mr. Alice
    
    // Filter
    opt := Some(42)
    filtered := opt.Filter(func(n int) bool {
        return n > 50
    })
    fmt.Println("Filtered present:", filtered.IsPresent())  // false
}
```

### Пример 5: Result (результат с ошибкой)

```go
package main

import (
    "errors"
    "fmt"
    "strconv"
)

// Result — либо значение, либо ошибка
type Result[T any] struct {
    value T
    err   error
}

func Ok[T any](value T) Result[T] {
    return Result[T]{value: value}
}

func Err[T any](err error) Result[T] {
    return Result[T]{err: err}
}

func (r Result[T]) IsOk() bool {
    return r.err == nil
}

func (r Result[T]) IsErr() bool {
    return r.err != nil
}

func (r Result[T]) Unwrap() T {
    if r.err != nil {
        panic(r.err)
    }
    return r.value
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
    if r.err != nil {
        return defaultValue
    }
    return r.value
}

func (r Result[T]) UnwrapOrElse(fn func(error) T) T {
    if r.err != nil {
        return fn(r.err)
    }
    return r.value
}

func (r Result[T]) Map(fn func(T) T) Result[T] {
    if r.err != nil {
        return Err[T](r.err)
    }
    return Ok(fn(r.value))
}

func (r Result[T]) AndThen(fn func(T) Result[T]) Result[T] {
    if r.err != nil {
        return Err[T](r.err)
    }
    return fn(r.value)
}

// Пример использования
func ParseInt(s string) Result[int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return Err[int](err)
    }
    return Ok(n)
}

func Divide(a, b int) Result[int] {
    if b == 0 {
        return Err[int](errors.New("division by zero"))
    }
    return Ok(a / b)
}

func main() {
    // Успешный результат
    result := ParseInt("42")
    fmt.Println("Is OK:", result.IsOk())     // true
    fmt.Println("Value:", result.Unwrap())   // 42
    
    // Ошибочный результат
    errResult := ParseInt("abc")
    fmt.Println("Is Err:", errResult.IsErr())           // true
    fmt.Println("Default:", errResult.UnwrapOr(0))      // 0
    
    // Цепочка операций
    finalResult := ParseInt("100").
        AndThen(func(n int) Result[int] {
            return Divide(n, 5)
        }).
        Map(func(n int) int {
            return n * 2
        })
    
    fmt.Println("Final:", finalResult.UnwrapOr(-1))  // 40
    
    // Обработка ошибки в цепочке
    errorChain := ParseInt("100").
        AndThen(func(n int) Result[int] {
            return Divide(n, 0)  // ошибка!
        }).
        Map(func(n int) int {
            return n * 2  // не выполнится
        })
    
    fmt.Println("Error chain:", errorChain.UnwrapOr(-1))  // -1
}
```

### Пример 6: Tree (бинарное дерево)

```go
package main

import (
    "cmp"
    "fmt"
)

type TreeNode[T cmp.Ordered] struct {
    Value T
    Left  *TreeNode[T]
    Right *TreeNode[T]
}

type BinarySearchTree[T cmp.Ordered] struct {
    root *TreeNode[T]
    size int
}

func NewBST[T cmp.Ordered]() *BinarySearchTree[T] {
    return &BinarySearchTree[T]{}
}

func (t *BinarySearchTree[T]) Insert(value T) {
    t.root = t.insertNode(t.root, value)
    t.size++
}

func (t *BinarySearchTree[T]) insertNode(node *TreeNode[T], value T) *TreeNode[T] {
    if node == nil {
        return &TreeNode[T]{Value: value}
    }
    
    if value < node.Value {
        node.Left = t.insertNode(node.Left, value)
    } else if value > node.Value {
        node.Right = t.insertNode(node.Right, value)
    }
    
    return node
}

func (t *BinarySearchTree[T]) Contains(value T) bool {
    return t.containsNode(t.root, value)
}

func (t *BinarySearchTree[T]) containsNode(node *TreeNode[T], value T) bool {
    if node == nil {
        return false
    }
    
    if value < node.Value {
        return t.containsNode(node.Left, value)
    } else if value > node.Value {
        return t.containsNode(node.Right, value)
    }
    
    return true
}

func (t *BinarySearchTree[T]) InOrder() []T {
    result := make([]T, 0, t.size)
    t.inOrderTraversal(t.root, &result)
    return result
}

func (t *BinarySearchTree[T]) inOrderTraversal(node *TreeNode[T], result *[]T) {
    if node == nil {
        return
    }
    
    t.inOrderTraversal(node.Left, result)
    *result = append(*result, node.Value)
    t.inOrderTraversal(node.Right, result)
}

func (t *BinarySearchTree[T]) Min() (T, bool) {
    var zero T
    if t.root == nil {
        return zero, false
    }
    
    node := t.root
    for node.Left != nil {
        node = node.Left
    }
    return node.Value, true
}

func (t *BinarySearchTree[T]) Max() (T, bool) {
    var zero T
    if t.root == nil {
        return zero, false
    }
    
    node := t.root
    for node.Right != nil {
        node = node.Right
    }
    return node.Value, true
}

func main() {
    // BST с числами
    tree := NewBST[int]()
    tree.Insert(5)
    tree.Insert(3)
    tree.Insert(7)
    tree.Insert(1)
    tree.Insert(9)
    tree.Insert(4)
    
    fmt.Println("Contains 7:", tree.Contains(7))  // true
    fmt.Println("Contains 10:", tree.Contains(10)) // false
    
    fmt.Println("In-order:", tree.InOrder())  // [1 3 4 5 7 9]
    
    min, _ := tree.Min()
    max, _ := tree.Max()
    fmt.Printf("Min: %d, Max: %d\n", min, max)  // 1, 9
    
    // BST со строками
    stringTree := NewBST[string]()
    stringTree.Insert("banana")
    stringTree.Insert("apple")
    stringTree.Insert("cherry")
    
    fmt.Println("Strings:", stringTree.InOrder())  // [apple banana cherry]
}
```

---

## ⚠️ Частые ошибки

### 1. Возврат zero value без указания

```go
// ❌ ПРОБЛЕМА — nil не всегда работает
func (q *Queue[T]) Dequeue() T {
    if q.IsEmpty() {
        return nil  // ошибка для int, string и т.д.
    }
}

// ✅ ПРАВИЛЬНО — zero value
func (q *Queue[T]) Dequeue() (T, bool) {
    var zero T
    if q.IsEmpty() {
        return zero, false
    }
    // ...
    return value, true
}
```

### 2. Сравнение с any

```go
// ❌ ОШИБКА — any нельзя сравнивать
type Set[T any] struct { ... }
func (s *Set[T]) Contains(item T) bool {
    return s.items[item]  // ошибка!
}

// ✅ ПРАВИЛЬНО — comparable
type Set[T comparable] struct { ... }
```

---

## 🏋️ Практические задания

### Задание 1: Generic Queue

Реализуйте обобщённую очередь (FIFO).

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Size: 3
Dequeue: first
Dequeue: second
Dequeue: third
```

**Баллы:** 10

### Задание 2: Generic Set

Реализуйте множество с базовыми операциями.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 15

### Задание 3: Generic Pair и Triple

Создайте обобщённые типы для пар и троек значений.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Pair: {First:name Second:42}
Swapped: {First:42 Second:name}
Alice is 30 years old
Bob is 25 years old
Charlie is 35 years old
```

**Баллы:** 10

### Задание 4: Optional type

Реализуйте Optional для безопасной работы с опциональными значениями.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
User 1: Alice
User 3: Unknown
Hello, Alice!
No greeting
```

**Баллы:** 15

### Задание 5: Generic Binary Search Tree

Реализуйте бинарное дерево поиска.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Contains 4: true
Contains 10: false
In-order: [1 3 4 5 6 7 8]
```

**Баллы:** 15

---

## 🔗 Полезные ссылки

- [Go Generic Data Structures](https://github.com/zyedidia/generic)
- [Go Collections](https://pkg.go.dev/golang.org/x/exp/slices)
