<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border p-6 w-[560px] max-h-[70vh] overflow-y-auto">
      <h2 class="text-base font-semibold text-text-primary mb-4">Alerts</h2>

      <p v-if="alerts.length === 0" class="text-text-secondary text-sm">No alerts.</p>

      <table v-else class="w-full text-sm">
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
                {{ a.severity }}
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
import type { Alert } from '@/lib/api'

defineProps<{ open: boolean; alerts: Alert[] }>()
defineEmits<{ (e: 'close'): void }>()
</script>
