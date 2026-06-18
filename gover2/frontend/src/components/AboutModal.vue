<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border p-6 w-[420px]">
      <h2 class="text-base font-semibold text-text-primary mb-2">About</h2>
      <p class="text-sm text-text-primary">Inventory</p>
      <p class="text-sm text-text-secondary mt-1">A desktop inventory manager built with Wails and Vue.</p>
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
defineProps<{ open: boolean }>()
defineEmits<{ (e: 'close'): void }>()
</script>
