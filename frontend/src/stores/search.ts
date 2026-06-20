import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSearchStore = defineStore('search', () => {
  const query = ref('')
  function setQuery(q: string) {
    query.value = q
  }
  function clear() {
    query.value = ''
  }
  return { query, setQuery, clear }
})
