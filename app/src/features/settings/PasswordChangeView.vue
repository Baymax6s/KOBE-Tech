<template>
  <v-container class="fill-height bg-grey-lighten-4" fluid>
    <v-row justify="center" align="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="pa-6">
          <v-card-title class="text-h6 text-center">パスワード変更</v-card-title>

          <v-card-text>
            <v-alert
              v-if="errorMessage"
              type="error"
              class="mb-4"
              closable
              @click:close="errorMessage = null"
            >
              {{ errorMessage }}
            </v-alert>

            <v-alert
              v-if="successMessage"
              type="success"
              class="mb-4"
              closable
              @click:close="successMessage = null"
            >
              {{ successMessage }}
            </v-alert>

            <v-form @submit.prevent="onSubmit">
              <v-text-field
                v-model="currentPassword"
                label="現在のパスワード"
                type="password"
                variant="outlined"
                required
                class="mb-4"
                :disabled="submitting"
              />

              <v-text-field
                v-model="newPassword"
                label="新しいパスワード"
                type="password"
                variant="outlined"
                required
                class="mb-4"
                :disabled="submitting"
              />

              <v-btn
                type="submit"
                color="primary"
                block
                :loading="submitting"
                :disabled="!currentPassword || !newPassword || submitting"
              >
                変更する
              </v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { api } from '@/api/client'

const currentPassword = ref('')
const newPassword = ref('')
const submitting = ref(false)
const errorMessage = ref<string | null>(null)
const successMessage = ref<string | null>(null)

const onSubmit = async () => {
  if (submitting.value) return
  submitting.value = true
  errorMessage.value = null
  successMessage.value = null
  try {
    await api.instance.put('/api/auth/password', {
      current_password: currentPassword.value,
      new_password: newPassword.value,
    })
    successMessage.value = 'パスワードを変更しました'
    currentPassword.value = ''
    newPassword.value = ''
  } catch (e) {
    if (axios.isAxiosError(e) && e.response?.status === 401) {
      errorMessage.value =
        e.response.data?.message ?? '現在のパスワードが正しくありません'
    } else if (axios.isAxiosError(e) && e.response?.status === 400) {
      errorMessage.value =
        e.response.data?.message ?? '入力内容に誤りがあります'
    } else {
      errorMessage.value =
        'パスワードの変更に失敗しました。時間をおいて再度お試しください'
    }
  } finally {
    submitting.value = false
  }
}
</script>
