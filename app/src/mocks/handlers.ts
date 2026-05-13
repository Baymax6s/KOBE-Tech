import { http, HttpResponse } from 'msw'

/**
 * =====================================================
 * MOCK DB（初期データ入り）
 * =====================================================
 */

const users = [
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

const tags = [
  { id: 1, name: 'dev' },
  { id: 2, name: 'frontend' },
  { id: 3, name: 'team' },
  { id: 4, name: 'vue' },
  { id: 5, name: 'typescript' },
]

const articles = [
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

const replies = [
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

const likes = new Set<string>()


/**
 * =====================================================
 * UTIL
 * =====================================================
 */

function auth() {
  const token = window.localStorage.getItem('mock_token')

  if (!token) return null

  const userId = Number(token.replace('mock-token-', ''))

  return users.find((u) => u.id === userId) ?? null
}

function now() {
  return new Date().toISOString()
}

function paginate<T>(list: T[], page = 1, limit = 10) {
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

  http.post('*/api/auth/login', async ({ request }) => {
    const body = (await request.json()) as {
      name: string
      password: string
    }

    const user = users.find(
      (u) => u.name === body.name && u.password === body.password,
    )

    if (!user) {
      return HttpResponse.json({ message: 'Invalid credentials' }, { status: 401 })
    }

    const token = `mock-token-${user.id}`

    // 👉 MSW内でブラウザlocalStorageに保存
    window.localStorage.setItem('mock_token', token)

    return HttpResponse.json({ token })
  }),

  http.get('*/api/auth/me', () => {
    const user = auth()

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    return HttpResponse.json({
      id: user.id,
      name: user.name,
      created_at: user.created_at,
      updated_at: user.updated_at,
    })
  }),

  /**
   * =========================
   * ARTICLES
   * =========================
   */

  http.get('*/api/articles', ({ request }) => {
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

  http.post('*/api/articles', async ({ request }) => {
    const user = auth()

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    const body = (await request.json()) as {
      title: string
      content: string
      tags?: string[]
    }

    const article = {
      id: articles.length + 1,
      title: body.title,
      content: body.content,
      user_id: user.id,
      likes_count: 0,
      tags: (body.tags || []).map((t, i) => ({
        id: tags.length + i + 1,
        name: t,
      })),
      created_at: now(),
      updated_at: now(),
    }

    articles.unshift(article)

    return HttpResponse.json(article, { status: 201 })
  }),

  http.get('*/api/articles/:id', ({ params }) => {
    const article = articles.find((a) => a.id === Number(params.id))

    if (!article) {
      return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    }

    const author = users.find((u) => u.id === article.user_id)

    return HttpResponse.json({
      ...article,
      author: author
        ? {
            id: author.id,
            name: author.name,
            created_at: author.created_at,
            updated_at: author.updated_at,
          }
        : null,
    })
  }),

  http.post('*/api/articles/:id/like', ({ params }) => {
    const user = auth()

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    const article = articles.find((a) => a.id === Number(params.id))

    if (!article) {
      return HttpResponse.json({ message: 'Not found' }, { status: 404 })
    }

    const key = `${user.id}-${article.id}`

    if (likes.has(key)) {
      return HttpResponse.json({ message: 'Already liked' }, { status: 409 })
    }

    likes.add(key)

    article.likes_count++

    return new HttpResponse(null, { status: 201 })
  }),

  /**
   * =========================
   * REPLIES
   * =========================
   */

  http.get('*/api/articles/:id/replies', ({ params }) => {
    return HttpResponse.json({
      replies: replies.filter((r) => r.article_id === Number(params.id)),
    })
  }),

  http.post('*/api/articles/:id/replies', async ({ request, params }) => {
    const user = auth()

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    const body = (await request.json()) as {
      body: string
      kind?: string
      parent_id?: number
    }

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
  }),

  /**
   * =========================
   * TAGS
   * =========================
   */

  http.get('*/api/tags', () => {
    const user = auth()

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    return HttpResponse.json({ tags })
  }),
]
