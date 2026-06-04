# 発表スライド（Marp）

`presentation.md` は [Marp](https://marp.app/) 形式の発表スライドです（Markdown で記述）。
構成・各スライドの意図は [docs/demo/README.md](../demo/README.md) と対になっています。

## プレビュー

### VS Code（おすすめ）

拡張機能 **「Marp for VS Code」** を入れて `presentation.md` を開き、右上のプレビューを表示。

### CLI（PDF / HTML / PPTX に書き出し）

```bash
# プレビュー（ライブリロード）
npx @marp-team/marp-cli@latest -p docs/slides/presentation.md

# 書き出し
npx @marp-team/marp-cli@latest docs/slides/presentation.md -o presentation.pdf
npx @marp-team/marp-cli@latest docs/slides/presentation.md -o presentation.html
npx @marp-team/marp-cli@latest docs/slides/presentation.md --pptx -o presentation.pptx
```

## 構成（15分）

| 区間 | スライド | 目安 |
| --- | --- | --- |
| 導入＋チーム紹介 | タイトル / チーム | 1分 |
| 製品 | 命題 / 4本柱目次 / 柱1〜4 / 通しデモ | 11分 |
| 作業プロセス | 回し方 / ベロシティ / 成長 | 2.5分 |
| まとめ | クロージング | 0.5分 |

- 各「柱」スライドの <span>▶ デモ</span> 箇所で、対応するデモ映像（`docs/demo/` の手順で録画）を無音ループ再生する想定。
- スピーカーノート（各スライド末尾の HTML コメント）に進行のコツを記載。Marp のプレゼンターモードで表示される。
