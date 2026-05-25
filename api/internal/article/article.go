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

// QuestionStatus は記事にぶら下がる質問返信（kind='question'）の解決状態を表す。
// solved > answered > unanswered の優先順位で、記事一覧クエリの CASE 式が一意に決定する。
// 質問返信が一つも無い通常記事は none。
type QuestionStatus string

const (
	QuestionStatusNone       QuestionStatus = "none"
	QuestionStatusUnanswered QuestionStatus = "unanswered"
	QuestionStatusAnswered   QuestionStatus = "answered"
	QuestionStatusSolved     QuestionStatus = "solved"
)

type Article struct {
	ID             int64
	Title          string
	Content        string
	UserID         int64
	Author         Author
	Tags           []Tag
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LikesCount     int64
	LikedByMe      bool
	QuestionStatus QuestionStatus
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
		normalizedTagName := strings.TrimSpace(tagName)
		if normalizedTagName == "" {
			return nil, ErrInvalidTagName
		}
		if utf8.RuneCountInString(normalizedTagName) > maxTagNameLength {
			return nil, ErrInvalidTagName
		}
		// 表示用に原文ケースを残す。重複判定は DB の UNIQUE(LOWER(name)) と揃えて
		// 大文字小文字を区別しない。
		dedupKey := strings.ToLower(normalizedTagName)
		if _, ok := seen[dedupKey]; ok {
			continue
		}

		seen[dedupKey] = struct{}{}
		normalizedTagNames = append(normalizedTagNames, normalizedTagName)
	}

	return normalizedTagNames, nil
}
