import { ref, type Ref } from 'vue'
import { apiroot } from './url_helper'

export type ColorPreset = {
  name: string
  ident: string
  colors: Array<{
    R: number
    G: number
    B: number
    A: number
    steps?: number
  }>
}

export type FractalPreset = {
  [key: string]: any
  name?: string
  iterFunc: 'Mandelbrot' | 'Julia' | 'Mandelbrot3' | 'Mandelbrot4'
  maxIterations: number
  centerCX: number
  centerCY: number
  diameterCX: number
  colorPreset: string
  colorPaletteRepeat: number
  colorPaletteLength: number
  colorPaletteReverse: boolean
  juliaKr: number
  juliaKi: number
}

export type FractalParams = FractalPreset & { width?: number; height?: number }

const isLoading = ref(false)
let loaded = false
const colorPresets: Ref<Array<ColorPreset>> = ref([])
const fractalPresets: Ref<Array<FractalPreset>> = ref([])

export async function loadPresets() {
  if (loaded || isLoading.value) {
    return
  }
  isLoading.value = true
  try {
    const data = await (await fetch(`${apiroot()}/presets.json`)).json()
    colorPresets.value = data.colorPresets
    fractalPresets.value = data.fractalPresets
    loaded = true
    isLoading.value = false
  } finally {
    isLoading.value = false
  }
}

export function useColorPresets() {
  loadPresets()
  return { loading: isLoading, presets: colorPresets }
}

export function useFractalPresets() {
  loadPresets()
  return { loading: isLoading, presets: fractalPresets }
}

export function fractalPresetByName(name: string): FractalPreset | null {
  for (const preset of fractalPresets.value) {
    if (preset.name === name) {
      return preset
    }
  }
  return null
}
