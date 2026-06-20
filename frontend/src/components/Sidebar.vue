<template>
  <nav
    class="flex flex-col bg-sidebar border-r border-sidebar-border"
    style="width: 200px"
  >
    <div class="p-4 text-sidebar-hi font-semibold text-sm border-b border-sidebar-border">
      Inventory
    </div>
    <div class="flex flex-col gap-0.5 p-2 flex-1">
      <SidebarItem to="/all" label="All Assets" />
      <SidebarItem to="/computers" label="Computers" />
      <SidebarItem to="/smartphones" label="Smartphones" />
      <SidebarItem to="/tablets" label="Tablets" />
      <div class="my-2 border-t border-sidebar-border" />
      <SidebarItem to="/windows-keys" label="Windows Keys" />
      <SidebarItem to="/antivirus" label="Antivirus" />
      <SidebarItem to="/other-software" label="Other Software" />
      <div class="my-2 border-t border-sidebar-border" />
      <SidebarItem to="/users" label="Users" />
    </div>
  </nav>
</template>

<script setup lang="ts">
import SidebarItem from '@/components/SidebarItem.vue'
</script>
