package reply

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const maxBodyLength = 2000

var (
	ErrInvalidBody = errors.New("body is required")
	ErrBodyTooLong = fmt.Errorf("body must be %d characters or less", maxBodyLength)
	ErrInvalidKind = errors.New("kind must be one of 'comment', 'question', 'answer'")
)

type Kind string

const (
	KindComment  Kind = "comment"
	KindQuestion Kind = "question"
	KindAnswer   Kind = "answer"
)

func ParseKind(s string) (Kind, bool) {
	switch Kind(s) {
	case KindComment, KindQuestion, KindAnswer:
		return Kind(s), true
	default:
		return "", false
	}
}

type Reply struct {
	ID        int64
	ArticleID int64
	ParentID  *int64
	Kind      Kind
	Body      string
	UserID    int64
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NormalizeCreateInput(body string, kindValue *string) (string, Kind, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return "", "", ErrInvalidBody
	}
	if utf8.RuneCountInString(body) > maxBodyLength {
		return "", "", ErrBodyTooLong
	}

	// kind 省略時は comment 扱い。値が指定された場合は comment / question / answer のいずれかを受け付ける。
	// 親 reply との整合性（例: ルートに answer は不可）は repository 側で検証する。
	kind := KindComment
	if kindValue != nil {
		parsed, ok := ParseKind(*kindValue)
		if !ok {
			return "", "", ErrInvalidKind
		}
		kind = parsed
	}

	return body, kind, nil
}
