import { http, HttpResponse } from 'msw'
import { db } from '../db'
import { auth, now } from '../utils'
import type { Reply } from '../db/replies'

function findReply(id: number): Reply | undefined {
  return db.replies.find((r) => r.id === id)
}

export const replyHandlers = [
  http.get('*/api/articles/:id/replies', ({ params }) => {
    return HttpResponse.json({
      replies: db.replies.filter((r) => r.article_id === Number(params.id)),
    })
  }),

  http.post('*/api/articles/:id/replies', async ({ request, params }) => {
    const user = auth()

    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    const body = (await request.json()) as {
      body: string
      kind?: string
      parent_id?: number
    }

    const reply: Reply = {
      id: db.replies.length + 1,
      article_id: Number(params.id),
      body: body.body,
      kind: body.kind || 'comment',
      parent_id: body.parent_id ?? null,
      user_id: user.id,
      user_name: user.name,
      is_best: false,
      created_at: now(),
      updated_at: now(),
    }

    db.replies.push(reply)

    return HttpResponse.json(reply, { status: 201 })
  }),

  http.post('*/api/replies/:reply_id/best', async ({ params }) => {
    const user = auth()
    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    const replyId = Number(params.reply_id)
    const reply = findReply(replyId)
    if (!reply) {
      return HttpResponse.json({ message: 'reply not found' }, { status: 404 })
    }

    if (reply.kind !== 'answer') {
      return HttpResponse.json(
        { message: 'reply is not an answer' },
        { status: 400 },
      )
    }

    if (reply.parent_id == null) {
      return HttpResponse.json(
        { message: 'reply is not an answer' },
        { status: 400 },
      )
    }

    const parent = findReply(reply.parent_id)
    if (!parent || parent.kind !== 'question') {
      return HttpResponse.json(
        { message: 'reply is not an answer' },
        { status: 400 },
      )
    }

    if (parent.user_id !== user.id) {
      return HttpResponse.json(
        { message: 'only the question author can mark a best answer' },
        { status: 403 },
      )
    }

    if (reply.is_best) {
      return HttpResponse.json(
        { message: 'best answer already exists for this question' },
        { status: 409 },
      )
    }

    reply.is_best = true

    return HttpResponse.json({ reply_id: replyId, is_best: true })
  }),

  http.put('*/api/replies/:reply_id/best', async ({ params }) => {
    const user = auth()
    if (!user) {
      return HttpResponse.json({ message: 'Unauthorized' }, { status: 401 })
    }

    const replyId = Number(params.reply_id)
    const reply = findReply(replyId)
    if (!reply) {
      return HttpResponse.json({ message: 'reply not found' }, { status: 404 })
    }

    if (reply.kind !== 'answer') {
      return HttpResponse.json(
        { message: 'reply is not an answer' },
        { status: 400 },
      )
    }

    if (!reply.is_best) {
      return HttpResponse.json(
        { message: 'reply is not marked as best answer' },
        { status: 409 },
      )
    }

    if (reply.parent_id == null) {
      return HttpResponse.json(
        { message: 'reply is not an answer' },
        { status: 400 },
      )
    }

    const parent = findReply(reply.parent_id)
    if (!parent || parent.kind !== 'question') {
      return HttpResponse.json(
        { message: 'reply is not an answer' },
        { status: 400 },
      )
    }

    if (parent.user_id !== user.id) {
      return HttpResponse.json(
        { message: 'only the question author can mark a best answer' },
        { status: 403 },
      )
    }

    reply.is_best = false

    return HttpResponse.json({ reply_id: replyId, is_best: false })
  }),
]
