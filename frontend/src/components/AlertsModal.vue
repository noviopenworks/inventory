<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border p-6 w-[560px] max-h-[70vh] overflow-y-auto">
      <h2 class="text-base font-semibold text-text-primary mb-2">Alerts</h2>

      <p v-if="alerts.length > 0" class="text-sm text-text-secondary mb-4">
        <b class="text-text-primary">{{ expiredCount }}</b> expired
        &nbsp;·&nbsp;
        <b class="text-text-primary">{{ expiringCount }}</b> expiring within {{ warningDays }} days
      </p>

      <p v-else class="text-text-secondary text-sm">No alerts.</p>

      <table v-if="alerts.length > 0" class="w-full text-sm">
        <thead class="bg-th-bg">
          <tr>
            <th class="text-left px-3 py-2 text-th-text font-medium">Category</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Name</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Expiry</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Days</th>
            <th class="text-left px-3 py-2 text-th-text font-medium">Severity</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in alerts" :key="a.category + '-' + a.id" class="border-t border-border">
            <td class="px-3 py-2 text-text-primary">{{ a.category }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.name }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.expiryDate }}</td>
            <td class="px-3 py-2 text-text-primary">{{ a.daysRemaining }}</td>
            <td class="px-3 py-2">
              <span
                class="px-2 py-0.5 rounded-badge text-xs text-white"
                :class="a.severity === 'expired' ? 'bg-s-expired' : 'bg-s-expiring'"
              >
                {{ severityLabel[a.severity] ?? a.severity }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <button
        class="mt-4 px-3 py-1.5 text-sm bg-s-active text-white hover:bg-accent-hover"
        @click="$emit('close')"
      >
        Close
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Alert } from '@/lib/api'

const props = withDefaults(defineProps<{
  open: boolean
  alerts: Alert[]
  warningDays?: number
}>(), { warningDays: 30 })

defineEmits<{ (e: 'close'): void }>()

const severityLabel: Record<string, string> = {
  expired: 'Expired',
  expiring: 'Expiring soon',
}

const expiredCount = computed(() => props.alerts.filter(a => a.severity === 'expired').length)
const expiringCount = computed(() => props.alerts.length - expiredCount.value)
</script>
