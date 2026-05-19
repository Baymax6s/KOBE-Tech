package article

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrInvalidArticle = errors.New("title and content are required")
	ErrInvalidTagName = errors.New("tags must contain tag names between 1 and 10 characters")
)

const maxTagNameLength = 10

type Author struct {
	ID   int64
	Name string
}

type Article struct {
	ID         int64
	Title      string
	Content    string
	UserID     int64
	Author     Author
	Tags       []Tag
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LikesCount int64
	LikedByMe  bool
}

func NormalizeCreateInput(title, content string, tagNames []string) (string, string, []string, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" || content == "" {
		return "", "", nil, ErrInvalidArticle
	}

	normalizedTagNames, err := NormalizeTagNames(tagNames)
	if err != nil {
		return "", "", nil, err
	}

	return title, content, normalizedTagNames, nil
}

func NormalizeTagNames(tagNames []string) ([]string, error) {
	normalizedTagNames := make([]string, 0, len(tagNames))
	seen := make(map[string]struct{}, len(tagNames))

	for _, tagName := range tagNames {
		normalizedTagName := strings.ToLower(strings.TrimSpace(tagName))
		if normalizedTagName == "" {
			return nil, ErrInvalidTagName
		}
		if utf8.RuneCountInString(normalizedTagName) > maxTagNameLength {
			return nil, ErrInvalidTagName
		}
		if _, ok := seen[normalizedTagName]; ok {
			continue
		}

		seen[normalizedTagName] = struct{}{}
		normalizedTagNames = append(normalizedTagNames, normalizedTagName)
	}

	return normalizedTagNames, nil
}
