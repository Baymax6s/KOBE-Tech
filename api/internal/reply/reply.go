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

type Kind int16

const (
	KindComment  Kind = 0
	KindQuestion Kind = 1
	KindAnswer   Kind = 2
)

func (k Kind) String() string {
	switch k {
	case KindComment:
		return "comment"
	case KindQuestion:
		return "question"
	case KindAnswer:
		return "answer"
	default:
		return ""
	}
}

func ParseKind(s string) (Kind, bool) {
	switch s {
	case "comment":
		return KindComment, true
	case "question":
		return KindQuestion, true
	case "answer":
		return KindAnswer, true
	default:
		return 0, false
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
		return "", 0, ErrInvalidBody
	}
	if utf8.RuneCountInString(body) > maxBodyLength {
		return "", 0, ErrBodyTooLong
	}

	// kind 省略時は comment 扱い。値が指定された場合は comment / question / answer のいずれかを受け付ける。
	// 親 reply との整合性（例: ルートに answer は不可）は repository 側で検証する。
	kind := KindComment
	if kindValue != nil {
		parsed, ok := ParseKind(*kindValue)
		if !ok {
			return "", 0, ErrInvalidKind
		}
		kind = parsed
	}

	return body, kind, nil
}
