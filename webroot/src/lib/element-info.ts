import { computed, type Ref, ref } from 'vue'

export type ElementInfo = {
  width: Ref<number>
  height: Ref<number>
  sizes: Ref<{ width: number; height: number }>
}
export function useElementResize(el: HTMLElement, delayMs: number = 100): ElementInfo {
  const width = ref(el.clientWidth)
  const height = ref(el.clientHeight)
  let timeoutRef = 0
  const observer = new ResizeObserver(() => {
    window.clearTimeout(timeoutRef)
    timeoutRef = window.setTimeout(() => {
      width.value = el.clientWidth
      height.value = el.clientHeight
    }, delayMs)
  })
  observer.observe(el)
  return {
    width,
    height,
    sizes: computed(() => ({ width: width.value, height: height.value })),
  }
}

export function screenSize(): {
  width: number
  height: number
  physicalWidth: number
  physicalHeight: number
  devicePixelRatio: number
} {
  const ratio = window.devicePixelRatio || 1
  return {
    width: window.screen?.width || 0,
    height: window.screen?.height || 0,
    physicalWidth: Math.floor((window.screen?.width || 0) * ratio),
    physicalHeight: Math.floor((window.screen?.height || 0) * ratio),
    devicePixelRatio: window.devicePixelRatio || 1,
  }
}
