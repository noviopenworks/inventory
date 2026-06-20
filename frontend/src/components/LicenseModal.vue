<template>
  <div v-if="open" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
    <div class="bg-surface border border-border p-6 w-[480px] max-h-[70vh] overflow-y-auto">
      <h2 class="text-base font-semibold text-text-primary mb-2">License</h2>
      <p class="text-sm text-text-secondary whitespace-pre-line">{{ licenseText }}</p>
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

const licenseText = 'This software is distributed under the GNU General Public License v3 (GPL v3). See the LICENSE file in the repository for the full terms.'
</script>
