<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Alerts</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { getAlerts } from '@/lib/api'
import type { Alert } from '@/lib/api'

const columns = [
  { key: 'category', label: 'Category' },
  { key: 'name', label: 'Name' },
  { key: 'expiryDate', label: 'Expiry Date' },
  { key: 'daysRemaining', label: 'Days Remaining' },
  { key: 'severity', label: 'Severity' },
]

const rows = ref<Alert[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await getAlerts()
  loading.value = false
})
</script>
