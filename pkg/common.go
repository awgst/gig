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

// PluralizeSnakeCase converts snake_case to plural snake_case
// Accepts snakeCase as string
// Returns string
// Example: snake_case -> snake_cases
func PluralizeSnakeCase(word string) string {
	// Check for some common pluralization rules
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") ||
		strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}

	if strings.HasSuffix(word, "y") && !isVowel(word[len(word)-2]) {
		return word[:len(word)-1] + "ies"
	}

	return word + "s"
}

// Helper function to check if a character is a vowel
func isVowel(c byte) bool {
	vowels := []byte{'a', 'e', 'i', 'o', 'u'}
	for _, v := range vowels {
		if v == c {
			return true
		}
	}
	return false
}
