package helper

import (
	"regexp"
	"strings"
)

var stopWords = map[string]bool{
	"va":  true,
	"ni":  true,
	"da":  true,
	"the": true,
	"is":  true,
	"in":  true,
	"on":  true,
}

var re = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func GenerateTags(title, desc string) []string {
	text := strings.ToLower(title + " " + desc)

	text = re.ReplaceAllString(text, " ")

	words := strings.Fields(text)

	tagMap := make(map[string]bool)
	var tags []string

	for _, w := range words {

		if len(w) < 3 {
			continue
		}

		if stopWords[w] {
			continue
		}

		if !tagMap[w] {
			tagMap[w] = true
			tags = append(tags, w)
		}
	}

	return tags
}
