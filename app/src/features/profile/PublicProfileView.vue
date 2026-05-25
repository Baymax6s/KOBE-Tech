<script setup lang="ts">
import axios from 'axios'
import { ref, onMounted } from 'vue'
import { api } from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const props = withDefaults(
  defineProps<{
    userId?: number
    isMe?: boolean
  }>(),
  {
    userId: 0,
    isMe: false,
  },
)

const auth = useAuthStore()

const profile = ref<{
  id?: number
  name?: string
  bio?: string
  is_owner?: boolean
} | null>(null)

const loading = ref(true)
const error = ref<string | null>(null)

const isEditing = ref(false)
const bio = ref('')
const submitting = ref(false)

const maxLength = 200

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    let targetId = props.userId
    if (props.isMe) {
      if (auth.userId === null) {
        const id = await auth.fetchUser()
        if (id === null) {
          error.value = 'ユーザー情報の取得に失敗しました'
          loading.value = false
          return
        }
        targetId = id
      } else {
        targetId = auth.userId
      }
    }
    const res = await api.api.profileDetail(targetId)
    profile.value = res.data
    bio.value = res.data.bio ?? ''
  } catch (err) {
    if (axios.isAxiosError(err) && err.response?.status === 404) {
      error.value = 'ユーザーが見つかりませんでした'
    } else {
      error.value = 'プロフィールの取得に失敗しました'
    }
  } finally {
    loading.value = false
  }
})

const saveBio = async () => {
  if (submitting.value) return
  if (bio.value.length > maxLength) return
  submitting.value = true

  try {
    const res = await api.api.profileBioUpdate({
      bio: bio.value,
    })
    profile.value = res.data
    bio.value = res.data.bio ?? ''
    isEditing.value = false
  } catch {
    error.value = '更新に失敗しました'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <v-container class="py-8">
    <v-row justify="center">
      <v-col cols="12" sm="10">
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

        <v-card v-if="profile" class="pa-6 text-center elevation-3">
          <v-avatar size="80" class="mx-auto mb-2" color="indigo-lighten-1">
            <v-icon size="50" color="white"> mdi-account-circle </v-icon>
          </v-avatar>

          <h2 class="text-h5 font-weight-bold mb-1">
            {{ profile.name }}
          </h2>

          <v-divider class="my-4" />

          <h3 class="text-subtitle-1 font-weight-bold mb-2">自己紹介</h3>

          <div class="text-left">
            <div v-if="!isEditing">
              <p class="mb-4">
                {{ profile.bio || '自己紹介はまだありません' }}
              </p>

              <v-btn
                v-if="profile.is_owner"
                variant="text"
                color="primary"
                @click="isEditing = true"
              >
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

              <v-btn
                color="primary"
                class="mr-2"
                :loading="submitting"
                :disabled="submitting"
                @click="saveBio"
              >
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
