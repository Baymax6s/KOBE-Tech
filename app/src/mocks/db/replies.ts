export type Reply = {
  id: number
  article_id: number
  body: string
  kind: string
  parent_id: number | null
  user_id: number
  user_name: string
  is_best: boolean
  created_at: string
  updated_at: string
}

export const replies: Reply[] = [
  // ===== Article 1 (by mock-user, id=1) =====
  {
    id: 1,
    article_id: 1,
    body: 'これ便利すぎる',
    kind: 'comment',
    parent_id: null,
    user_id: 2,
    user_name: 'alice',
    is_best: false,
    created_at: '2026-05-01T10:00:00.000Z',
    updated_at: '2026-05-01T10:00:00.000Z',
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
    body: '導入事例としてはかなり増えてきてます。MSW の公式サイトにも事例が載ってますよ。',
    kind: 'answer',
    parent_id: 2,
    user_id: 1,
    user_name: 'mock-user',
    is_best: true,
    created_at: '2026-05-02T11:00:00.000Z',
    updated_at: '2026-05-02T11:00:00.000Z',
  },
  {
    id: 4,
    article_id: 1,
    body: '参考になります！',
    kind: 'comment',
    parent_id: 1,
    user_id: 1,
    user_name: 'mock-user',
    is_best: false,
    created_at: '2026-05-03T08:00:00.000Z',
    updated_at: '2026-05-03T08:00:00.000Z',
  },

  // ===== Article 2 (by alice, id=2) =====
  {
    id: 5,
    article_id: 2,
    body: 'Composition API と Options API はどちらを使うべきですか？',
    kind: 'question',
    parent_id: null,
    user_id: 1,
    user_name: 'mock-user',
    is_best: false,
    created_at: '2026-05-05T14:00:00.000Z',
    updated_at: '2026-05-05T14:00:00.000Z',
  },
  {
    id: 6,
    article_id: 2,
    body: 'プロジェクトの規模次第ですが、新規なら Composition API 一択でいいと思います。型推論が効くので TS との相性も良いです。',
    kind: 'answer',
    parent_id: 5,
    user_id: 2,
    user_name: 'alice',
    is_best: false,
    created_at: '2026-05-05T16:00:00.000Z',
    updated_at: '2026-05-05T16:00:00.000Z',
  },
]
