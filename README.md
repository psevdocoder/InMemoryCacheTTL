# Golang in-memory cache

Second training project (Goroutines and channels).

Provides a simple in-memory cache with TTL support.

## Installation
```bash
go get "github.com/psevdocoder/InMemoryCacheTTL"
```

## Usage example

```go
package main

import (
	"fmt"
	cache "github.com/psevdocoder/InMemoryCacheTTL"
	"time"
)

func main() {
	c := cache.New(time.Second * 5)

	// Добавляем значение в кэш с TTL
	c.Set("key1", "value1", 4*time.Second)

	// Получаем значение из кэша
	result, ok := c.Get("key1")
	if ok {
		fmt.Println("Value:", result)
	} else {
		fmt.Println("Key not found")
	}

	// Ждем несколько секунд, чтобы TTL истек
	time.Sleep(5 * time.Second)

	// Попытка получить значение после истечения TTL
	result, ok = c.Get("key1")
	if ok {
		fmt.Println("Value:", result)
	} else {
		fmt.Println("Key not found (TTL expired)")
	}

}

```