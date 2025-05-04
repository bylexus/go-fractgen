<script setup lang="ts">
import { useFractalPresets, type FractalParams } from '@/lib/use-presets'
import Dialog from './Dialog.vue'
import { computed, reactive } from 'vue'
import { apiroot, queryStr } from '@/lib/url_helper'
import { screenSize } from '@/lib/element-info'
import ColorPaletteDisplay from './ColorPaletteDisplay.vue'
const props = defineProps<{
  fractParams: FractalParams
}>()

const emit = defineEmits(['update:fractParams'])

const { physicalWidth, physicalHeight } = screenSize()

const state = reactive({
  imgWidth: physicalWidth || props.fractParams.width,
  imgHeight: physicalHeight || props.fractParams.height,
})

const isSecureContext = computed(() => {
  return window.isSecureContext
})

function fractParamsAsQueryParams(inputObj: { [key: string]: any }) {
  return queryStr({ ...inputObj })
}

function changeFractalParams(params: Partial<FractalParams>) {
  emit('update:fractParams', { ...props.fractParams, ...params })
}

const jpegImageLink = computed(() => {
  const params = { ...props.fractParams }
  params.width = state.imgWidth
  params.height = state.imgHeight
  return `${apiroot()}/fractal-image/jpg?${fractParamsAsQueryParams(params)}`
})
const pngImageLink = computed(() => {
  const params = { ...props.fractParams }
  params.width = state.imgWidth
  params.height = state.imgHeight
  return `${apiroot()}/fractal-image/png?${fractParamsAsQueryParams(props.fractParams)}`
})

function copyParamsAsJson() {
  window.navigator.clipboard.writeText(JSON.stringify(props.fractParams, null, 2))
}

function onResetBtn() {
  const firstPreset = useFractalPresets().presets.value[0]
  changeFractalParams(firstPreset)
}
</script>

<template>
  <Dialog>
    <div class="settings-dialog">
      <fieldset>
        <legend>Fractal Parameters</legend>
        <div class="flex gap-1">
          <div class="label-field">
            <label for="fractFunc">Fractal Function</label>
            <select
              id="fractFunc"
              :value="fractParams.iterFunc"
              @change="
                changeFractalParams({
                  iterFunc: ($event.target as HTMLSelectElement).value as FractalParams['iterFunc'],
                })
              "
            >
              <option value="Mandelbrot">Mandelbrot</option>
              <option value="Mandelbrot3">Mandelbrot ^ 3</option>
              <option value="Mandelbrot4">Mandelbrot ^ 4</option>
              <option value="Julia">Julia</option>
            </select>
          </div>
          <div v-if="fractParams.iterFunc === 'Julia'" class="label-field">
            <label for="juliaKr">Julia: K(r)</label>
            <input
              type="number"
              id="juliaKr"
              :value="fractParams.juliaKr"
              @change.lazy="
                changeFractalParams({ juliaKr: Number(($event.target as HTMLInputElement).value) })
              "
            />
          </div>
          <div v-if="fractParams.iterFunc === 'Julia'" class="label-field">
            <label for="juliaKi">Julia: K(i)</label>
            <input
              type="number"
              id="juliaKi"
              :value="fractParams.juliaKi"
              @change.lazy="
                changeFractalParams({ juliaKi: Number(($event.target as HTMLInputElement).value) })
              "
            />
          </div>
        </div>

        <div class="flex flex-wrap gap-1">
          <div class="label-field">
            <label for="centerCX">Center CX(r)</label>
            <input
              type="number"
              id="centerCX"
              :value="fractParams.centerCX"
              @change.lazy="
                changeFractalParams({ centerCX: Number(($event.target as HTMLInputElement).value) })
              "
            />
          </div>
          <div class="label-field">
            <label for="centerCY">Center CY(i)</label>
            <input
              type="number"
              id="centerCY"
              :value="fractParams.centerCY"
              @change.lazy="
                changeFractalParams({ centerCY: Number(($event.target as HTMLInputElement).value) })
              "
            />
          </div>
          <div class="label-field">
            <label for="diameterCX">Diameter CX(r)</label>
            <input
              type="number"
              id="diameterCX"
              :value="fractParams.diameterCX"
              @change.lazy="
                changeFractalParams({
                  diameterCX: Number(($event.target as HTMLInputElement).value),
                })
              "
            />
          </div>
        </div>
        <div class="label-field">
          <label>Actual Color Palette</label>
          <ColorPaletteDisplay
            :color-preset="fractParams.colorPreset"
            :palette-length="fractParams.colorPaletteLength"
            :palette-repeat="fractParams.colorPaletteRepeat"
            :max-iterations="fractParams.maxIterations"
          >
          </ColorPaletteDisplay>
        </div>

        <div>
          <button type="button" @click="onResetBtn">Reset</button>
        </div>
      </fieldset>
      <fieldset>
        <legend>Export Image</legend>
        <div class="label-field">
          <label for="imgWidth">Image Width (px)</label>
          <input type="number" v-model="state.imgWidth" id="imgWidth" />
        </div>
        <div class="label-field">
          <label for="imgHeight">Image Height (px)</label>
          <input type="number" v-model="state.imgHeight" id="imgHeight" />
        </div>
        <div>
          Image URLs:
          <ul>
            <li>
              <a :href="jpegImageLink" target="_blank" rel="noopener noreferrer">jpeg-Link </a>
            </li>
            <li>
              <a :href="pngImageLink" target="_blank" rel="noopener noreferrer">png-Link </a>
            </li>
          </ul>
        </div>
      </fieldset>
      <fieldset>
        <legend>Export Fractal params as JSON</legend>
        <div>
          <div class="json-code">{{ JSON.stringify(props.fractParams, null, 2) }}</div>
          <button v-if="isSecureContext" type="button" @click="copyParamsAsJson()">copy</button>
        </div>
      </fieldset>
    </div>
  </Dialog>
</template>

<style lang="scss" scoped>
.settings-dialog {
  color: white;
  a {
    color: white;
  }
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  overflow: scroll;
  height: 100%;
}

.json-code {
  background-color: rgba(255, 255, 255, 0.1);
  padding: 0.5rem;
  font-family: monospace;
  white-space: pre;
}

.label-field {
  display: inline-flex;
  flex-direction: column;

  label {
    color: white;
    text-shadow: 1px 1px 2px black;
    font-size: 0.75rem;
  }
}
</style>
