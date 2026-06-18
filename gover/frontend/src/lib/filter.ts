export function matchesQuery(
  row: Record<string, unknown>,
  columns: { key: string }[],
  query: string,
): boolean {
  const q = query.trim().toLowerCase()
  if (q === '') return true
  return columns.some((col) => {
    const value = row[col.key] ?? ''
    return String(value).toLowerCase().includes(q)
  })
}
