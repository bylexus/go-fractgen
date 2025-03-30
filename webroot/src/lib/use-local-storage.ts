import { readonly, ref, watch, type Ref } from 'vue'

export function useLocalStorageVariable<T>(key: string, defaultValue: T) {
  const value: Ref<T> = ref(JSON.parse(localStorage.getItem(key) || '""') || defaultValue)

  watch(
    value,
    () => {
      localStorage.setItem(key, JSON.stringify(value.value))
    },
    { deep: true },
  )

  return value
}
