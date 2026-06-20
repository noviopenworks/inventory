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
          <th class="text-left px-4 py-2 text-th-text font-medium">Delete</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="loading">
          <td :colspan="columns.length + 1" class="px-4 py-8 text-center text-text-secondary">
            Loading...
          </td>
        </tr>
        <tr v-else-if="rows.length === 0">
          <td :colspan="columns.length + 1" class="px-4 py-8 text-center text-text-secondary">
            No items to display.
          </td>
        </tr>
        <tr
          v-else
          v-for="(row, i) in rows"
          :key="i"
          class="border-t border-border even:bg-row-alt hover:bg-row-hover cursor-pointer"
          @click="$emit('row-click', row)"
        >
          <td
            v-for="col in columns"
            :key="col.key"
            :class="['px-4', cellPad, 'text-text-primary']"
          >
            {{ row[col.key] ?? '' }}
          </td>
          <td :class="['px-4', cellPad]" @click.stop>
            <button
              class="px-2 py-1 text-xs text-white bg-red-600 hover:bg-red-700 rounded"
              @click="$emit('delete', row)"
            >
              Delete
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUiStore } from '@/stores/ui'

defineProps<{
  columns: { key: string; label: string }[]
  rows: Record<string, unknown>[]
  loading: boolean
}>()

defineEmits<{
  'row-click': [row: Record<string, unknown>]
  'delete': [row: Record<string, unknown>]
}>()

const ui = useUiStore()
const cellPad = computed(() => (ui.density === 'compact' ? 'py-1' : 'py-2'))
</script>
