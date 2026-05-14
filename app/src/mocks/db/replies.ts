type Reply = {
  id: number
  article_id: number
  body: string
  kind: string
  parent_id: number | null
  user_id: number
  user_name: string
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
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
]
