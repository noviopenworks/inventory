<template>
  <div class="flex flex-col gap-1">
    <label class="text-xs font-medium text-text-secondary">{{ label }}</label>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
}>()
</script>
