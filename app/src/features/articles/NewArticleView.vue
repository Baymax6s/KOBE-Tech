<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { MdEditor, type ToolbarNames } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { api } from '@/api/client'
import { useArticleNotificationStore } from '@/stores/articleNotification'
import MarkdownContent from './MarkdownContent.vue'

defineOptions({
  name: 'NewArticleView',
})

type BodyTab = 'edit' | 'preview'

type ArticleFormRef = {
  validate: () => Promise<{ valid: boolean }>
  reset: () => void
}

const formRef = ref<ArticleFormRef | null>(null)

const form = reactive({
  title: '',
  tags: [] as string[],
  body: '',
})

const router = useRouter()
const notificationStore = useArticleNotificationStore()

const submitting = ref(false)
const submitError = ref<string | null>(null)
const tagCandidates = ref<string[]>([])
const bodyTab = ref<BodyTab>('edit')

// プレビューは既存の MarkdownContent (記事詳細と同じレンダラ) で表示するため、
// MdEditor 側のプレビュー / 全画面 / カタログ系ツールバーは除外する。
const editorToolbars: ToolbarNames[] = [
  'bold',
  'underline',
  'italic',
  'strikeThrough',
  '-',
  'title',
  'sub',
  'sup',
  'quote',
  'unorderedList',
  'orderedList',
  'task',
  '-',
  'codeRow',
  'code',
  'link',
  'table',
  '-',
  'revoke',
  'next',
]

const canSubmit = computed(
  () => !submitting.value && !!form.title.trim() && !!form.body.trim(),
)

const submit = async () => {
  if (submitting.value) return
  if (!formRef.value) return

  const { valid } = await formRef.value.validate()
  if (!valid) return

  submitting.value = true
  submitError.value = null

  try {
    await api.api.articlesCreate({
      title: form.title,
      content: form.body,
      tags: form.tags,
    })

    notificationStore.markCreated()

    await router.push('/articles')
  } catch {
    submitError.value = '投稿に失敗しました。もう一度お試しください。'
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    const response = await api.api.tagsList()
    tagCandidates.value = response.data.tags?.map((tag) => tag.name) ?? []
  } catch {
    tagCandidates.value = []
  }
})
</script>

<template>
  <v-sheet color="grey-lighten-4" class="page-bg">
    <v-container max-width="1200" class="editor-wrap">
      <v-form ref="formRef" class="article-form" @submit.prevent="submit">
        <div class="d-flex justify-end mb-3">
          <v-btn
            type="submit"
            color="black"
            class="px-6"
            :loading="submitting"
            :disabled="!canSubmit"
          >
            投稿
          </v-btn>
        </div>

        <v-alert
          v-if="submitError"
          type="error"
          class="mb-3"
          closable
          @click:close="submitError = null"
        >
          {{ submitError }}
        </v-alert>

        <v-sheet rounded class="editor-card pa-4 pa-md-6">
          <v-text-field
            v-model="form.title"
            placeholder="タイトルを入力"
            aria-label="タイトル"
            variant="plain"
            density="comfortable"
            maxlength="200"
            :rules="[(v) => !!v || 'タイトルは必須です']"
            validate-on="input"
            class="title-input"
          />

          <v-divider />

          <v-combobox
            v-model="form.tags"
            placeholder="タグを入力してEnter"
            aria-label="タグ"
            variant="plain"
            density="comfortable"
            :items="tagCandidates"
            multiple
            chips
            closable-chips
            prepend-inner-icon="mdi-tag-outline"
          />

          <v-divider class="mb-2" />

          <v-tabs
            v-model="bodyTab"
            density="comfortable"
            color="primary"
            class="mb-2"
          >
            <v-tab value="edit">編集</v-tab>
            <v-tab value="preview">プレビュー</v-tab>
          </v-tabs>

          <!-- v-window のスライドアニメは不要、かつ編集中のカーソル位置・履歴を保つため v-show で切替 -->
          <div class="body-area">
            <MdEditor
              v-show="bodyTab === 'edit'"
              v-model="form.body"
              language="en-US"
              theme="light"
              :toolbars="editorToolbars"
              :preview="false"
              :show-code-row-number="false"
              placeholder="本文を入力（Markdown）"
              class="body-editor"
            />

            <div v-show="bodyTab === 'preview'" class="body-preview pa-2">
              <MarkdownContent
                v-if="bodyTab === 'preview' && form.body.trim()"
                :source="form.body"
              />
              <div
                v-else-if="bodyTab === 'preview'"
                class="text-medium-emphasis"
              >
                本文がまだ入力されていません。
              </div>
            </div>
          </div>
        </v-sheet>
      </v-form>
    </v-container>
  </v-sheet>
</template>

<style scoped>
/* Qiita スタイル: ページ自体はスクロールせず、本文エリアの内部だけがスクロールする。
   そのために v-main → page-bg → editor-wrap → article-form → editor-card → body-area
   と高さ 100% / flex column を連鎖させる必要がある。 */
.page-bg {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-wrap {
  flex: 1 1 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

/* v-form / editor-card 直下の要素は基本「自然な高さ」で積み、
   一番下の領域 (.editor-card, .body-area) だけが残り高さを取る。 */
.article-form,
.editor-card,
.body-area {
  flex: 1 1 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
}
.article-form > *,
.editor-card > * {
  flex: 0 0 auto;
}
.article-form > .editor-card,
.editor-card > .body-area {
  flex: 1 1 0;
  min-height: 0;
}

/* タイトル入力は記事のタイトルらしい大きさに */
.title-input :deep(input) {
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1.3;
}

/* MdEditor は内部でツールバー + エディタの flex レイアウトを完結させているので、
   外側に 100% の高さを渡すだけで残り領域いっぱいに広がる。 */
.body-editor {
  height: 100%;
}

.body-preview {
  flex: 1 1 0;
  min-height: 0;
  overflow-y: auto;
}
</style>
