package Slug

import "strings"

func Create(title string) string {
	return strings.Join(strings.Split(title, " "), "-")
}
