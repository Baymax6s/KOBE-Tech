import { execFileSync } from "node:child_process"
import { extname } from "node:path"

const EDIT_TOOLS = new Set(["edit", "write", "multiedit", "apply_patch"])
const FORMATTABLE = new Set([".vue", ".ts", ".tsx", ".js", ".mjs", ".cjs", ".json", ".css", ".scss", ".md"])
const LINTABLE = new Set([".vue", ".ts", ".tsx", ".js", ".mjs", ".cjs"])

const pickFilePath = (args: Record<string, unknown> | undefined): string | undefined => {
  if (!args) return undefined
  for (const key of ["filePath", "file_path", "path"]) {
    const value = args[key]
    if (typeof value === "string") return value
  }
  return undefined
}

const GENERATED_BLOCK_MESSAGE =
  "Refusing to edit generated API client.\n" +
  "`app/src/api/generated/` は swagger-typescript-api が自動生成しています。\n" +
  "Go 側を編集し、`make swagger` → `npm --prefix app run generate:api` で再生成してください。"

export const KobeTechHooks = async (ctx: { project?: { worktree?: string } }) => {
  const root = ctx.project?.worktree ?? process.cwd()
  const appDir = `${root}/app`

  return {
    "tool.execute.before": async (input: { tool: string }, output: { args?: Record<string, unknown> }) => {
      if (!EDIT_TOOLS.has(input.tool)) return
      const filePath = pickFilePath(output.args)
      if (!filePath) return
      if (filePath.includes("/app/src/api/generated/")) {
        throw new Error(GENERATED_BLOCK_MESSAGE)
      }
    },

    "tool.execute.after": async (input: { tool: string }, output: { args?: Record<string, unknown> }) => {
      if (!EDIT_TOOLS.has(input.tool)) return
      const filePath = pickFilePath(output.args)
      if (!filePath) return
      if (!filePath.includes("/app/")) return
      if (filePath.includes("/app/src/api/generated/")) return

      const ext = extname(filePath)
      try {
        if (FORMATTABLE.has(ext)) {
          execFileSync("npx", ["--no-install", "prettier", "--write", "--log-level=warn", filePath], {
            cwd: appDir,
            stdio: "inherit",
          })
        }
        if (LINTABLE.has(ext)) {
          execFileSync("npx", ["--no-install", "eslint", "--fix", filePath], {
            cwd: appDir,
            stdio: "inherit",
          })
        }
      } catch {
        // フォーマッタ失敗で編集を巻き戻さない
      }
    },
  }
}
