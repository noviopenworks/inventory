<template>
  <div class="p-4 text-text-primary">
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-lg font-semibold">Computers</h1>
      <button @click="onAdd"
              class="px-3 py-1.5 bg-accent text-white text-sm rounded hover:opacity-90">
        Add
      </button>
    </div>
    <DataTable
      :columns="columns"
      :rows="rows as Record<string, unknown>[]"
      :loading="loading"
      @row-click="onRowClick"
      @delete="onDeleteClick"
    />
    <EditPanel
      v-if="panelOpen"
      :category="'computers'"
      :mode="panelMode"
      :row="selectedRow"
      :users="users"
      :devices="devices"
      @saved="onSaved"
      @cancelled="panelOpen = false"
    />
    <ConfirmDialog
      v-if="confirmOpen"
      message="Delete this computer?"
      @confirmed="onDeleteConfirmed"
      @cancelled="confirmOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '@/components/DataTable.vue'
import EditPanel from '@/components/EditPanel.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { listComputers, deleteComputer, listUsersForDropdown, listDevicesForDropdown } from '@/lib/api'
import type { Computer, DropdownItem, DeviceDropdownItem } from '@/lib/api'

const columns = [
  { key: 'name', label: 'Name' },
  { key: 'model', label: 'Model' },
  { key: 'status', label: 'Status' },
  { key: 'purchaseDate', label: 'Purchase Date' },
  { key: 'warrantyExpiry', label: 'Warranty Expiry' },
]

const rows = ref<Computer[]>([])
const loading = ref(true)

const panelOpen = ref(false)
const panelMode = ref<'create' | 'edit'>('create')
const selectedRow = ref<Record<string, unknown> | null>(null)
const confirmOpen = ref(false)
const users = ref<DropdownItem[]>([])
const devices = ref<DeviceDropdownItem[]>([])

onMounted(async () => {
  ;[rows.value, users.value, devices.value] = await Promise.all([
    listComputers(), listUsersForDropdown(), listDevicesForDropdown(),
  ])
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
  rows.value = await listComputers()
}

function onDeleteClick(row: Record<string, unknown>) {
  selectedRow.value = row
  confirmOpen.value = true
}

async function onDeleteConfirmed() {
  confirmOpen.value = false
  if (selectedRow.value?.id != null) {
    await deleteComputer(selectedRow.value.id as number)
    rows.value = await listComputers()
  }
}
</script>
