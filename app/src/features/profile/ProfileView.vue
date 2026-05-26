<script setup lang="ts">
import axios from 'axios'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArticleCard from '@/features/articles/ArticleCard.vue'
import AvatarEditDialog from '@/features/profile/AvatarEditDialog.vue'
import { api } from '@/api/client'
import type {
  ServerArticleJSONResponse,
  ServerProfileJSON,
} from '@/api/generated/apiSchema'
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

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const profile = ref<ServerProfileJSON | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

const postedArticles = ref<ServerArticleJSONResponse[]>([])
const likedArticles = ref<ServerArticleJSONResponse[]>([])
const listsLoading = ref(false)
const listsError = ref<string | null>(null)

const isEditing = ref(false)
const avatarDialogOpen = ref(false)
const bio = ref('')
const submitting = ref(false)
const maxLength = 200

// 記事一覧の絞り込みと同じ思想で、開いているタブを URL に持たせる。
// リロード・共有でタブ状態を再現できるようにするため。
const activeTab = computed<'articles' | 'likes'>({
  get() {
    return route.query.tab === 'likes' ? 'likes' : 'articles'
  },
  set(next) {
    void router.replace({
      query: { ...route.query, tab: next === 'likes' ? 'likes' : undefined },
    })
  },
})

const resolveTargetUserId = async (): Promise<number | null> => {
  if (!props.isMe) return props.userId
  if (auth.userId !== null) return auth.userId
  return auth.fetchUser()
}

const fetchLists = async (targetId: number) => {
  listsLoading.value = true
  listsError.value = null
  try {
    const [posted, liked] = await Promise.all([
      api.api.profileArticlesList(targetId),
      api.api.profileLikedArticlesList(targetId),
    ])
    postedArticles.value = posted.data.articles ?? []
    likedArticles.value = liked.data.articles ?? []
  } catch {
    listsError.value = '記事一覧の取得に失敗しました'
  } finally {
    listsLoading.value = false
  }
}

const fetchProfile = async () => {
  loading.value = true
  error.value = null
  // 別ユーザーへ遷移した直後にエラーが出ても前のプロフィールが残らないよう、
  // 取得開始時に表示状態をリセットする。
  profile.value = null
  postedArticles.value = []
  likedArticles.value = []
  listsError.value = null
  isEditing.value = false
  try {
    const targetId = await resolveTargetUserId()
    if (targetId === null) {
      error.value = 'ユーザー情報の取得に失敗しました'
      return
    }
    const res = await api.api.profileDetail(targetId)
    profile.value = res.data
    bio.value = res.data.bio ?? ''
    await fetchLists(targetId)
  } catch (err) {
    if (axios.isAxiosError(err) && err.response?.status === 404) {
      error.value = 'ユーザーが見つかりませんでした'
    } else {
      error.value = 'プロフィールの取得に失敗しました'
    }
  } finally {
    loading.value = false
  }
}

watch(() => [props.userId, props.isMe], fetchProfile, { immediate: true })

const saveBio = async () => {
  if (submitting.value) return
  if (bio.value.length > maxLength) return
  submitting.value = true
  try {
    const res = await api.api.profileBioUpdate({ bio: bio.value })
    profile.value = res.data
    bio.value = res.data.bio ?? ''
    isEditing.value = false
  } catch {
    error.value = '更新に失敗しました'
  } finally {
    submitting.value = false
  }
}

// アバターは編集モードのときだけ変更できる（オーバーレイのカメラから開く）。
const onAvatarClick = () => {
  if (isEditing.value) avatarDialogOpen.value = true
}

// アバター更新後は presigned GET URL を更新したいだけなので、記事一覧は再取得せず
// プロフィール本体だけを読み直す。
const reloadProfile = async () => {
  const targetId = await resolveTargetUserId()
  if (targetId === null) return
  const res = await api.api.profileDetail(targetId)
  profile.value = res.data
}

// プロフィールではタグ絞り込みを持たないので、タグをクリックしたら
// 記事一覧画面の絞り込みへ遷移させる。
const goToTag = (tagName: string) => {
  void router.push({ path: '/articles', query: { tag: tagName } })
}
</script>

<template>
  <v-sheet color="grey-lighten-4" min-height="100%">
    <v-container class="py-8">
      <v-row justify="center">
        <v-col cols="12" sm="10" md="8">
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

          <template v-else-if="profile">
            <v-card class="pa-6 mb-6" rounded="lg">
              <div class="d-flex ga-6 align-start">
                <div
                  class="avatar-edit"
                  :class="{ 'avatar-edit--active': isEditing }"
                  @click="onAvatarClick"
                >
                  <v-avatar size="96" color="indigo-lighten-1">
                    <v-img
                      v-if="profile.avatar_url"
                      :src="profile.avatar_url"
                      alt="アバター"
                    />
                    <v-icon v-else size="64" color="white">
                      mdi-account-circle
                    </v-icon>
                  </v-avatar>
                  <div v-if="isEditing" class="avatar-edit__overlay">
                    <v-icon color="white">mdi-camera</v-icon>
                  </div>
                </div>

                <div class="flex-grow-1">
                  <h1 class="text-h5 font-weight-bold mb-2">
                    {{ profile.name }}
                  </h1>

                  <p v-if="!isEditing" class="text-body-2 text-medium-emphasis">
                    {{ profile.bio || '自己紹介はまだありません' }}
                  </p>

                  <template v-else>
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
                      rows="3"
                      class="mb-2"
                    />
                    <div class="d-flex ga-2">
                      <v-btn
                        color="primary"
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
                  </template>
                </div>

                <v-btn
                  v-if="profile.is_owner && !isEditing"
                  variant="text"
                  color="primary"
                  size="small"
                  prepend-icon="mdi-pencil"
                  @click="isEditing = true"
                >
                  編集
                </v-btn>
              </div>
            </v-card>

            <AvatarEditDialog
              v-if="profile.is_owner"
              v-model="avatarDialogOpen"
              @uploaded="reloadProfile"
            />

            <v-tabs v-model="activeTab" color="primary" class="mb-4">
              <v-tab value="articles">投稿した記事</v-tab>
              <v-tab value="likes">いいねした記事</v-tab>
            </v-tabs>

            <div v-if="listsLoading" class="d-flex justify-center py-12">
              <v-progress-circular indeterminate color="primary" />
            </div>

            <v-alert v-else-if="listsError" type="error">
              {{ listsError }}
            </v-alert>

            <v-window v-else v-model="activeTab">
              <v-window-item value="articles">
                <div class="d-flex flex-column ga-4">
                  <ArticleCard
                    v-for="article in postedArticles"
                    :key="article.id"
                    :article="article"
                    :selected-tags="[]"
                    @select-tag="goToTag"
                  />
                  <v-alert
                    v-if="postedArticles.length === 0"
                    type="info"
                    variant="tonal"
                  >
                    まだ投稿した記事がありません
                  </v-alert>
                </div>
              </v-window-item>

              <v-window-item value="likes">
                <div class="d-flex flex-column ga-4">
                  <ArticleCard
                    v-for="article in likedArticles"
                    :key="article.id"
                    :article="article"
                    :selected-tags="[]"
                    @select-tag="goToTag"
                  />
                  <v-alert
                    v-if="likedArticles.length === 0"
                    type="info"
                    variant="tonal"
                  >
                    まだいいねした記事がありません
                  </v-alert>
                </div>
              </v-window-item>
            </v-window>
          </template>
        </v-col>
      </v-row>
    </v-container>
  </v-sheet>
</template>

<style scoped>
/* 編集モード時にアバター上へカメラを重ねる。Vuetify にアバター用の
   オーバーレイ表現が無いため、ここだけ最小限の scoped CSS で実装する。 */
.avatar-edit {
  position: relative;
  display: inline-flex;
  border-radius: 50%;
}

.avatar-edit--active {
  cursor: pointer;
}

.avatar-edit__overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: rgba(0, 0, 0, 0.45);
  transition: background-color 0.2s ease;
}

.avatar-edit--active:hover .avatar-edit__overlay {
  background-color: rgba(0, 0, 0, 0.6);
}
</style>
