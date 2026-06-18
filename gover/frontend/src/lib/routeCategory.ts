export const routeCategory: Record<string, string> = {
  '/computers': 'computers',
  '/smartphones': 'smartphones',
  '/tablets': 'tablets',
  '/windows-keys': 'windowskeys',
  '/antivirus': 'antivirus',
  '/other-software': 'othersoftware',
  '/users': 'users',
}

export function categoryForRoute(path: string): string | null {
  return routeCategory[path] ?? null
}
