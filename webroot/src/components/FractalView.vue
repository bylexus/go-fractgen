<script lang="ts" setup>
import { computed, type Ref, ref, watch } from 'vue'
import ColorPresetSelect from './ColorPresetSelect.vue'
import { useSessionStorageVariable } from '@/lib/use-session-storage'
import FractalPresetsSelect from './FractalPresetsSelect.vue'
import { fractalPresetByName, useFractalPresets, type FractalParams } from '@/lib/use-presets'
import SettingsDialog from './SettingsDialog.vue'
import AboutDialog from './AboutDialog.vue'
import FractalMap from './FractalMap.vue'

const fractalPreset = ref('')
const settingsOverlay = ref<HTMLDivElement>()
const fractalMap = ref<InstanceType<typeof FractalMap>>()
const hudVisible = ref(true)
const showSettingsDlg = ref(false)
const showAboutDlg = ref(false)
const mapClickMode: Ref<'click' | 'center'> = ref('click')
const fractalPresets = useFractalPresets()
const fractalPresetFilter = ref('')

const fractalParams: Ref<FractalParams> = useSessionStorageVariable(
  'fractalParams',
  Object.assign(
    {
      width: 0, // will be updated on mount, based on the container size
      height: 0,
    },
    // we use the first preset as default:
    { ...fractalPresets.presets.value[0] },
  ),
)
// Initial values:
fractalPreset.value = fractalParams.value.name || ''

const filteredPresets = computed(() => {
  return fractalPresets.presets.value.filter((p) =>
    p.name?.toLowerCase().includes(fractalPresetFilter.value.toLowerCase()),
  )
})

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
}

function centerOnClick() {
  if (mapClickMode.value === 'click') {
    mapClickMode.value = 'center'
  } else {
    mapClickMode.value = 'click'
  }
}
</script>

<template>
  <div class="display-container">
    <FractalMap
      ref="fractalMap"
      v-model:fractalParams="fractalParams"
      :color-preset="fractalParams.colorPreset"
      :show-hud="hudVisible"
      :click-mode="mapClickMode"
      @map-single-click="hudVisible = !hudVisible"
      @map-centered="mapClickMode = 'click'"
    ></FractalMap>

    <div ref="settingsOverlay" :class="{ 'settings-overlay': true, hidden: !hudVisible }">
      <div class="glass-container"></div>
      <div class="flex flex-wrap align-end gap-1 padding-1">
        <div class="label-field">
          <label>Preset:</label>
          <FractalPresetsSelect v-model="fractalPreset"></FractalPresetsSelect>
        </div>
        <div class="label-field">
          <label>Color Palette:</label>
          <ColorPresetSelect
            :model-value="fractalParams.colorPreset"
            @update:model-value="changeFractalParams({ colorPreset: $event })"
          ></ColorPresetSelect>
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
          <label for="paletteRepeat">Palette Repeat: </label>
          <input
            type="checkbox"
            :checked="fractalParams.colorPaletteLength > 0"
            @change="onFixedPaletteCBChange"
          />
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
          <label for="paletteLength">Fixed Palette Length:</label>
          <input
            type="checkbox"
            :checked="fractalParams.colorPaletteLength > 0"
            @change="onFixedPaletteCBChange"
          />
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
        <div class="label-field">
          <label for="paletteReverse">Reverse Palette:</label>
          <input
            type="checkbox"
            :checked="fractalParams.colorPaletteReverse"
            @change.lazy="
              (e: Event) =>
                changeFractalParams({
                  colorPaletteReverse: Boolean((e.target as HTMLInputElement).checked),
                })
            "
          />
        </div>
        <button type="button" @click="showSettingsDlg = true" title="Settings">⚙️</button>
        <button
          type="button"
          :style="{ backgroundColor: mapClickMode === 'click' ? '' : '#aaa' }"
          @click="centerOnClick"
          title="Center on click"
        >
          🎯️
        </button>
        <button type="button" @click="showAboutDlg = true" title="About">ℹ️</button>
      </div>
    </div>
    <SettingsDialog
      v-model="showSettingsDlg"
      :fract-params="fractalParams"
      @update:fract-params="changeFractalParams($event)"
    />
    <AboutDialog v-model="showAboutDlg"></AboutDialog>
  </div>
</template>

<style scoped lang="scss">
.display-container {
  width: 100%;
  height: 100%;
}

.settings-overlay {
  /** Do not place transition, filter here, as it prevents childs from position:fixed! */
  position: absolute;
  z-index: 1;
  bottom: 0;
  left: 0;
  width: 100%;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
  transition:
    opacity 0.3s ease-in-out,
    transform 0.3s ease-in-out;
  &.hidden {
    opacity: 0;
    transform: translateY(100%);
  }

  .glass-container {
    position: absolute;
    z-index: -1;
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.3);
    backdrop-filter: blur(10px);
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
