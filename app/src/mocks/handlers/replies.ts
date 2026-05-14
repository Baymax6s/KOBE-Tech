import { http, HttpResponse } from "msw";
import { db } from "../db";
import { auth,now} from "../utils";

export const replyHandlers = [
  http.get("*/api/articles/:id/replies", ({ params }) => {
    return HttpResponse.json({
      replies: db.replies.filter(r => r.article_id === Number(params.id)),
    });
  }),

  http.post("*/api/articles/:id/replies", async ({ request, params }) => {
    const user = auth();

    if (!user) {
      return HttpResponse.json({ message: "Unauthorized" }, { status: 401 });
    }

    const body = (await request.json()) as {
      body: string;
      kind?: string;
      parent_id?: number;
    };

    const reply = {
      id: db.replies.length + 1,
      article_id: Number(params.id),
      body: body.body,
      kind: body.kind || "comment",
      parent_id: body.parent_id,
      user_id: user.id,
      user_name: user.name,
      created_at: now(),
      updated_at: now(),
    };

    db.replies.push(reply);

    return HttpResponse.json(reply, { status: 201 });
  }),
];