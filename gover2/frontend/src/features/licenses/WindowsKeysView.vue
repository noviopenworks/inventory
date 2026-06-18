<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Windows Keys</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listWindowsKeys } from '@/lib/api'
import type { WindowsKey } from '@/lib/api'

const columns = [
  { key: 'licenseKey', label: 'License Key' },
  { key: 'computerId', label: 'Computer ID' },
  { key: 'status', label: 'Status' },
  { key: 'notes', label: 'Notes' },
]

const rows = ref<WindowsKey[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listWindowsKeys()
  loading.value = false
})
</script>
