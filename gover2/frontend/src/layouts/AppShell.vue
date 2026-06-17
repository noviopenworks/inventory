<template>
  <div class="flex h-screen bg-page font-ui">
    <Sidebar />
    <div class="flex flex-col flex-1 overflow-hidden">
      <Topbar />
      <main class="flex-1 overflow-auto p-4">
        <slot />
      </main>
      <Statusbar />
    </div>
  </div>
</template>

<script setup lang="ts">
import Sidebar from '@/components/Sidebar.vue'
import Topbar from '@/components/Topbar.vue'
import Statusbar from '@/components/Statusbar.vue'
</script>
