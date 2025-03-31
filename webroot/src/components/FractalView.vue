<script lang="ts" setup>
import { useElementResize } from '@/lib/element-info'
import { apiroot, queryStr } from '@/lib/url_helper'
import { onMounted, reactive, type Ref, ref, watch, watchEffect } from 'vue'
import ColorPresetsSelect from './ColorPresetsSelect.vue'
import { useLocalStorageVariable } from '@/lib/use-local-storage'
import FractalPresetsSelect from './FractalPresetsSelect.vue'
import { fractalPresetByName, useFractalPresets, type FractalPreset } from '@/lib/use-presets'
import { usePointer } from '@/lib/use-pointer'

const imageUrl: Ref<string> = ref('')
const loading = ref(true)
const imgContainer: Ref<HTMLElement | null> = ref(null)
const image: Ref<HTMLImageElement | null> = ref(null)
const windowSizes = useElementResize(window.document.body, 1000)
const colorPreset = ref('')
const fractalPreset = ref('')

const fractalParams: Ref<FractalPreset & { width: number; height: number }> =
  useLocalStorageVariable(
    'fractalParams',
    Object.assign(
      {
        width: 0, // will be updated on mount, based on the container size
        height: 0,
      },
      // we use the first preset as default:
      { ...useFractalPresets().presets.value[0] },
    ),
  )

// Initial values:
colorPreset.value = fractalParams.value.colorPreset || ''
fractalPreset.value = fractalParams.value.name || ''

onMounted(() => {
  loading.value = true
  fractalParams.value.width = imgContainer.value?.clientWidth || 0
  fractalParams.value.height = imgContainer.value?.clientHeight || 0

  let imagePointerHandler = usePointer(image.value!)

  imagePointerHandler.onDoubleClick((ev: PointerEvent) => {
    ev.preventDefault()
    centerZoom(ev, 2)
  })
  // on drag: move the image
  watch(imagePointerHandler.pointerMoveWhileDragging, ({ dx, dy }) => {
    if (image.value) {
      image.value.style.transform = `translate(${dx}px, ${dy}px)`
    }
  })
  imagePointerHandler.onDragEnd((ev: PointerEvent, { dx, dy }: { dx: number; dy: number }) => {
    ev.preventDefault()
    onDragEnd({ dx, dy })
  })
})

watch(windowSizes.sizes, ({ width, height }) => {
  fractalParams.value.width = width
  fractalParams.value.height = height
})

watch(colorPreset, () => {
  fractalParams.value.colorPreset = colorPreset.value || ''
})

watch(fractalPreset, () => {
  const preset = fractalPresetByName(fractalPreset.value)
  if (preset) {
    fractalParams.value.iterFunc = preset.iterFunc
    fractalParams.value.maxIterations = preset.maxIterations
    fractalParams.value.centerCX = preset.centerCX
    fractalParams.value.centerCY = preset.centerCY
    fractalParams.value.diameterCX = preset.diameterCX
    fractalParams.value.colorPreset = preset.colorPreset
    fractalParams.value.colorPaletteRepeat = preset.colorPaletteRepeat
    fractalParams.value.name = preset.name || ''
    colorPreset.value = preset.colorPreset
  }
})

watchEffect(() => {
  calcImage(fractalParams.value)
})

function recalcIterations(diameterCX: number) {
  // approximation of the number of iterations, based on the following formula,
  // which seems to work well:
  // maxIter = 50 * (log10(scale))^1.25
  // where scale is pixelWidth/complexPlaneWidth e.g. 1280/5
  return Math.ceil(50 * Math.pow(Math.log10(fractalParams.value.width / diameterCX), 1.25))
}

function afterImageLoad() {
  loading.value = false
  if (image.value) {
    image.value.style.transform = ''
  }
}

function zoomIn() {
  image.value!.style.transform = 'scale(2)'
  fractalParams.value.diameterCX /= 2.0
  fractalParams.value.maxIterations = recalcIterations(fractalParams.value.diameterCX)
}

function zoomOut() {
  image.value!.style.transform = 'scale(0.5)'
  fractalParams.value.diameterCX *= 2.0
  fractalParams.value.maxIterations = recalcIterations(fractalParams.value.diameterCX)
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

function onDragEnd({ dx, dy }: { dx: number; dy: number }) {
  if (dx === 0 && dy === 0) {
    return
  }
  const aspect = fractalParams.value.width / fractalParams.value.height
  const fractDiameterCY = fractalParams.value.diameterCX / aspect
  const moveFactorX = dx / fractalParams.value.width
  const moveFactorY = dy / fractalParams.value.height
  const fractDistX = fractalParams.value.diameterCX * moveFactorX
  const fractDistY = fractDiameterCY * moveFactorY
  fractalParams.value.centerCX -= fractDistX
  fractalParams.value.centerCY += fractDistY
  calcImage(fractalParams)
}

function centerZoom(ev: PointerEvent, scale: number) {
  const imgRect = image.value!.getBoundingClientRect()
  // calc the distance from the clicked coord to the center:
  const dX = ev.clientX - imgRect.left - fractalParams.value.width / 2.0
  const dY = ev.clientY - imgRect.top - fractalParams.value.height / 2.0

  const aspect = fractalParams.value.width / fractalParams.value.height
  const fractDiameterCY = fractalParams.value.diameterCX / aspect
  const moveFactorX = dX / fractalParams.value.width
  const moveFactorY = dY / fractalParams.value.height
  const fractDistX = fractalParams.value.diameterCX * moveFactorX
  const fractDistY = fractDiameterCY * moveFactorY
  fractalParams.value.centerCX += fractDistX
  fractalParams.value.centerCY -= fractDistY
  fractalParams.value.diameterCX /= scale
  fractalParams.value.maxIterations = recalcIterations(fractalParams.value.diameterCX)

  image.value!.style.transform = `scale(${scale}) translate(${-1*dX}px, ${-1*dY}px)`

  calcImage(fractalParams)
}
</script>

<template>
  <div class="display-container" ref="imgContainer">
    <img
      ref="image"
      :src="imageUrl"
      alt="Fractal Image"
      style="touch-action: none"
      @load="afterImageLoad"
    />
    <div class="settings-overlay">
      <FractalPresetsSelect v-model="fractalPreset"></FractalPresetsSelect>
      <ColorPresetsSelect v-model="colorPreset"></ColorPresetsSelect>
      <div class="label-field">
        <label for="iterations">Max. Iterations</label>
        <input type="number" v-model.lazy="fractalParams.maxIterations" id="iterations" />
      </div>
      <div class="label-field">
        <label for="paletteRepeat">Palette Repeat:</label>
        <input type="number" v-model.lazy="fractalParams.colorPaletteRepeat" id="paletteRepeat" />
      </div>
      <button type="button" @click="zoomIn">+</button>
      <button type="button" @click="zoomOut">-</button>
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
