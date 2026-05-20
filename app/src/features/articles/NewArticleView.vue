<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
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
  <v-container fluid class="px-md-8 px-4 py-4 new-article">
    <v-form ref="formRef" class="article-form" @submit.prevent="submit">
      <!-- ページ上部アクションバー: 投稿ボタンのみ -->
      <div class="action-bar d-flex justify-end py-2 mb-2">
        <v-btn
          type="submit"
          color="black"
          size="default"
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

      <!-- タイトル -->
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

      <!-- タグ -->
      <v-combobox
        v-model="form.tags"
        placeholder="タグを入力してEnter（スペース区切りで複数可）"
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

      <!-- 本文（編集 / プレビュー タブ切替） -->
      <v-tabs
        v-model="bodyTab"
        density="comfortable"
        color="primary"
        class="mb-2"
      >
        <v-tab value="edit">編集</v-tab>
        <v-tab value="preview">プレビュー</v-tab>
      </v-tabs>

      <!-- 編集 / プレビューの切替は v-show で行う。
           v-window はスライドアニメーションが副作用で残るため不採用。
           textarea を v-show で残し続けると入力状態 (caret / スクロール位置) も保たれる。 -->
      <div class="editor-area">
        <v-textarea
          v-show="bodyTab === 'edit'"
          v-model="form.body"
          placeholder="本文を入力（Markdown）"
          aria-label="本文"
          variant="plain"
          density="comfortable"
          no-resize
          maxlength="10000"
          :rules="[(v) => !!v || '本文は必須です']"
          validate-on="input"
          class="body-input"
        />

        <div v-show="bodyTab === 'preview'" class="body-preview pa-2">
          <MarkdownContent
            v-if="bodyTab === 'preview' && form.body.trim()"
            :source="form.body"
          />
          <div v-else-if="bodyTab === 'preview'" class="text-medium-emphasis">
            本文がまだ入力されていません。
          </div>
        </div>
      </div>
    </v-form>
  </v-container>
</template>

<style scoped>
/* v-main の中で viewport を使い切るレイアウト。
   ページ自体はスクロールせず、編集エリアの内部だけがスクロールする
   (Qiita 同様)。
   - height: 100% で v-main 領域を満たす
   - overflow: hidden で外側のスクロールバーを抑止
   - flex column で「上から固定の入力欄が積まれ、最後の編集エリアが残り高さを取る」構造にする */
.new-article {
  max-width: 1200px;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* v-container の唯一の子である v-form が flex column のチェーンを
   引き継がないと、内部の editor-area が「残り高さ」を計算できない。
   form 自身を flex column 化して、編集エリアまで高さの伝搬を繋ぐ。 */
.article-form {
  flex: 1 1 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

/* Vuetify の .v-input は CSS で `flex: 1 1 auto` を持っているため、
   そのまま flex column の子にすると全部が縦に均等ストレッチされてしまう。
   タイトル / タグ / 区切り線 / タブなどの「上に積む要素」は自然な高さに固定し、
   .editor-area だけが残り高さを取るようにする。 */
.article-form > * {
  flex: 0 0 auto;
}
.article-form > .editor-area {
  flex: 1 1 0;
  min-height: 0;
}

/* タイトル入力は記事のタイトルらしい大きさに */
.title-input :deep(input) {
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1.3;
}

/* 編集 / プレビュー領域: flex 子要素として残り高さ全部を取り、内部だけスクロール */
.editor-area {
  flex: 1 1 0;
  min-height: 0; /* flex 子要素が縮める許可。これが無いと内部スクロールが効かない */
  display: flex;
  flex-direction: column;
}

/* v-textarea の内部要素まで高さを伝搬させ、textarea 本体だけがスクロールするようにする。
   auto-grow を切り、no-resize で手動リサイズも禁止。
   :deep() を多段に書く必要があるのは Vuetify が v-input → v-field → ... と
   入れ子の wrapper を持つため。 */
.body-input,
.body-input :deep(.v-input__control),
.body-input :deep(.v-field),
.body-input :deep(.v-field__field) {
  height: 100%;
}
.body-input :deep(textarea) {
  height: 100% !important;
  overflow-y: auto !important;
  resize: none;
}

.body-preview {
  flex: 1 1 0;
  min-height: 0;
  overflow-y: auto;
}
</style>
