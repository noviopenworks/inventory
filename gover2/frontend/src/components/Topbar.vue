<template>
  <header class="flex items-center px-4 bg-surface border-b border-border" style="height: 48px">
    <span class="text-sm text-text-secondary">{{ title }}</span>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const title = computed(() =>
  String(route.path).replace('/', '').replace(/-/g, ' ') || 'Inventory'
)
</script>
