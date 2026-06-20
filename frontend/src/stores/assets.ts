import { defineStore } from 'pinia'
import { ref } from 'vue'
import type {
  Antivirus,
  Computer,
  OtherSoftware,
  Smartphone,
  Tablet,
  User,
  WindowsKey,
} from '@/lib/types'

export type AssetCategory =
  | 'computers'
  | 'smartphones'
  | 'tablets'
  | 'all'
  | 'windows-keys'
  | 'antivirus'
  | 'other-software'
  | 'users'

export type AssetItem =
  | Computer
  | Smartphone
  | Tablet
  | WindowsKey
  | Antivirus
  | OtherSoftware
  | User

export const useAssetsStore = defineStore('assets', () => {
  const currentCategory = ref<AssetCategory>('computers')
  const items = ref<AssetItem[]>([])
  const selectedId = ref<number | null>(null)
  const searchQuery = ref('')
  const statusFilter = ref('')

  function setCategory(category: AssetCategory) {
    currentCategory.value = category
    selectedId.value = null
    items.value = []
  }

  function setItems(newItems: AssetItem[]) {
    items.value = newItems
  }

  function setSelectedId(id: number | null) {
    selectedId.value = id
  }

  return {
    currentCategory,
    items,
    selectedId,
    searchQuery,
    statusFilter,
    setCategory,
    setItems,
    setSelectedId,
  }
})
