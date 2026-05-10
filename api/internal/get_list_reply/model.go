package reply

import "time"

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
