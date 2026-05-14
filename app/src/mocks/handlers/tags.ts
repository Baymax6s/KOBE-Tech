import { http, HttpResponse } from "msw";
import { db } from "../db/index";
import { auth } from "../utils";

export const tagHandlers = [
  http.get("*/api/tags", () => {
    const user = auth();

    if (!user) {
      return HttpResponse.json(
        { message: "Unauthorized" },
        { status: 401 }
      );
    }

    return HttpResponse.json({
      tags: db.tags,
    });
  }),
];