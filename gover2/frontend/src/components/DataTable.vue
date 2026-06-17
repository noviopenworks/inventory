<template>
  <div class="bg-surface border border-border rounded-none overflow-hidden">
    <table class="w-full text-sm">
      <thead class="bg-th-bg">
        <tr>
          <th class="text-left px-4 py-2 text-th-text font-medium">Name</th>
          <th class="text-left px-4 py-2 text-th-text font-medium">Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="items.length === 0">
          <td colspan="2" class="px-4 py-8 text-center text-text-secondary">
            No items to display.
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  items: unknown[]
}>()
</script>
