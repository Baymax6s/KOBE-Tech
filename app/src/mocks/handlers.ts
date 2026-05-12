/* eslint-disable */
/* tslint:disable */
// @ts-nocheck

import { http, HttpResponse } from 'msw'

/**
 * =====================================================
 * MOCK DB（初期データ入り）
 * =====================================================
 */

let users = [
  {
    id: 1,
    name: 'mock-user',
    password: 'password',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: 2,
    name: 'alice',
    password: 'password',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
]

let sessions = new Map<string, number>()

const MOCK_TOKEN = 'mock-jwt-token'

let tags = [
  { id: 1, name: 'dev' },
  { id: 2, name: 'frontend' },
  { id: 3, name: 'team' },
  { id: 4, name: 'vue' },
  { id: 5, name: 'typescript' },
]

let articles = [
  {
    id: 1,
    title: 'MSWを導入してみた',
    content: 'フロント開発がかなり楽になる',
    user_id: 1,
    likes_count: 3,
    tags: [
      { id: 1, name: 'dev' },
      { id: 2, name: 'frontend' },
    ],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: 2,
    title: 'Vue3 + TS構成まとめ',
    content: 'Composition APIの整理が重要',
    user_id: 2,
    likes_count: 7,
    tags: [{ id: 5, name: 'typescript' }],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: 3,
    title: 'チーム開発の罠',
    content: 'APIモックずれが一番危ない',
    user_id: 1,
    likes_count: 12,
    tags: [{ id: 3, name: 'team' }],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
]

let replies = [
  {
    id: 1,
    article_id: 1,
    body: 'これ便利すぎる',
    kind: 'comment',
    parent_id: undefined,
    user_id: 1,
    user_name: 'mock-user',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: 2,
    article_id: 1,
    body: '確かに開発早くなる',
    kind: 'comment',
    parent_id: 1,
    user_id: 2,
    user_name: 'alice',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
]

/**
 * =====================================================
 * UTIL
 * =====================================================
 */

function auth(request: Request) {
  const token = request.headers.get('Authorization')?.replace('Bearer ', '')

  if (token !== MOCK_TOKEN) return null

  return users.find((u) => u.id === 1)
}

function now() {
  return new Date().toISOString()
}

function paginate(list: any[], page = 1, limit = 10) {
  const start = (page - 1) * limit
  return {
    data: list.slice(start, start + limit),
    total: list.length,
    page,
    limit,
  }
}

/**
 * =====================================================
 * HANDLERS
 * =====================================================
 */

export const handlers = [
  /**
   * =========================
   * AUTH
   * =========================
   */

  http.post('http://localhost:8080/api/auth/login', async ({ request }) => {
    const body = await request.json()

    const user = users.find(
      (u) => u.name === body.name && u.password === body.password,
    )

    if (!user) {
      return HttpResponse.json(
        { message: 'Invalid credentials' },
        { status: 401 },
      )
    }

    sessions.set(MOCK_TOKEN, user.id)

    return HttpResponse.json({
      token: MOCK_TOKEN,
    })
  }),

  http.get('http://localhost:8080/api/auth/me', ({ request }) => {
    const user = auth(request)

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    return HttpResponse.json(user)
  }),

  /**
   * =========================
   * ARTICLES
   * =========================
   */

  http.get('http://localhost:8080/api/articles', ({ request }) => {
    const url = new URL(request.url)
    const page = Number(url.searchParams.get('page') || 1)
    const limit = Number(url.searchParams.get('limit') || 10)
    const q = url.searchParams.get('q')

    let result = articles

    if (q) {
      result = result.filter((a) =>
        a.title.toLowerCase().includes(q.toLowerCase()),
      )
    }

    const paginated = paginate(result, page, limit)

    return HttpResponse.json({
      articles: paginated.data,
      total: paginated.total,
      page: paginated.page,
      limit: paginated.limit,
    })
  }),

  http.post('http://localhost:8080/api/articles', async ({ request }) => {
    const user = auth(request)

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    const body = await request.json()

    const article = {
      id: articles.length + 1,
      title: body.title,
      content: body.content,
      user_id: user.id,
      likes_count: 0,
      tags: (body.tags || []).map((t, i) => ({
        id: i + 1,
        name: t,
      })),
      created_at: now(),
      updated_at: now(),
    }

    articles.unshift(article)

    return HttpResponse.json(article, { status: 201 })
  }),

  http.get('http://localhost:8080/api/articles/:id', ({ params }) => {
    const article = articles.find((a) => a.id === Number(params.id))

    if (!article) {
      return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    }

    return HttpResponse.json({
      ...article,
      author: users.find((u) => u.id === article.user_id),
    })
  }),

  http.post(
    'http://localhost:8080/api/articles/:id/like',
    ({ request, params }) => {
      const user = auth(request)
      if (!user) {
        return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
      }

      const article = articles.find((a) => a.id === Number(params.id))

      if (!article) {
        return HttpResponse.json({ message: 'Not found' }, { status: 404 })
      }

      article.likes_count++

      return new HttpResponse(null, { status: 204 })
    },
  ),

  /**
   * =========================
   * REPLIES
   * =========================
   */

  http.get('http://localhost:8080/api/articles/:id/replies', ({ params }) => {
    return HttpResponse.json({
      replies: replies.filter((r) => r.article_id === Number(params.id)),
    })
  }),

  http.post(
    'http://localhost:8080/api/articles/:id/replies',
    async ({ request, params }) => {
      const user = auth(request)

      if (!user) {
        return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
      }

      const body = await request.json()

      const reply = {
        id: replies.length + 1,
        article_id: Number(params.id),
        body: body.body,
        kind: body.kind || 'comment',
        parent_id: body.parent_id,
        user_id: user.id,
        user_name: user.name,
        created_at: now(),
        updated_at: now(),
      }

      replies.push(reply)

      return HttpResponse.json(reply, { status: 201 })
    },
  ),

  /**
   * =========================
   * TAGS
   * =========================
   */

  http.get('http://localhost:8080/api/tags', ({ request }) => {
    const user = auth(request)

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    return HttpResponse.json({ tags })
  }),
]
