import { computed, type Ref, ref } from 'vue'

export type ElementInfo = {
  width: Ref<number>
  height: Ref<number>
  sizes: Ref<{ width: number; height: number }>
}
export function useElementResize(el: HTMLElement, delayMs: number = 0): ElementInfo {
  const width = ref(el.clientWidth)
  const height = ref(el.clientHeight)
  let timeoutRef = 0
  el.onresize = (ev: UIEvent) => {
    window.clearTimeout(timeoutRef)
    timeoutRef = window.setTimeout(() => {
      width.value = el.clientWidth
      height.value = el.clientHeight
    }, delayMs)
  }
  return {
    width,
    height,
    sizes: computed(() => ({ width: width.value, height: height.value })),
  }
}
