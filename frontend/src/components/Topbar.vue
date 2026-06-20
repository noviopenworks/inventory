<template>
  <header class="flex items-center gap-3 px-4 bg-surface border-b border-border" style="height: 48px">
    <span class="text-sm text-text-secondary">{{ title }}</span>

    <input
      v-model="search.query"
      type="search"
      placeholder="Search..."
      class="ml-2 px-2 py-1 text-sm bg-page text-text-primary border border-border rounded w-56"
    />

    <div class="ml-auto flex items-center gap-2">
      <button
        class="px-2 py-1 text-sm text-text-primary hover:bg-page rounded"
        title="Alerts"
        @click="$emit('open-alerts')"
      >
        ⚠ Alerts
      </button>

      <button
        v-if="exportCategory"
        class="px-2 py-1 text-sm text-text-primary hover:bg-page rounded"
        title="Export CSV"
        @click="onExport"
      >
        ↓ CSV
      </button>

      <button
        class="px-2 py-1 text-sm text-text-primary hover:bg-page rounded"
        title="Toggle dark mode"
        @click="ui.toggleDarkMode()"
      >
        {{ ui.darkMode ? '☀' : '🌙' }}
      </button>

      <div class="relative">
        <button
          class="px-2 py-1 text-sm text-text-primary hover:bg-page rounded"
          title="More"
          @click="menuOpen = !menuOpen"
        >
          ⋯
        </button>
        <div
          v-if="menuOpen"
          class="absolute right-0 mt-1 w-44 bg-surface border border-border shadow z-50 text-sm"
        >
          <button
            data-menu-item
            class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page"
            @click="onNewDatabase"
          >
            New Database
          </button>
          <button
            data-menu-item
            class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page"
            @click="onOpenDatabase"
          >
            Open Database
          </button>
          <button
            data-menu-item
            data-testid="backup-database"
            class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page"
            @click="onBackupDatabase"
          >
            Backup database…
          </button>
          <div class="border-t border-border"></div>
          <button
            data-menu-item
            class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page"
            @click="onAbout"
          >
            About
          </button>
          <button
            data-menu-item
            class="block w-full text-left px-3 py-2 text-text-primary hover:bg-page"
            @click="onLicense"
          >
            License
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useSearchStore } from '@/stores/search'
import { useUiStore } from '@/stores/ui'
import { categoryForRoute } from '@/lib/routeCategory'
import { backupDatabaseDialog, exportCSV, newDatabaseDialog, openDatabaseDialog } from '@/lib/api'

const emit = defineEmits<{
  (e: 'open-alerts'): void
  (e: 'open-about'): void
  (e: 'open-license'): void
  (e: 'db-changed'): void
}>()

const route = useRoute()
const search = useSearchStore()
const ui = useUiStore()
const menuOpen = ref(false)

const title = computed(() =>
  String(route.path).replace('/', '').replace(/-/g, ' ') || 'Inventory',
)
const exportCategory = computed(() => categoryForRoute(route.path))

async function onExport() {
  if (exportCategory.value) {
    await exportCSV(exportCategory.value)
  }
}

async function onNewDatabase() {
  menuOpen.value = false
  await newDatabaseDialog()
  emit('db-changed')
}

async function onOpenDatabase() {
  menuOpen.value = false
  await openDatabaseDialog()
  emit('db-changed')
}

async function onBackupDatabase() {
  menuOpen.value = false
  await backupDatabaseDialog()
}

function onAbout() {
  menuOpen.value = false
  emit('open-about')
}

function onLicense() {
  menuOpen.value = false
  emit('open-license')
}
</script>
