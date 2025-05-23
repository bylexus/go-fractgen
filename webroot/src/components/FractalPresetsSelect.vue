<script setup lang="ts">
import { apiroot, queryStr } from '@/lib/url_helper'
import { useFractalPresets, type FractalParams, type FractalPreset } from '@/lib/use-presets'
import { computed, reactive, ref, watch } from 'vue'
import EnhancedSelect from './EnhancedSelect.vue'
const fractalPresets = useFractalPresets()
const fractalPreset = defineModel<string>()

const state = reactive({
  filter: '',
})

watch(fractalPresets.presets, () => {
  if (!fractalPreset.value) {
    fractalPreset.value = fractalPresets.presets.value[0].name
  }
})

const items = computed(() => {
  return fractalPresets.presets.value.filter((p) =>
    p.name?.toLowerCase().includes(state.filter.toLowerCase()),
  )
})

function smallPrevLink(fractalParams: FractalParams) {
  const params = { ...fractalParams }
  params.width = 36
  params.height = 36
  return `${apiroot()}/fractal-image/jpg?${queryStr(params)}`
}
</script>

<template>
  <EnhancedSelect
    :items="items"
    v-model:value="fractalPreset"
    v-model:filter="state.filter"
    value-property="name"
    display-property="name"
    filter-placeholder="e.g. 'Mandelbrot'"
  >
    <template #display>
      <span class="value">{{ fractalPreset }}</span>
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
.search-field {
  display: flex;
  justify-content: space-between;
  width: 100%;
  input {
    flex-grow: 1;
    padding: 0.3rem 0.2rem;
  }
  button {
    border: none;
    background-color: transparent;
    padding: none;
    cursor: pointer;
    color: white;
    font-size: 1.25rem;
  }
}

.selector-item {
  flex-grow: 1;
  padding: 0.3rem 0.2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  }
}
</style>
