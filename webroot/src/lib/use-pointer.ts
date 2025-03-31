import { readonly, ref, type Ref } from 'vue'

export function usePointer(el: HTMLElement) {
  const singleClickListeners: ((ev: PointerEvent) => void)[] = []
  const doubleClickListeners: ((ev: PointerEvent) => void)[] = []
  const dragEndListeners: ((ev: PointerEvent, { dx, dy }: { dx: number; dy: number }) => void)[] =
    []

  let pointerDownPos: { x: number; y: number } = { x: 0, y: 0 }
  let pointerUpCount = 0
  let pointerIsDown = false
  const pointerMove: Ref<{ dx: number; dy: number }> = ref({ dx: 0, dy: 0 })

  el.addEventListener('pointerdown', (ev: PointerEvent) => {
    ev.preventDefault()
    pointerDownPos = { x: ev.clientX, y: ev.clientY }
    pointerMove.value = { dx: 0, dy: 0 }
    pointerIsDown = true
  })
  el.addEventListener('pointerup', (ev: PointerEvent) => {
    ev.preventDefault()
    pointerIsDown = false
    const pointerUpPos = { x: ev.clientX, y: ev.clientY }
    if (pointerUpPos.x !== pointerDownPos?.x || pointerUpPos.y !== pointerDownPos?.y) {
      pointerUpCount = 0
      dragEndListeners.forEach((cb) => cb(ev, pointerMove.value))
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
  el.addEventListener('pointermove', (ev: PointerEvent) => {
    if (pointerIsDown) {
      ev.preventDefault()
      pointerMove.value = { dx: ev.clientX - pointerDownPos.x, dy: ev.clientY - pointerDownPos.y }
    }
  })

  return {
    pointerMoveWhileDragging: readonly(pointerMove),

    onDragEnd: (callback: (ev: PointerEvent, { dx, dy }: { dx: number; dy: number }) => void) => {
      dragEndListeners.push(callback)
    },

    onSingleClick: (callback: (ev: PointerEvent) => void) => {
      singleClickListeners.push(callback)
    },
    onDoubleClick: (callback: (ev: PointerEvent) => void) => {
      doubleClickListeners.push(callback)
    },
  }
}
