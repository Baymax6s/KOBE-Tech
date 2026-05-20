<script setup lang="ts">
import { SUPPORTED_LANGUAGES } from './markdownLanguages'

defineOptions({ name: 'SupportedLanguagesDialog' })

defineProps<{ modelValue: boolean }>()
defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    max-width="640"
    scrollable
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center ga-2">
        <v-icon icon="mdi-code-tags" />
        <span>シンタックスハイライト対応言語</span>
        <v-spacer />
        <v-btn
          icon="mdi-close"
          variant="text"
          density="comfortable"
          aria-label="閉じる"
          @click="$emit('update:modelValue', false)"
        />
      </v-card-title>

      <v-divider />

      <v-card-text class="pa-0">
        <v-list density="compact">
          <v-list-item v-for="lang in SUPPORTED_LANGUAGES" :key="lang.name">
            <template #title>
              <code class="text-body-2">{{ lang.name }}</code>
            </template>
            <template v-if="lang.aliases.length" #subtitle>
              <span class="text-caption text-medium-emphasis">
                エイリアス: {{ lang.aliases.join(', ') }}
              </span>
            </template>
          </v-list-item>
        </v-list>
      </v-card-text>

      <v-divider />

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="$emit('update:modelValue', false)">
          閉じる
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
