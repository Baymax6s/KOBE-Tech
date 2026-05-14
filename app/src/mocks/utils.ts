import { db } from "./db";
/**
 * =====================================================
 * MOCK UTILS
 * =====================================================
 */

/**
 * 現在時刻（ISO形式）
 */
export function now() {
  return new Date().toISOString();
}

/**
 * 認証ユーザー取得（mock_tokenベース）
 */
export function auth() {
  const token = window.localStorage.getItem("mock_token");

  if (!token) return null;

  const userId = Number(token.replace("mock-token-", ""));

  const user = db.users.find((u) => u.id === userId);

  return user ?? null;
}

/**
 * ページネーション処理
 */
export function paginate<T>(list: T[], page = 1, limit = 10) {
  const start = (page - 1) * limit;

  return {
    data: list.slice(start, start + limit),
    total: list.length,
    page,
    limit,
  };
}