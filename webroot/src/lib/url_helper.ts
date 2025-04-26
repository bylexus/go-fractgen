export function apiroot(): string {
  return import.meta.env.VITE_API_ROOT
}
export function queryStr(obj: { [key: string]: any }) {
  const parts = []
  for (const key in obj) {
    const value = obj[key]
    if (value instanceof Object) {
      parts.push(`${key}=${encodeURIComponent(JSON.stringify(value))}`)
    } else {
      parts.push(`${key}=${encodeURIComponent(value)}`)
    }
  }
  if (parts.length > 0) {
    return `${parts.join('&')}`
  }
}
