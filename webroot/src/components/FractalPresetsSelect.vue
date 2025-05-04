<script setup lang="ts">
import { apiroot, queryStr } from '@/lib/url_helper'
import { useFractalPresets, type FractalParams, type FractalPreset } from '@/lib/use-presets'
import { computed, reactive, ref, watch } from 'vue'
const fractalPresets = useFractalPresets()
const fractalPreset = defineModel<string>()
const filterInput = ref<HTMLInputElement>()

const state = reactive({
  selectorOpen: false,
  filter: '',
})

watch(fractalPresets.presets, () => {
  if (!fractalPreset.value) {
    fractalPreset.value = fractalPresets.presets.value[0].name
  }
})

watch(
  () => state.selectorOpen,
  (newVal, oldVal) => {
    if (newVal && !oldVal) {
      state.filter = ''
      filterInput.value?.focus()
    }
  },
)

const items = computed(() => {
  return fractalPresets.presets.value.filter((p) =>
    p.name?.toLowerCase().includes(state.filter.toLowerCase()),
  )
})

function selectPreset(name: string | null) {
  if (!name) {
    return
  }
  fractalPreset.value = name
  state.selectorOpen = false
}

function smallPrevLink(fractalParams: FractalParams) {
  const params = { ...fractalParams }
  params.width = 36
  params.height = 36
  return `${apiroot()}/fractal-image/jpg?${queryStr(params)}`
}
</script>

<template>
  <div class="container" @click="state.selectorOpen = true">
    <span class="value">{{ fractalPreset }}</span>
    <!--
    <select v-model="fractalPreset" @click="state.selectorOpen = !state.selectorOpen">
      <option
        v-for="preset in fractalPresets.presets.value"
        :value="preset.name"
        :key="preset.name"
      >
        {{ preset.name }}
      </option>
    </select>
    -->
    <div :class="{ 'selector-overlay': true, open: state.selectorOpen }">
      <div class="search-field">
        <input type="text" placeholder="e.g. Mandelbrot" v-model="state.filter" ref="filterInput" />
        <button type="button" @click.stop="state.selectorOpen = false">🅧</button>
      </div>
      <div class="items-list">
        <div
          v-for="preset in items"
          :key="preset.name"
          class="selector-item"
          :class="{ selected: preset.name === fractalPreset }"
          @click.stop="selectPreset(preset.name || null)"
        >
          <div>
            {{ preset.name }}
          </div>
          <img class="preview-img" :src="smallPrevLink(preset)" alt="Preview Image" />
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
.container {
  position: relative;
  display: inline-block;
  border: 1px solid black;
  border-radius: 3px;
  padding: 0.2rem 0.3rem;
  background-color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  .value {
    color: #333;
  }

  .selector-overlay {
    color: white;
    position: absolute;
    width: 400px;
    max-width: 100vw;
    height: 400px;
    max-height: 400px;
    opacity: 0;
    left: 0;
    bottom: 0;
    z-index: 2;
    background-color: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(10px);
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
    border: 1px solid #aaa;
    border-radius: 3px;
    overflow: none;

    display: flex;
    flex-direction: column;
    align-items: start;
    gap: 0.2rem;
    padding: 0.5rem;

    transition: opacity 0.2s ease-in-out;
    transform: translateY(100vh);

    &.open {
      opacity: 1;
      transform: translateY(0);
    }

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

    .items-list {
      flex-grow: 1;
      width: 100%;
      overflow: auto;
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
        }
      }
    }
  }
}
</style>
