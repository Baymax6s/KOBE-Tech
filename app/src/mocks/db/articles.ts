export const articles = [
  {
    id: 1,
    title: 'MSWを導入してみた',
    content: 'フロント開発がかなり楽になる',
    user_id: 1,
    likes_count: 3,
    tags: [{ id: 1, name: 'dev' }],
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
]
