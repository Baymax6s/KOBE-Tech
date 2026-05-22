<script setup lang="ts">
defineOptions({ name: 'MarkdownCheatsheet' })

type CheatItem = {
  label: string
  syntax: string
}

const items: CheatItem[] = [
  { label: '見出し', syntax: '## 見出し2\n### 見出し3\n#### 見出し4' },
  { label: '強調', syntax: '**強調**' },
  { label: 'インラインコード', syntax: '`インラインコード`' },
  { label: 'コードブロック', syntax: '```ts\nconst x = 1\n```' },
  { label: 'リンク', syntax: '[リンクテキスト](http://...)' },
  { label: 'リスト', syntax: '- リスト1\n- リスト2\n- リスト3' },
  {
    label: '番号付きリスト',
    syntax: '1. 番号付きリスト1\n2. 番号付きリスト2',
  },
  { label: '引用', syntax: '> 引用文' },
]
</script>

<template>
  <v-card variant="flat" color="transparent">
    <v-card-title class="text-subtitle-1 font-weight-bold px-0">
      文章を見やすく装飾する方法
    </v-card-title>

    <v-row dense>
      <v-col v-for="item in items" :key="item.label" cols="12" sm="6" md="4">
        <div class="text-body-2 font-weight-medium mb-1">
          {{ item.label }}
        </div>
        <v-sheet color="grey-lighten-4" rounded class="pa-3">
          <pre class="cheatsheet-pre text-body-2">{{ item.syntax }}</pre>
        </v-sheet>
      </v-col>
    </v-row>
  </v-card>
</template>

<style scoped>
.cheatsheet-pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family:
    'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
}
</style>
