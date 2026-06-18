<template>
  <div class="p-4 text-text-primary">
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-lg font-semibold">Users</h1>
      <button @click="onAdd"
              class="px-3 py-1.5 bg-accent text-white text-sm rounded hover:opacity-90">
        Add
      </button>
    </div>
    <DataTable
      :columns="columns"
      :rows="filteredRows"
      :loading="loading"
      @row-click="onRowClick"
      @delete="onDeleteClick"
    />
    <EditPanel
      v-if="panelOpen"
      :category="'users'"
      :mode="panelMode"
      :row="selectedRow"
      :users="[]"
      :devices="[]"
      @saved="onSaved"
      @cancelled="panelOpen = false"
    />
    <ConfirmDialog
      v-if="confirmOpen"
      message="Delete this user?"
      @confirmed="onDeleteConfirmed"
      @cancelled="confirmOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import DataTable from '@/components/DataTable.vue'
import EditPanel from '@/components/EditPanel.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { listUsers, deleteUser } from '@/lib/api'
import type { User } from '@/lib/api'
import { useSearchStore } from '@/stores/search'
import { matchesQuery } from '@/lib/filter'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'surname', label: 'Surname' },
  { key: 'status', label: 'Status' },
  { key: 'notes', label: 'Notes' },
]

const rows = ref<User[]>([])
const search = useSearchStore()
const filteredRows = computed(() =>
  (rows.value as Record<string, unknown>[]).filter((r) => matchesQuery(r, columns, search.query)),
)
const loading = ref(true)

const panelOpen = ref(false)
const panelMode = ref<'create' | 'edit'>('create')
const selectedRow = ref<Record<string, unknown> | null>(null)
const confirmOpen = ref(false)

onMounted(async () => {
  rows.value = await listUsers()
  loading.value = false
})

function onRowClick(row: Record<string, unknown>) {
  selectedRow.value = row
  panelMode.value = 'edit'
  panelOpen.value = true
}

function onAdd() {
  selectedRow.value = null
  panelMode.value = 'create'
  panelOpen.value = true
}

async function onSaved() {
  panelOpen.value = false
  rows.value = await listUsers()
}

function onDeleteClick(row: Record<string, unknown>) {
  selectedRow.value = row
  confirmOpen.value = true
}

async function onDeleteConfirmed() {
  confirmOpen.value = false
  if (selectedRow.value?.id != null) {
    await deleteUser(selectedRow.value.id as number)
    rows.value = await listUsers()
  }
}
</script>
