import { http, HttpResponse } from "msw";
import { db } from "../db";
import { auth, now, paginate } from "../utils";

export const articleHandlers = [
  http.get("*/api/articles", ({ request }) => {
    const url = new URL(request.url);
    const page = Number(url.searchParams.get("page") || 1);
    const limit = Number(url.searchParams.get("limit") || 10);

    const result = paginate(db.articles, page, limit);

    return HttpResponse.json({
      articles: result.data,
      total: result.total,
      page,
      limit,
    });
  }),

  http.get("*/api/articles/:id", ({ params }) => {
    const article = db.articles.find(a => a.id === Number(params.id));

    if (!article) {
      return HttpResponse.json({ message: "Not found" }, { status: 404 });
    }

    return HttpResponse.json(article);
  }),

  http.post("*/api/articles", async ({ request }) => {
  const user = auth();

  if (!user) {
    return HttpResponse.json({ message: "Unauthorized" }, { status: 401 });
  }

  const body = (await request.json()) as {
    title: string;
    content: string;
    tags?: string[];
  };

  const article = {
    id: db.articles.length + 1,
    title: body.title,
    content: body.content,
    user_id: user.id,
    likes_count: 0,

    // ★ここが修正ポイント
    tags: (body.tags || []).map(name => {
      const tag = db.tags.find(t => t.name === name);
      return tag ?? { id: -1, name };
    }),

    created_at: now(),
    updated_at: now(),
  };

  db.articles.unshift(article);

  return HttpResponse.json(article, { status: 201 });
}),

  http.post("*/api/articles/:id/like", ({ params }) => {
    const user = auth();

    if (!user) {
      return HttpResponse.json({ message: "Unauthorized" }, { status: 401 });
    }

    const article = db.articles.find(a => a.id === Number(params.id));

    if (!article) {
      return HttpResponse.json({ message: "Not found" }, { status: 404 });
    }

    const key = `${user.id}-${article.id}`;

    if (db.likes.has(key)) {
      return HttpResponse.json({ message: "Already liked" }, { status: 409 });
    }

    db.likes.add(key);
    article.likes_count++;

    return new HttpResponse(null, { status: 201 });
  }),
];