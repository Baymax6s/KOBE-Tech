<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'
import { useArticleNotificationStore } from '@/stores/articleNotification'

defineOptions({
  name: 'NewArticleView'
})

type ArticleFormRef = {
  validate: () => Promise<{ valid: boolean }>
  reset: () => void
}

const formRef = ref<ArticleFormRef | null>(null)

const form = reactive({
  title: '',
  body: '',
})

const router = useRouter()
const notificationStore = useArticleNotificationStore()

const submitting = ref(false)
const submitError = ref<string | null>(null)

const canSubmit = computed(
  () => !submitting.value && !!form.title.trim() && !!form.body.trim()
)

const submit = async () => {
  if (submitting.value) return
  if (!formRef.value) return
  const { valid } = await formRef.value.validate()
  if (!valid) return

  submitting.value = true
  submitError.value = null

  try {
    await api.api.articlesCreate({ title: form.title, content: form.body })
    notificationStore.markCreated()
    await router.push('/articles')
  } catch {
    submitError.value = '投稿に失敗しました。もう一度お試しください。'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <v-container class="py-12" fluid>
    <v-row justify="center" align="start">
      <v-col cols="12" sm="10" md="8" lg="6">
        <v-card elevation="4" rounded="lg">
          <v-card-title class="text-h4 font-weight-bold pa-6 pb-2">
            新規記事作成
          </v-card-title>

          <v-card-text class="pa-6">
            <v-alert v-if="submitError" type="error" class="mb-4" closable @click:close="submitError = null">
              {{ submitError }}
            </v-alert>

            <v-form ref="formRef" @submit.prevent="submit" class="form-wrapper">
              <div class="form-fields">
                <v-text-field
                  v-model="form.title"
                  label="タイトル"
                  placeholder="タイトルを入力"
                  variant="outlined"
                  density="comfortable"
                  clearable
                  counter
                  maxlength="200"
                  :rules="[v => !!v || 'タイトルは必須です']"
                  validate-on="input"
                />

                <v-textarea
                  v-model="form.body"
                  label="本文"
                  placeholder="本文を入力"
                  rows="10"
                  variant="outlined"
                  density="comfortable"
                  counter
                  maxlength="10000"
                  :rules="[v => !!v || '本文は必須です']"
                  validate-on="input"
                />

                <div class="button-row mt-4">
                  <v-btn
                    type="submit"
                    color="black"
                    size="large"
                    class="px-8"
                    :loading="submitting"
                    :disabled="!canSubmit"
                  >
                    投稿
                  </v-btn>
                </div>
              </div>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.form-wrapper {
  display: block;
}

.form-fields {
  display: grid;
  gap: 24px;
}

.button-row {
  display: flex;
  justify-content: flex-start;
}
</style>
