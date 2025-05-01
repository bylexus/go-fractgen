<script lang="ts" setup>
import { type Ref, ref, watch } from 'vue'
import ColorPresetsSelect from './ColorPresetsSelect.vue'
import { useSessionStorageVariable } from '@/lib/use-session-storage'
import FractalPresetsSelect from './FractalPresetsSelect.vue'
import { fractalPresetByName, useFractalPresets, type FractalParams } from '@/lib/use-presets'
import SettingsDialog from './SettingsDialog.vue'
import FractalMap from './FractalMap.vue'

const fractalPreset = ref('')
const settingsOverlay = ref<HTMLDivElement>()
const fractalMap = ref<InstanceType<typeof FractalMap>>()
const hudVisible = ref(true)
const showSettingsDlg = ref(false)

const fractalParams: Ref<FractalParams> = useSessionStorageVariable(
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
fractalPreset.value = fractalParams.value.name || ''

watch(fractalPreset, () => {
  const preset = fractalPresetByName(fractalPreset.value)
  if (preset) {
    changeFractalParams(preset)
  }
})

watch(fractalParams, (newVal) => {
  console.log('changed params: ', newVal)
})

function changeFractalParams(params: Partial<FractalParams>) {
  fractalParams.value = { ...fractalParams.value, ...params }
}

function recalcIterations() {
  // approximation of the number of iterations, based on the following formula,
  // maxIterations := int(40 * math.Pow(1.3, float64(zoomLevel)))

  // which seems to work well:
  // maxIter = 50 * (log10(scale))^1.25
  // where scale is pixelWidth/complexPlaneWidth e.g. 1280/5
  console.log('width: ', fractalParams.value.width, 'diameterCX: ', fractalParams.value.diameterCX)
  const optimalIterationsForScale = Math.ceil(
    60 *
      Math.pow(Math.log10((fractalParams.value.width || 0) / fractalParams.value.diameterCX), 1.25),
  )
  changeFractalParams({
    maxIterations: optimalIterationsForScale,
  })
}

function onFixedPaletteCBChange(e: Event) {
  const checked = (e.target as HTMLInputElement).checked
  if (checked) {
    if (fractalParams.value.colorPaletteLength <= 0) {
      changeFractalParams({ colorPaletteLength: 256 })
    }
  } else {
    changeFractalParams({
      colorPaletteLength: -1,
      colorPaletteRepeat: fractalParams.value.colorPaletteRepeat || 1,
    })
  }
  console.log(checked)
}
</script>

<template>
  <div class="display-container">
    <FractalMap
      ref="fractalMap"
      v-model:fractalParams="fractalParams"
      :color-preset="fractalParams.colorPreset"
      :show-hud="hudVisible"
      @map-single-click="hudVisible = !hudVisible"
    ></FractalMap>

    <div ref="settingsOverlay" :class="{ 'settings-overlay': true, hidden: !hudVisible }">
      <div class="label-field">
        <label>Preset:</label>
        <FractalPresetsSelect v-model="fractalPreset"></FractalPresetsSelect>
      </div>
      <div class="label-field">
        <label>Color Palette:</label>
        <ColorPresetsSelect
          :model-value="fractalParams.colorPreset"
          @change.lazy="
            (e: Event) => changeFractalParams({ colorPreset: (e.target as HTMLInputElement).value })
          "
        ></ColorPresetsSelect>
      </div>
      <div class="label-field">
        <label for="iterations">Max. Iterations</label>
        <div style="display: flex; gap: 0.25rem; align-items: end">
          <input
            type="number"
            :value="fractalParams.maxIterations"
            @change.lazy="
              (e: Event) =>
                changeFractalParams({
                  maxIterations: parseInt((e.target as HTMLInputElement).value),
                })
            "
            id="iterations"
          />
          <button type="button" @click="recalcIterations()">🧮</button>
        </div>
      </div>
      <div v-if="fractalParams.colorPaletteLength <= 0" class="label-field">
        <label for="paletteRepeat"
          >Palette Repeat:
          <input
            type="checkbox"
            :checked="fractalParams.colorPaletteLength > 0"
            @change="onFixedPaletteCBChange"
          />
        </label>
        <input
          type="number"
          :value="fractalParams.colorPaletteRepeat"
          id="paletteRepeat"
          @change.lazy="
            (e: Event) =>
              changeFractalParams({
                colorPaletteRepeat: parseInt((e.target as HTMLInputElement).value),
              })
          "
        />
      </div>
      <div v-if="fractalParams.colorPaletteLength > 0" class="label-field">
        <label for="paletteLength"
          >Fixed Palette Length:
          <input
            type="checkbox"
            :checked="fractalParams.colorPaletteLength > 0"
            @change="onFixedPaletteCBChange"
        /></label>
        <input
          type="number"
          :value="fractalParams.colorPaletteLength"
          id="paletteLength"
          @change.lazy="
            (e: Event) =>
              changeFractalParams({
                colorPaletteLength: parseInt((e.target as HTMLInputElement).value),
              })
          "
        />
      </div>
      <button type="button" @click="hudVisible = !hudVisible">hide HUD</button>
      <button type="button" @click="showSettingsDlg = true" title="Settings">⚙️</button>
    </div>
    <SettingsDialog
      v-model="showSettingsDlg"
      :fract-params="fractalParams"
      @update:fract-params="changeFractalParams($event)"
    />
  </div>
</template>

<style scoped lang="scss">
.display-container {
  width: 100%;
  height: 100%;
}

.settings-overlay {
  position: absolute;
  padding: 0.2rem;
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  gap: 0.2rem;
  z-index: 1;
  bottom: 0;
  left: 0;
  width: 100%;
  background-color: rgba(255, 255, 255, 0.3);
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
  transition:
    opacity 0.2s ease-in-out,
    bottom 0.2s ease-in-out;
  &.hidden {
    opacity: 0;
    bottom: -100%;
  }
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
