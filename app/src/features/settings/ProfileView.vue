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
      // 取得したURLを変数にセットする
      avatarImageUrl.value = urlRes.data.url
    }

    // 👇 ✨ここにこれを追加して、ブラウザのコンソールで中身を見てみてください！
    console.log('現在のユーザープロフィールデータ:', res.data)
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

    // APIから返ってきた最新のデータで更新
    user.value = res.data
    bio.value = res.data.bio ?? ''

    isEditing.value = false
  } catch {
    error.value = '更新に失敗しました'
  }
}

// ...既存のコード（saveBio関数の終わりなど）のすぐ下に追加...

// 1. ファイル入力欄（DOM）にアクセスするための参照
const fileInput = ref<HTMLInputElement | null>(null)

// 2. 選択した画像の一時URLを入れておく箱
const avatarPreview = ref<string | null>(null)

// 3. アバターがクリックされたら、隠しインプットを代わりにクリックする関数
const triggerFileInput = () => {
  fileInput.value?.click()
}

// 4. 画像が選ばれたら動く関数
const handleImageUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return

  const file = target.files[0]
  if (!file) return

  // フロントエンド側での事前のサイズ・形式バリデーション（Goのバリデーションと合わせておく）
  if (file.size > 2 * 1024 * 1024) {
    error.value = '画像サイズは2MB以下にしてください'
    return
  }
  if (file.type !== 'image/png' && file.type !== 'image/jpeg') {
    error.value = '対応している画像フォーマットはPNGまたはJPEGのみです'
    return
  }

  loading.value = true // 画面をローディング表示にする
  error.value = null

  try {
    // 1. プレビューを即座に表示（ユーザー体験のため）
    avatarPreview.value = URL.createObjectURL(file)

    // 2. [ステップ 1 & 2] Goのバックエンドから MinIO アップロード用のURLを生成してもらう
    const presignRes = await api.api.profileAvatarPresignCreate({
      size: file.size,
      contentType: file.type,
    })

    const { url, objectKey } = presignRes.data

    // 👇 ✨この2行を追加！ url または objectKey が無かったら処理を中断する
    if (!url || !objectKey) {
      throw new Error('サーバーから有効なURLが返されませんでした')
    }

    // 3. [ステップ 3] これで url と objectKey が「絶対に string である」と確定するのでエラーが消えます！
    await axios.put(url, file, {
      headers: { 'Content-Type': file.type },
    })

    // 4. [ステップ 4 & 5] Goのバックエンドに「アップロード完了したよ」と伝える
    await api.api.profileAvatarCompleteCreate({
      objectKey: objectKey, // 👈 ここもエラーが消えます！
    })

    // 5. 最後にプロフィール情報を再取得して、画面を最新の状態にする
    const res = await api.api.profileList()
    user.value = res.data
    bio.value = res.data.bio ?? ''

  } catch (err) {
    console.error('アップロードエラー:', err)
    error.value = '画像のアップロードに失敗しました'
    // 失敗したらプレビューを消す、などの処理を入れても良いです
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
