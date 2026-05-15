<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'

type ProfileResponse = {
  id?: number
  name?: string
  bio?: string
}

const user = ref<ProfileResponse | null>(null)

const loading = ref(false)
const error = ref<string | null>(null)

const isEditing = ref(false)
const bio = ref('')

const maxLength = 200

onMounted(async () => {
  loading.value = true
  error.value = null

  try {
    const res = await api.api.profileList()
    user.value = res.data
    bio.value = res.data.bio ?? ''
  } catch {
    error.value = 'プロフィールの取得に失敗しました'
  } finally {
    loading.value = false
  }
})

const saveBio = async () => {
  if (bio.value.length > maxLength) return

  try {
    await api.api.profileUpdate({
      bio: bio.value,
    })

    if (user.value) {
      user.value.bio = bio.value
    }

    isEditing.value = false
  } catch {
    error.value = '更新に失敗しました'
  }
}
</script>

<template>
  <v-container class="py-8">
    <v-row justify="center">
      <v-col cols="12" md="8" lg="6">
        <v-alert
          v-if="error"
          type="error"
          class="mb-4"
          closable
          @click:close="error = null"
        >
          {{ error }}
        </v-alert>

        <div v-if="loading" class="d-flex justify-center py-12">
          <v-progress-circular indeterminate color="primary" />
        </div>

        <v-card v-if="user" class="pa-6 text-center elevation-3">
          <v-avatar size="70" class="mx-auto mb-2" color="indigo-lighten-1">
            <v-icon size="40" color="white"> mdi-account-circle </v-icon>
          </v-avatar>

          <h2 class="text-h5 font-weight-bold mb-1">
            {{ user.name }}
          </h2>

          <v-divider class="my-4" />

          <h3 class="text-subtitle-1 font-weight-bold mb-2">自己紹介</h3>

          <div class="text-left">
            <div v-if="!isEditing">
              <p class="mb-4">
                {{ user.bio || '自己紹介はまだありません' }}
              </p>

              <v-btn variant="text" color="primary" @click="isEditing = true">
                編集
              </v-btn>
            </div>

            <div v-else>
              <v-textarea
                v-model="bio"
                :counter="maxLength"
                :rules="[
                  (v) =>
                    v?.length <= maxLength ||
                    `${maxLength}文字以内で入力してください`,
                ]"
                label="自己紹介"
                variant="outlined"
                class="mb-3"
              />

              <v-btn color="primary" class="mr-2" @click="saveBio">
                完了
              </v-btn>

              <v-btn variant="text" @click="isEditing = false">
                キャンセル
              </v-btn>
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
