package product

import (
	"regexp"
	"strings"
)

var (
	nonAlnum  = regexp.MustCompile(`[^a-z0-9\s-]`)
	multiDash = regexp.MustCompile(`-+`)
	space     = regexp.MustCompile(`\s+`)
)

func Slugify(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)

	// remove special characters
	text = nonAlnum.ReplaceAllString(text, "")

	// spaces → hyphens
	text = space.ReplaceAllString(text, "-")

	// collapse multiple hyphens
	text = multiDash.ReplaceAllString(text, "-")

	return text
}
