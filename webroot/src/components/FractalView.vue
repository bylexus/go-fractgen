<script lang="ts" setup>
import { type ElementInfo, useElementResize } from '@/lib/element-info'
import { apiroot, queryStr } from '@/lib/url_helper'
import { onMounted, reactive, type Ref, ref, watch, watchEffect } from 'vue'
import ColorPresetsSelect from './ColorPresetsSelect.vue'
import { useLocalStorageVariable } from '@/lib/use-local-storage'
import FractalPresetsSelect from './FractalPresetsSelect.vue'
import { fractalPresetByName, type FractalPreset } from '@/lib/use-presets'

const imageUrl: Ref<string> = ref('')
const loading = ref(true)
const imgContainer: Ref<HTMLElement | null> = ref(null)
const image: Ref<HTMLImageElement | null> = ref(null)
const windowSizes = useElementResize(window.document.body, 1000)
const colorPreset = useLocalStorageVariable('colorPreset', '')
const fractalPreset = useLocalStorageVariable('fractalPreset', '')

const fractalParams: FractalPreset & { width: number; height: number } = reactive(
  Object.assign(
    {
      width: 0,
      height: 0,
      iterFunc: 'Mandelbrot',
      maxIterations: 40,
      centerCY: 0,
      centerCX: -0.7,
      diameterCX: 4,
      colorPreset: colorPreset.value,
      colorPaletteRepeat: 1,
      juliaKr: 0,
      juliaKi: 0,
    },
    fractalPresetByName(fractalPreset.value || ''),
  ),
)

onMounted(() => {
  loading.value = true
  fractalParams.width = imgContainer.value?.clientWidth || 0
  fractalParams.height = imgContainer.value?.clientHeight || 0
})

watch(windowSizes.sizes, ({ width, height }) => {
  fractalParams.width = width
  fractalParams.height = height
})

watch(colorPreset, () => {
  fractalParams.colorPreset = colorPreset.value || ''
})

watch(fractalPreset, () => {
  const preset = fractalPresetByName(fractalPreset.value)
  if (preset) {
    fractalParams.iterFunc = preset.iterFunc
    fractalParams.maxIterations = preset.maxIterations
    fractalParams.centerCX = preset.centerCX
    fractalParams.centerCY = preset.centerCY
    fractalParams.diameterCX = preset.diameterCX
    fractalParams.colorPreset = preset.colorPreset
    fractalParams.colorPaletteRepeat = preset.colorPaletteRepeat
    colorPreset.value = preset.colorPreset
  }
})

watchEffect(() => {
  calcImage(fractalParams)
})

function afterImageLoad() {
  loading.value = false
  if (image.value) {
    image.value.style.transform = ''
  }
}

function calcImage(fractalParams: any) {
  const apiRoot = apiroot()
  loading.value = true
  // check for the most important values to be present:
  if (!fractalParams.width || !fractalParams.height) return
  if (!fractalParams.colorPreset) return
  fractalParams.ts = new Date().getTime()
  imageUrl.value = `${apiRoot}/fractal-image.png?${queryStr(fractalParams)}`
}

let dragStartPos: { x: number; y: number } | null = null
let dragDistance: { dx: number; dy: number } | null = null
function onDragStart(ev: PointerEvent) {
  ev.preventDefault()
  dragStartPos = { x: ev.clientX, y: ev.clientY }
  dragDistance = { dx: 0, dy: 0 }
}

function onDrag(ev: PointerEvent) {
  ev.preventDefault()
  if (dragStartPos) {
    const dx = ev.clientX - dragStartPos.x
    const dy = ev.clientY - dragStartPos.y
    dragDistance = { dx, dy }
    if (image.value) {
      image.value.style.transform = `translate(${dx}px, ${dy}px)`
    }
  }
}

function onDragEnd(ev: PointerEvent) {
  ev.preventDefault()
  if (!dragDistance || (dragDistance.dx === 0 && dragDistance.dy === 0)) {
    dragStartPos = null
    dragDistance = null
    return
  }
  const aspect = fractalParams.width / fractalParams.height
  const fractDiameterCY = fractalParams.diameterCX / aspect
  const moveFactorX = dragDistance ? dragDistance.dx / fractalParams.width : 0
  const moveFactorY = dragDistance ? dragDistance.dy / fractalParams.height : 0
  const fractDistX = fractalParams.diameterCX * moveFactorX
  const fractDistY = fractDiameterCY * moveFactorY
  fractalParams.centerCX -= fractDistX
  fractalParams.centerCY += fractDistY
  calcImage(fractalParams)
  dragStartPos = null
  dragDistance = null
}
</script>

<template>
  <div class="display-container" ref="imgContainer">
    <img
      ref="image"
      :src="imageUrl"
      alt="Fractal Image"
      style="touch-action: none;"
      @load="afterImageLoad"
      @pointerdown="onDragStart"
      @pointermove="onDrag"
      @pointerup="onDragEnd"
    />
    <div class="settings-overlay">
      <FractalPresetsSelect v-model="fractalPreset"></FractalPresetsSelect>
      <ColorPresetsSelect v-model="colorPreset"></ColorPresetsSelect>
      <div class="label-field">
        <label for="paletteRepeat">Palette Repeat:</label>
        <input type="number" v-model.lazy="fractalParams.colorPaletteRepeat" id="paletteRepeat" />
      </div>
    </div>
    <div v-if="loading" class="loading-overlay">Calculating...</div>
  </div>
</template>

<style scoped>
.display-container {
  width: 100%;
  height: 100%;
}

.loading-overlay {
  position: absolute;
  z-index: 2;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  color: white;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.5);
}

.settings-overlay {
  position: absolute;
  z-index: 1;
  top: 0;
  left: 0;
  width: 100%;
}

.label-field {
  display: inline-flex;
  flex-direction: column;
  label {
    font-size: 0.875rem;
  }
}
</style>
