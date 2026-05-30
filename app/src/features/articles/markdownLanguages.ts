// SupportedLanguagesDialog で「対応言語」として案内する言語の一覧。
// 実際のシンタックスハイライトは md-editor-v3 が CDN から読み込む
// highlight.js が担うため、ここは表示用のメタデータに徹する。
// 追加・削除すればダイアログの表示もそのまま追随する。
export const SUPPORTED_LANGUAGE_NAMES = [
  'bash',
  'c',
  'clojure',
  'cmake',
  'cpp',
  'csharp',
  'css',
  'dart',
  'diff',
  'dockerfile',
  'elixir',
  'erlang',
  'fsharp',
  'go',
  'graphql',
  'groovy',
  'haskell',
  'ini',
  'java',
  'javascript',
  'json',
  'kotlin',
  'latex',
  'less',
  'lua',
  'makefile',
  'markdown',
  'nginx',
  'objectivec',
  'ocaml',
  'perl',
  'php',
  'plaintext',
  'powershell',
  'properties',
  'protobuf',
  'python',
  'r',
  'ruby',
  'rust',
  'scala',
  'scss',
  'shell',
  'sql',
  'swift',
  'typescript',
  'vbnet',
  'xml',
  'yaml',
]

// フェンスのよく使われる短縮表記を canonical な言語名へ寄せて表示する。
export const LANGUAGE_ALIASES: Record<string, string> = {
  js: 'javascript',
  jsx: 'javascript',
  ts: 'typescript',
  tsx: 'typescript',
  py: 'python',
  sh: 'bash',
  zsh: 'bash',
  yml: 'yaml',
  md: 'markdown',
  rb: 'ruby',
  kt: 'kotlin',
  cs: 'csharp',
  golang: 'go',
  html: 'xml',
  vue: 'xml',
  text: 'plaintext',
  txt: 'plaintext',
  toml: 'ini',
}

// 「canonical 名」と「その alias 群」の組で公開する。
export type SupportedLanguage = {
  name: string
  aliases: string[]
}

const aliasesByCanonical = Object.entries(LANGUAGE_ALIASES).reduce<
  Record<string, string[]>
>((acc, [alias, canonical]) => {
  if (!acc[canonical]) acc[canonical] = []
  acc[canonical].push(alias)
  return acc
}, {})

export const SUPPORTED_LANGUAGES: SupportedLanguage[] = [
  ...SUPPORTED_LANGUAGE_NAMES,
]
  .sort()
  .map((name) => ({
    name,
    aliases: (aliasesByCanonical[name] ?? []).sort(),
  }))
