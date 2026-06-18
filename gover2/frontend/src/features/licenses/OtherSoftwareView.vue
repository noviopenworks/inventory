<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Other Software</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listOtherSoftware } from '@/lib/api'
import type { OtherSoftware } from '@/lib/api'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'licenseKey', label: 'License Key' },
  { key: 'status', label: 'Status' },
  { key: 'expiryDate', label: 'Expiry Date' },
]

const rows = ref<OtherSoftware[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listOtherSoftware()
  loading.value = false
})
</script>
