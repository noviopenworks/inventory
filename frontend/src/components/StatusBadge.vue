<template>
  <span
    class="inline-block px-1.5 py-0.5 text-xs font-medium text-white rounded-badge"
    :class="badgeClass"
  >
    {{ status }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
}>()

const badgeClass = computed(() => {
  const map: Record<string, string> = {
    active: 'bg-s-active',
    repair: 'bg-s-repair',
    spare: 'bg-s-spare',
    retired: 'bg-s-retired',
    missing: 'bg-s-missing',
    expiring: 'bg-s-expiring',
    expired: 'bg-s-expired',
  }
  return map[props.status.toLowerCase()] ?? 'bg-s-retired'
})
</script>
