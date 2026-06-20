<template>
  <AppShell>
    <router-view :key="ui.dbVersion" />
  </AppShell>
</template>

<script setup lang="ts">
import AppShell from '@/layouts/AppShell.vue'
import { useUiStore } from '@/stores/ui'

const ui = useUiStore()
</script>
