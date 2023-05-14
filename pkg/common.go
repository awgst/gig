package pkg

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SnakeToPascal(snakeCase string) string {
	words := strings.Split(snakeCase, "_")
	caser := cases.Title(language.AmericanEnglish)

	for i := 0; i < len(words); i++ {
		words[i] = StringTitle(caser.String(words[i]))
	}

	return strings.Join(words, "")
}

func StringTitle(str string) string {
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(str)
}
