import {
  addComputer, updateComputer,
  addSmartphone, updateSmartphone,
  addTablet, updateTablet,
  addWindowsKey, updateWindowsKey,
  addAntivirus, updateAntivirus,
  addOtherSoftware, updateOtherSoftware,
  addUser, updateUser,
} from '@/lib/api'

export async function callApi(category: string, isEdit: boolean, id: number, payload: unknown): Promise<void> {
  const cat = category

  if (cat === 'computers') {
    return isEdit
      ? updateComputer(id, payload as Parameters<typeof updateComputer>[1])
      : addComputer(payload as Parameters<typeof addComputer>[0])
  }
  if (cat === 'smartphones') {
    return isEdit
      ? updateSmartphone(id, payload as Parameters<typeof updateSmartphone>[1])
      : addSmartphone(payload as Parameters<typeof addSmartphone>[0])
  }
  if (cat === 'tablets') {
    return isEdit
      ? updateTablet(id, payload as Parameters<typeof updateTablet>[1])
      : addTablet(payload as Parameters<typeof addTablet>[0])
  }
  if (cat === 'windowskeys') {
    return isEdit
      ? updateWindowsKey(id, payload as Parameters<typeof updateWindowsKey>[1])
      : addWindowsKey(payload as Parameters<typeof addWindowsKey>[0])
  }
  if (cat === 'antivirus') {
    return isEdit
      ? updateAntivirus(id, payload as Parameters<typeof updateAntivirus>[1])
      : addAntivirus(payload as Parameters<typeof addAntivirus>[0])
  }
  if (cat === 'othersoftware') {
    return isEdit
      ? updateOtherSoftware(id, payload as Parameters<typeof updateOtherSoftware>[1])
      : addOtherSoftware(payload as Parameters<typeof addOtherSoftware>[0])
  }
  if (cat === 'users') {
    return isEdit
      ? updateUser(id, payload as Parameters<typeof updateUser>[1])
      : addUser(payload as Parameters<typeof addUser>[0])
  }
}
