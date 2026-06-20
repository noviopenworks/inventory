<template>
  <!-- full-viewport overlay -->
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
       @click.self="$emit('cancelled')">
    <div class="bg-surface border border-border rounded-lg p-6 max-w-sm w-full mx-4">
      <p class="text-text-primary mb-4">{{ message }}</p>
      <div class="flex gap-3 justify-end">
        <button @click="$emit('cancelled')"
                class="px-4 py-2 text-text-secondary hover:text-text-primary border border-border rounded">
          Cancel
        </button>
        <button @click="$emit('confirmed')"
                class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded">
          Delete
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ message: string }>()
defineEmits<{ confirmed: []; cancelled: [] }>()
</script>
