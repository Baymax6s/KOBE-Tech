type Reply = {
  id: number
  article_id: number
  body: string
  kind: string
  parent_id: number | null
  user_id: number
  user_name: string
  is_best?: boolean
  created_at: string
  updated_at: string
}

export const replies: Reply[] = [
  {
    id: 1,
    article_id: 1,
    body: 'これ便利すぎる',
    kind: 'comment',
    parent_id: null,
    user_id: 1,
    user_name: 'mock-user',
    is_best: false,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: 2,
    article_id: 1,
    body: '本番環境でも使ってる人は多いんですか？',
    kind: 'question',
    parent_id: null,
    user_id: 2,
    user_name: 'alice',
    is_best: false,
    created_at: '2026-05-02T09:00:00.000Z',
    updated_at: '2026-05-02T09:00:00.000Z',
  },
  {
    id: 3,
    article_id: 1,
    body: 'かなり増えてきてますよ。MSW の事例も増えてきてます。',
    kind: 'answer',
    parent_id: 2,
    user_id: 1,
    user_name: 'mock-user',
    is_best: true,
    created_at: '2026-05-02T11:00:00.000Z',
    updated_at: '2026-05-02T11:00:00.000Z',
  },
]
