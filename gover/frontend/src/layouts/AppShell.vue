<template>
  <div class="flex h-screen bg-page font-ui">
    <Sidebar />
    <div class="flex flex-col flex-1 overflow-hidden">
      <Topbar
        @open-alerts="onOpenAlerts"
        @open-about="aboutOpen = true"
        @open-license="licenseOpen = true"
        @db-changed="onDbChanged"
      />
      <main class="flex-1 overflow-auto p-4">
        <slot />
      </main>
      <Statusbar />
    </div>

    <AlertsModal
      :open="alertsOpen"
      :alerts="alerts"
      :warning-days="ui.expiryWarningDays"
      @close="alertsOpen = false"
    />
    <AboutModal :open="aboutOpen" @close="aboutOpen = false" />
    <LicenseModal :open="licenseOpen" @close="licenseOpen = false" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'
import Statusbar from '@/components/Statusbar.vue'
import AlertsModal from '@/components/AlertsModal.vue'
import AboutModal from '@/components/AboutModal.vue'
import LicenseModal from '@/components/LicenseModal.vue'
import { useUiStore } from '@/stores/ui'
import { useSearchStore } from '@/stores/search'
import { getAlerts } from '@/lib/api'
import type { Alert } from '@/lib/api'

const ui = useUiStore()
const search = useSearchStore()
const route = useRoute()

const alerts = ref<Alert[]>([])
const alertsOpen = ref(false)
const aboutOpen = ref(false)
const licenseOpen = ref(false)

onMounted(async () => {
  // Theme/density are loaded pre-mount in main.ts; here we only handle alerts.
  alerts.value = await getAlerts()
  if (alerts.value.length > 0) alertsOpen.value = true
})

async function onOpenAlerts() {
  alerts.value = await getAlerts()
  alertsOpen.value = true
}

async function onDbChanged() {
  // DB switched: remount the active view (via the dbVersion :key on <router-view>)
  // so it refetches rows from the new database, and refresh alerts to match.
  ui.bumpDbVersion()
  alerts.value = await getAlerts()
}

watch(() => route.path, () => search.clear())
</script>
