<script setup lang="ts">
import axios from 'axios'
import { ref, watch } from 'vue'
import { Cropper, CircleStencil } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'
import { api } from '@/api/client'

// アバターは 1 ユーザー 1 枚で上書き保存される。サーバ側は 2MB / png・jpeg のみ
// 受け付けるため、ここでクロップ結果を 512px 四方の JPEG に正規化して制約を満たす。
const MAX_OUTPUT_SIZE = 512
const MAX_UPLOAD_BYTES = 2 * 1024 * 1024
const JPEG_QUALITY = 0.9

const open = defineModel<boolean>({ required: true })
const emit = defineEmits<{ uploaded: [] }>()

const cropperRef = ref<InstanceType<typeof Cropper>>()
const selectedFile = ref<File>()
const imageSrc = ref<string>()
const submitting = ref(false)
const error = ref<string>()

// 選択ファイルが変わるたびにプレビュー用 Object URL を張り替え、前の URL は解放する。
watch(selectedFile, (file) => {
  if (imageSrc.value) URL.revokeObjectURL(imageSrc.value)
  imageSrc.value = file ? URL.createObjectURL(file) : undefined
  error.value = undefined
})

watch(open, (isOpen) => {
  if (isOpen) return
  if (imageSrc.value) URL.revokeObjectURL(imageSrc.value)
  imageSrc.value = undefined
  selectedFile.value = undefined
  error.value = undefined
})

const cropToJpegBlob = (): Promise<Blob | null> =>
  new Promise((resolve) => {
    const canvas = cropperRef.value?.getResult().canvas
    if (!canvas) {
      resolve(null)
      return
    }
    canvas.toBlob((blob) => resolve(blob), 'image/jpeg', JPEG_QUALITY)
  })

const save = async () => {
  if (submitting.value) return
  submitting.value = true
  error.value = undefined
  try {
    const blob = await cropToJpegBlob()
    if (!blob) {
      error.value = '画像を選択してください'
      return
    }
    if (blob.size > MAX_UPLOAD_BYTES) {
      error.value = '画像が大きすぎます（2MB以下にしてください）'
      return
    }

    const { data: presign } = await api.api.profileAvatarPresignCreate({
      size: blob.size,
      contentType: 'image/jpeg',
    })
    if (!presign.url || !presign.objectKey) {
      error.value = 'アバターの更新に失敗しました'
      return
    }

    // presigned PUT は MinIO へ直接送る。api クライアント（Authorization 付与）を
    // 通すと署名対象外のヘッダが混ざるため、素の axios でアップロードする。
    await axios.put(presign.url, blob, {
      headers: { 'Content-Type': 'image/jpeg' },
    })

    await api.api.profileAvatarCompleteCreate({ objectKey: presign.objectKey })

    emit('uploaded')
    open.value = false
  } catch {
    error.value = 'アバターの更新に失敗しました'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <v-dialog v-model="open" max-width="480" persistent>
    <v-card rounded="lg">
      <v-card-title>アイコンを変更</v-card-title>

      <v-card-text>
        <v-alert
          v-if="error"
          type="error"
          variant="tonal"
          class="mb-4"
          density="compact"
        >
          {{ error }}
        </v-alert>

        <v-file-input
          v-model="selectedFile"
          accept="image/png, image/jpeg"
          label="画像を選択"
          prepend-icon="mdi-camera"
          variant="outlined"
          density="comfortable"
          class="mb-4"
        />

        <div v-if="imageSrc" class="avatar-cropper rounded-lg overflow-hidden">
          <Cropper
            ref="cropperRef"
            :src="imageSrc"
            :stencil-component="CircleStencil"
            :stencil-props="{ aspectRatio: 1 }"
            :canvas="{ maxWidth: MAX_OUTPUT_SIZE, maxHeight: MAX_OUTPUT_SIZE }"
            image-restriction="stencil"
          />
        </div>
        <v-sheet
          v-else
          color="grey-lighten-3"
          rounded="lg"
          class="d-flex align-center justify-center text-medium-emphasis"
          height="240"
        >
          画像を選択してください
        </v-sheet>
      </v-card-text>

      <v-card-actions class="px-4 pb-4">
        <v-spacer />
        <v-btn variant="text" :disabled="submitting" @click="open = false">
          キャンセル
        </v-btn>
        <v-btn
          color="primary"
          variant="flat"
          :loading="submitting"
          :disabled="submitting || !imageSrc"
          @click="save"
        >
          保存
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.avatar-cropper {
  height: 320px;
  background-color: rgb(var(--v-theme-surface-variant));
}
</style>
