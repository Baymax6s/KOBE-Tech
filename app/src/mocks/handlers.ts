import { profileHandlers } from "./handlers/profiles";
import { articleHandlers } from "./handlers/articles";
import { replyHandlers } from "./handlers/replies";
import { tagHandlers } from "./handlers/tags";

export const handlers = [
  ...profileHandlers,
  ...articleHandlers,
  ...replyHandlers,
  ...tagHandlers,
];