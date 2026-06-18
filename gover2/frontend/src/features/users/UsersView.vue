<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Users</h1>
    <DataTable :columns="columns" :rows="rows as Record<string, unknown>[]" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import { listUsers } from '@/lib/api'
import type { User } from '@/lib/api'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'surname', label: 'Surname' },
  { key: 'status', label: 'Status' },
  { key: 'notes', label: 'Notes' },
]

const rows = ref<User[]>([])
const loading = ref(true)

onMounted(async () => {
  rows.value = await listUsers()
  loading.value = false
})
</script>
