package handler

import "github.com/Baymax6s/KOBE-Tech/api/internal/article"

func newListArticlesJSONResponse(articles []article.Article) ListArticlesJSONResponse {
	response := ListArticlesJSONResponse{
		Articles: make([]ArticleListItemJSON, 0, len(articles)),
	}

	for _, item := range articles {
		response.Articles = append(response.Articles, ArticleListItemJSON{
			ID:         item.ID,
			Title:      item.Title,
			Content:    item.Content,
			UserID:     item.UserID,
			Tags:       newArticleTagJSONs(item.Tags),
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			LikesCount: item.LikesCount,
		})
	}

	return response
}

func newArticleTagJSONs(tags []article.Tag) []ArticleTagJSON {
	response := make([]ArticleTagJSON, 0, len(tags))
	for _, tag := range tags {
		response = append(response, ArticleTagJSON{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return response
}

func newListTagsJSONResponse(tags []article.Tag) ListTagsJSONResponse {
	response := ListTagsJSONResponse{
		Tags: make([]TagJSON, 0, len(tags)),
	}

	for _, tag := range tags {
		response.Tags = append(response.Tags, TagJSON{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return response
}
