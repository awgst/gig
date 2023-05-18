// Package pkg implements list function and variable that can be used by other packages
package pkg

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// SnakeToCamel converts snake_case to camelCase
// Accepts snakeCase as string
// Returns string
// Example: snake_case -> SnakeCase
func SnakeToPascal(snakeCase string) string {
	words := strings.Split(snakeCase, "_")
	caser := cases.Title(language.AmericanEnglish)

	for i := 0; i < len(words); i++ {
		words[i] = StringTitle(caser.String(words[i]))
	}

	return strings.Join(words, "")
}

// StringTitle converts a string to title case
// Accepts str as string
// Returns string
// Example: hello world -> Hello World
func StringTitle(str string) string {
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(str)
}

// SnakeToCamel converts snake_case to camelCase
// Accepts snakeCase as string
// Returns string
// Example: snake_case -> snakeCase
func SnakeToCamel(snakeCase string) string {
	words := strings.Split(snakeCase, "_")
	caser := cases.Title(language.AmericanEnglish)

	for i := 1; i < len(words); i++ {
		words[i] = StringTitle(caser.String(words[i]))
	}

	return strings.Join(words, "")
}
