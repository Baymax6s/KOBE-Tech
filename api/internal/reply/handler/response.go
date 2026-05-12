package handler

import "github.com/Baymax6s/KOBE-Tech/api/internal/reply"

func newReplyJSONs(replies []reply.Reply) []ReplyJSON {
	items := make([]ReplyJSON, 0, len(replies))
	for _, item := range replies {
		items = append(items, newReplyJSON(item))
	}
	return items
}

func newReplyJSON(item reply.Reply) ReplyJSON {
	return ReplyJSON{
		ID:        item.ID,
		ArticleID: item.ArticleID,
		ParentID:  item.ParentID,
		Kind:      item.Kind.String(),
		Body:      item.Body,
		UserID:    item.UserID,
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
