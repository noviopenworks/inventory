import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const density = ref<'comfortable' | 'compact'>('comfortable')
  const darkMode = ref(false)
  const sidebarCollapsed = ref(false)

  function toggleDarkMode() {
    darkMode.value = !darkMode.value
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setDensity(d: 'comfortable' | 'compact') {
    density.value = d
  }

  return {
    density,
    darkMode,
    sidebarCollapsed,
    toggleDarkMode,
    toggleSidebar,
    setDensity,
  }
})
