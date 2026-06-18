<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">All Assets</h1>
    <DataTable :columns="columns" :rows="filteredRows" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listComputers, listSmartphones, listTablets } from '@/lib/api'
import { useSearchStore } from '@/stores/search'
import { matchesQuery } from '@/lib/filter'

const columns = [
  { key: 'category', label: 'Category' },
  { key: 'name', label: 'Name' },
  { key: 'model', label: 'Model' },
  { key: 'status', label: 'Status' },
]

const rows = ref<Record<string, unknown>[]>([])
const search = useSearchStore()
const filteredRows = computed(() =>
  rows.value.filter((r) => matchesQuery(r, columns, search.query)),
)
const loading = ref(true)

onMounted(async () => {
  const [computers, smartphones, tablets] = await Promise.all([
    listComputers(),
    listSmartphones(),
    listTablets(),
  ])
  rows.value = [
    ...computers.map(r => ({ ...r, category: 'Computer' })),
    ...smartphones.map(r => ({ ...r, category: 'Smartphone' })),
    ...tablets.map(r => ({ ...r, category: 'Tablet' })),
  ]
  loading.value = false
})
</script>
