export function usePointer(el: HTMLElement) {
  const singleClickListeners: ((ev: PointerEvent) => void)[] = []
  const doubleClickListeners: ((ev: PointerEvent) => void)[] = []
  let pointerDownPos: { x: number; y: number } | null = null
  let pointerUpCount = 0

  el.addEventListener('pointerdown', (ev: PointerEvent) => {
    pointerDownPos = { x: ev.clientX, y: ev.clientY }
  })
  el.addEventListener('pointerup', (ev: PointerEvent) => {
    const pointerUpPos = { x: ev.clientX, y: ev.clientY }
    if (pointerUpPos.x !== pointerDownPos?.x || pointerUpPos.y !== pointerDownPos?.y) {
      pointerUpCount = 0
      return
    }
    pointerUpCount++
    setTimeout(() => {
      if (pointerUpCount === 1) {
        pointerUpCount = 0
        singleClickListeners.forEach((cb) => cb(ev))
      }
    }, 250)
    if (pointerUpCount === 2) {
      pointerUpCount = 0
      doubleClickListeners.forEach((cb) => cb(ev))
    }
  })
  el.addEventListener('pointermove', (ev: PointerEvent) => {})

  return {
    onSingleClick: (callback: (ev: PointerEvent) => void) => {
      singleClickListeners.push(callback)
    },
    onDoubleClick: (callback: (ev: PointerEvent) => void) => {
      doubleClickListeners.push(callback)
    },
  }
}
