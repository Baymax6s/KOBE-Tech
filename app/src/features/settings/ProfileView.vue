<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import axios from 'axios'

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

const auth = useAuthStore()
const router = useRouter()

const avatarImageUrl = ref<string | null | undefined>(null)

onMounted(async () => {
  if (!auth.isAuthenticated) {
    await router.push('/login')
    return
  }

  loading.value = true
  error.value = null

  

  try {
    const res = await api.api.profileList()
    user.value = res.data
    bio.value = res.data.bio ?? ''

    if (res.data.objectKey) {
      const urlRes = await api.api.profileAvatarDownloadCreate({
        objectKey: res.data.objectKey
      })
      avatarImageUrl.value = urlRes.data.url
    }

  } catch {
    error.value = 'プロフィールの取得に失敗しました'
  } finally {
    loading.value = false
  }
})

const saveBio = async () => {
  if (bio.value.length > maxLength) return

  try {
    const res = await api.api.profileBioUpdate({
      bio: bio.value,
    })

    user.value = res.data
    bio.value = res.data.bio ?? ''

    isEditing.value = false
  } catch {
    error.value = '更新に失敗しました'
  }
}


const fileInput = ref<HTMLInputElement | null>(null)

const avatarPreview = ref<string | null>(null)

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleImageUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const file = target.files[0]
  if (!file) return

  if (file.size > 2 * 1024 * 1024) {
    error.value = '画像サイズは2MB以下にしてください'
    return
  }
  if (file.type !== 'image/png' && file.type !== 'image/jpeg') {
    error.value = '対応している画像フォーマットはPNGまたはJPEGのみです'
    return
  }

  loading.value = true
  error.value = null

  try {
    avatarPreview.value = URL.createObjectURL(file)

    const presignRes = await api.api.profileAvatarPresignCreate({
      size: file.size,
      contentType: file.type,
    })

    const { url, objectKey } = presignRes.data

    if (!url || !objectKey) {
      throw new Error('サーバーから有効なURLが返されませんでした')
    }

    await axios.put(url, file, {
      headers: { 'Content-Type': file.type },
    })

    await api.api.profileAvatarCompleteCreate({
      objectKey: objectKey,
    })

    const res = await api.api.profileList()
    user.value = res.data
    bio.value = res.data.bio ?? ''

  } catch (err) {
    console.error('アップロードエラー:', err)
    error.value = '画像のアップロードに失敗しました'
    avatarPreview.value = null
  } finally {
    loading.value = false
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
          
          <input
            ref="fileInput"
            type="file"
            accept="image/png, image/jpeg"
            style="display: none"
            @change="handleImageUpload"
          />

          <v-avatar 
            size="70" 
            class="mx-auto mb-2 cursor-pointer" 
            color="indigo-lighten-1"
            @click="triggerFileInput"
          >
            <v-img v-if="avatarPreview" :src="avatarPreview" alt="Avatar Preview" />
            
            <v-img v-else-if="avatarImageUrl" :src="avatarImageUrl" alt="Avatar Image" />
            
            <v-icon v-else size="40" color="white"> mdi-account-circle </v-icon>
          </v-avatar>

          <div class="text-caption text-grey mb-2">アバターをクリックして変更</div>

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
