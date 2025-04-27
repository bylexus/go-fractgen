import { readonly, ref, watch, type Ref } from 'vue'

export function useSessionStorageVariable<T>(key: string, defaultValue: T) {
  const value: Ref<T> = ref(JSON.parse(sessionStorage.getItem(key) || '""') || defaultValue)

  watch(
    value,
    () => {
      sessionStorage.setItem(key, JSON.stringify(value.value))
    },
    { deep: true },
  )

  return value
}
