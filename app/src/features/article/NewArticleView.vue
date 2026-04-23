<script setup lang="ts">
import { reactive } from 'vue'

defineOptions({
  name: 'NewArticleView'
})

const form = reactive({
  title: '',
  body: '',
  // tags フィールドは将来追加用に構造化（現在は非表示）
})

const submit = () => {
  console.log('投稿内容', {
    title: form.title,
    body: form.body,
  })

  // ここで API 送信や状態更新を行います。
  form.title = ''
  form.body = ''
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
            <v-form @submit.prevent="submit" class="form-wrapper">
              <div class="form-fields">
                <!-- タイトル入力フィールド -->
                <v-text-field
                  v-model="form.title"
                  label="タイトル"
                  placeholder="タイトルを入力"
                  variant="outlined"
                  density="comfortable"
                  clearable
                  counter
                  maxlength="200"
                />

                <!-- 本文入力フィールド -->
                <v-textarea
                  v-model="form.body"
                  label="本文"
                  placeholder="本文を入力"
                  rows="10"
                  variant="outlined"
                  density="comfortable"
                  counter
                  maxlength="10000"
                />

                <!-- 投稿ボタン -->
                <div class="button-row mt-4">
                  <v-btn
                    type="submit"
                    color="black"
                    size="large"
                    class="px-8"
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
