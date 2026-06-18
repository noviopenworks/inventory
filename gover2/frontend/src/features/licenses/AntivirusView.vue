<template>
  <div class="p-4 text-text-primary">
    <h1 class="text-lg font-semibold mb-4">Antivirus</h1>
    <DataTable :columns="[]" :rows="[]" :loading="false" />
  </div>
</template>

<script setup lang="ts">
import DataTable from '@/components/DataTable.vue'
</script>
