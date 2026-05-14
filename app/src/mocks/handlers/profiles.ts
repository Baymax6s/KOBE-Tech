import { http, HttpResponse } from "msw";
import { db } from "../db";
import { auth } from "../utils";

export const profileHandlers = [
  http.post("*/api/auth/login", async ({ request }) => {
    const body = (await request.json()) as {
      name: string;
      password: string;
    };

    const user = db.users.find(
      u => u.name === body.name && u.password === body.password
    );

    if (!user) {
      return HttpResponse.json(
        { message: "Invalid credentials" },
        { status: 401 }
      );
    }

    const token = `mock-token-${user.id}`;
    window.localStorage.setItem("mock_token", token);

    return HttpResponse.json({ token });
  }),

  http.get("*/api/auth/me", () => {
    const user = auth();

    if (!user) {
      return HttpResponse.json(
        { message: "Unauthorized" },
        { status: 401 }
      );
    }

    return HttpResponse.json(user);
  }),
];