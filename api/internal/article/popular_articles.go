package article

// PopularArticle はタグ内ランキングの 1 記事分。Rank はタグ内のいいね数順位（RANK() の値）。
type PopularArticle struct {
	ArticleID  int64
	Title      string
	LikesCount int64
	Rank       int64
}

// TagRanking は 1 タグと、そのタグ内の上位記事を束ねたもの。
type TagRanking struct {
	TagID    int64
	TagName  string
	Articles []PopularArticle
}
