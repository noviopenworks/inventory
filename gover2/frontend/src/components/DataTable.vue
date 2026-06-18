<template>
  <div class="bg-surface border border-border overflow-hidden">
    <table class="w-full text-sm">
      <thead class="bg-th-bg">
        <tr>
          <th
            v-for="col in columns"
            :key="col.key"
            class="text-left px-4 py-2 text-th-text font-medium"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-text-secondary">
            Loading...
          </td>
        </tr>
        <tr v-else-if="rows.length === 0">
          <td :colspan="columns.length" class="px-4 py-8 text-center text-text-secondary">
            No items to display.
          </td>
        </tr>
        <tr
          v-else
          v-for="(row, i) in rows"
          :key="i"
          class="border-t border-border hover:bg-sidebar/10"
        >
          <td
            v-for="col in columns"
            :key="col.key"
            class="px-4 py-2 text-text-primary"
          >
            {{ row[col.key] ?? '' }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  columns: { key: string; label: string }[]
  rows: Record<string, unknown>[]
  loading: boolean
}>()
</script>
