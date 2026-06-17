<template>
  <router-link
    :to="to"
    class="block px-3 py-1.5 rounded text-sm text-sidebar-text hover:bg-sidebar-hover hover:text-sidebar-hi transition-colors"
    active-class="bg-sidebar-active text-sidebar-hi"
  >
    {{ label }}
  </router-link>
</template>

<script setup lang="ts">
defineProps<{
  to: string
  label: string
}>()
</script>
