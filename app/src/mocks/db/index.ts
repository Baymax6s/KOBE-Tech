import { users } from './users';
import { tags } from './tags';
import { articles } from './articles';
import { replies } from './replies';

export const db = {
  users,
  tags,
  articles,
  replies,
  likes: new Set<string>(),
};