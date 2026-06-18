<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Antivirus</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listAntivirus } from '@/lib/api'
import type { Antivirus } from '@/lib/api'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'licenseKey', label: 'License Key' },
  { key: 'status', label: 'Status' },
  { key: 'expiryDate', label: 'Expiry Date' },
]

const rows = ref<Antivirus[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listAntivirus()
  loading.value = false
})
</script>
