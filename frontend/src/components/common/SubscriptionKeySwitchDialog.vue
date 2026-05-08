<template>
  <BaseDialog
    :show="show && !!result && result.action !== 'none'"
    :title="dialogTitle"
    width="normal"
    @close="emit('close')"
  >
    <div class="space-y-3 text-sm leading-6 text-gray-600 dark:text-gray-300">
      <p>{{ dialogDescription }}</p>
      <div class="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-amber-700 dark:border-amber-900/40 dark:bg-amber-900/10 dark:text-amber-300">
        {{ t('subscriptionKeySwitch.warning') }}
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <button type="button" class="btn btn-secondary" @click="emit('close')">
          {{ t('common.cancel') }}
        </button>
        <button type="button" class="btn btn-primary" @click="emit('goKeys')">
          {{ t('subscriptionKeySwitch.goToKeys') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import type { SubscriptionKeySwitchResult } from '@/types/subscriptionKeySwitch'

const props = defineProps<{
  show: boolean
  result: SubscriptionKeySwitchResult | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'goKeys'): void
}>()

const { t } = useI18n()

const targetGroupName = computed(() => props.result?.target_group_name || t('subscriptionKeySwitch.defaultGroupName'))

const dialogTitle = computed(() => {
  switch (props.result?.action) {
    case 'auto_switched':
      return t('subscriptionKeySwitch.autoTitle')
    case 'no_api_keys':
      return t('subscriptionKeySwitch.noKeyTitle')
    case 'manual_switch_required':
      return t('subscriptionKeySwitch.manualTitle')
    default:
      return ''
  }
})

const dialogDescription = computed(() => {
  const groupName = targetGroupName.value
  switch (props.result?.action) {
    case 'auto_switched':
      return t('subscriptionKeySwitch.autoDescription', {
        groupName,
        keyName: props.result?.auto_switched_key_name || t('subscriptionKeySwitch.singleKeyFallbackName'),
      })
    case 'no_api_keys':
      return t('subscriptionKeySwitch.noKeyDescription', { groupName })
    case 'manual_switch_required':
      return t('subscriptionKeySwitch.manualDescription', { groupName })
    default:
      return ''
  }
})
</script>
