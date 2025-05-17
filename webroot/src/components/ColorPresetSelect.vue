<script setup lang="ts">
import { apiroot, queryStr } from '@/lib/url_helper'
import { colorPresetByIdent, useColorPresets, type ColorPreset } from '@/lib/use-presets'
import { computed, reactive, ref, watch } from 'vue'
import EnhancedSelect from './EnhancedSelect.vue'
const colorPresets = useColorPresets()
const colorPreset = defineModel<string>('colorPreset')
const paletteRepeat = defineModel<number>('paletteRepeat', { required: true })
const paletteLength = defineModel<number>('paletteLength', { required: true })
const paletteReverse = defineModel<boolean>('paletteReverse', { required: true })
const paletteHardStops = defineModel<boolean>('paletteHardStops', { required: true })

const state = reactive({
  filter: '',
})

const actPresetObj = computed(() => {
  return colorPresetByIdent(colorPreset.value || '')
})

watch(colorPresets.presets, () => {
  if (!colorPreset.value) {
    colorPreset.value = colorPresets.presets.value[0].name
  }
})

const items = computed(() => {
  return colorPresets.presets.value.filter((p) =>
    p.name?.toLowerCase().includes(state.filter.toLowerCase()),
  )
})

function smallPrevLink(preset: ColorPreset) {
  const params = {
    width: 400,
    height: 8,
    maxIterations: 256,
    paletteRepeat: 1,
    paletteLength: -1,
    dir: 'horizontal',
    colorPreset: preset.ident,
  }
  return `${apiroot()}/paletteViewer?${queryStr(params)}`
}
function onPaletteModeChanged(e: Event) {
  if (e.target instanceof HTMLInputElement) {
    if (e.target.value === 'fixed') {
      paletteLength.value = 256
    } else {
      paletteLength.value = -1
    }
  }
}
</script>

<template>
  <EnhancedSelect
    :items="items"
    v-model:value="colorPreset"
    v-model:filter="state.filter"
    value-property="ident"
    display-property="name"
    filter-placeholder="e.g. 'Fire!'"
  >
    <template #display>
      <span class="value">{{ actPresetObj?.name }}</span>
    </template>
    <template #item="{ item, selected }">
      <div class="selector-item" :class="{ selected: selected }">
        <div class="item-name">
          {{ item.name }}
        </div>
        <img class="preview-img" :src="smallPrevLink(item)" loading="lazy" alt="Preview Image" />
      </div>
    </template>
    <template #settings>
      <div class="settings">
        <div class="flex gap-1">
          <div class="mode-radio">
            <div>Palette mode:</div>
            <label class="d-block">
              <input
                type="radio"
                name="paletteMode"
                value="fixed"
                :checked="paletteLength! >= 1"
                @change="onPaletteModeChanged"
              />
              fixed
            </label>
            <label class="d-block">
              <input
                type="radio"
                name="paletteMode"
                value="dynamic"
                :checked="paletteLength! < 1"
                @change="onPaletteModeChanged"
              />
              iteration-based
            </label>
          </div>
          <div class="mode-setting flex-basis-0">
            <div v-if="paletteLength < 1" class="label-field">
              <label
                >Palette Repeat:
                <input
                  type="number"
                  :value="paletteRepeat"
                  @change.lazy="
                    (e: Event) => (paletteRepeat = parseInt((e.target as HTMLInputElement).value))
                  "
                />
              </label>
            </div>

            <div v-if="paletteLength > 0" class="label-field">
              <label
                >Fixed Palette Length:
                <input
                  type="number"
                  :value="paletteLength"
                  @change.lazy="
                    (e: Event) => (paletteLength = parseInt((e.target as HTMLInputElement).value))
                  "
                />
              </label>
            </div>
          </div>
        </div>

        <div class="flex-col">
          <div class="label-field">
            <label>
              <input
                id="paletteReverse"
                type="checkbox"
                :checked="paletteReverse"
                @change.lazy="
                  (e: Event) => (paletteReverse = Boolean((e.target as HTMLInputElement).checked))
                "
              />
              Reverse Palette
            </label>
          </div>
          <div class="label-field">
            <label>
              <input
                id="paletteHardStops"
                type="checkbox"
                :checked="paletteHardStops"
                @change.lazy="
                  (e: Event) => (paletteHardStops = Boolean((e.target as HTMLInputElement).checked))
                "
              />
              Hard Color Stops
            </label>
          </div>
        </div>
      </div>
    </template>
  </EnhancedSelect>
</template>

<style lang="css" scoped>
.selector-item {
  padding: 0.3rem 0.2rem;
  &.selected {
    background-color: rgba(255, 255, 255, 0.7);
  }
  &:hover {
    background-color: rgba(255, 255, 255, 0.3);
    cursor: pointer;
  }
  .item-name {
    font-size: 0.825rem;
  }
  .preview-img {
    border: 1px solid black;
    border-radius: 3px;
    max-width: 100%;
  }
}
.settings {
  font-size: 0.75rem;
}
</style>
