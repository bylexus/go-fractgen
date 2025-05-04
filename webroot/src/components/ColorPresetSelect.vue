<script setup lang="ts">
import { apiroot, queryStr } from '@/lib/url_helper'
import {
  colorPresetByIdent,
  useColorPresets,
  type ColorPreset,
  type FractalParams,
} from '@/lib/use-presets'
import { computed, reactive, ref, watch } from 'vue'
import EnhancedSelect from './EnhancedSelect.vue'
const colorPresets = useColorPresets()
const colorPreset = defineModel<string>()

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
        <div>
          {{ item.name }}
        </div>
        <img class="preview-img" :src="smallPrevLink(item)" loading="lazy" alt="Preview Image" />
      </div>
    </template>
  </EnhancedSelect>
</template>

<style lang="css" scoped>
.selector-item {
  flex-grow: 1;
  padding: 0.3rem 0.2rem;
  &.selected {
    background-color: rgba(255, 255, 255, 0.7);
  }
  &:hover {
    background-color: rgba(255, 255, 255, 0.3);
    cursor: pointer;
  }
  .preview-img {
    border: 1px solid black;
    border-radius: 3px;
    max-width: 100%;
  }
}
</style>
