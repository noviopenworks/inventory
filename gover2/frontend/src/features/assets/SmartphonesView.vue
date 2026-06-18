<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Smartphones</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listSmartphones } from '@/lib/api'
import type { Smartphone } from '@/lib/api'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'model', label: 'Model' },
  { key: 'status', label: 'Status' },
  { key: 'purchaseDate', label: 'Purchase Date' },
  { key: 'warrantyExpiry', label: 'Warranty Expiry' },
]

const rows = ref<Smartphone[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listSmartphones()
  loading.value = false
})
</script>
