package main

import (
	"fmt"
	"strings"
)

func main() {
	// strings package examples
	text := "Hello, Go Programming"

	fmt.Println("Contains:", strings.Contains(text, "Go"))
	fmt.Println("Count:", strings.Count(text, "o"))
	fmt.Println("HasPrefix:", strings.HasPrefix(text, "Hello"))
	fmt.Println("HasSuffix:", strings.HasSuffix(text, "Programming"))
	fmt.Println("Index:", strings.Index(text, "Go"))
	fmt.Println("Join:", strings.Join([]string{"Hello", "World"}, " "))
	fmt.Println("Repeat:", strings.Repeat("Go", 3))
	fmt.Println("Replace:", strings.Replace(text, "Go", "Golang", 1))
	fmt.Println("Split:", strings.Split(text, ","))
	fmt.Println("ToLower:", strings.ToLower(text))
	fmt.Println("ToUpper:", strings.ToUpper(text))
	fmt.Println("Trim:", strings.Trim("  spaces  ", " "))
}
