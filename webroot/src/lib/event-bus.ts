export enum GlobalEvent {
  hideHud,
  showHud,
}

const listeners = new Map<GlobalEvent, ((data?: any) => void)[]>()

export function onEvent<T>(event: GlobalEvent, cb: (data?: T) => void) {
  if (listeners.has(event)) {
    const existing = listeners.get(event)!
    if (!existing.includes(cb)) {
      existing.push(cb)
    }
  } else {
    listeners.set(event, [cb])
  }
}

export function emitEvent(event: GlobalEvent, data?: any) {
  if (listeners.has(event)) {
    listeners.get(event)!.forEach((cb) => cb(data))
  }
}

export function offEvent(event: GlobalEvent, cb: (data?: any) => void) {
  if (listeners.has(event)) {
    const existing = listeners.get(event)!
    if (existing.includes(cb)) {
      listeners.set(
        event,
        existing.filter((c) => c !== cb),
      )
    }
  }
}
