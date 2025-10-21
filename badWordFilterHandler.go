package main

import (
	"strings"
)

func badWordFilterHandler(msg string) string {
	cleanWords := []string{}
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitWords := strings.Fields(msg)
	for _, word := range splitWords {
		isBword := false
		for _, bword := range badWords {
			if strings.EqualFold(strings.ToLower(word), strings.ToLower(bword)) {
				cleanWords = append(cleanWords, "****")
				isBword = true
				break
			}
		}
		if !isBword {
			cleanWords = append(cleanWords, word)
		}
	}
	cleanString := strings.Join(cleanWords, " ")
	return cleanString
}
