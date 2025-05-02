<script setup lang="ts">
import { apiroot, queryStr } from '@/lib/url_helper'
import { computed } from 'vue'

const props = defineProps<{
  colorPreset: string
  paletteRepeat: number
  paletteLength: number
  maxIterations: number
}>()

const imgUrl = computed(() => {
  return `${apiroot()}/paletteViewer?${queryStr({
    width: 400,
    height: 36,
    dir: 'horizontal',
    colorPreset: props.colorPreset,
    paletteRepeat: props.paletteRepeat,
    paletteLength: props.paletteLength,
    maxIterations: props.maxIterations,
  })}`
})
</script>

<template>
  <div class="flex-col gap-1">
    <a :href="imgUrl" target="_blank" rel="noopener noreferrer">
      <img :src="imgUrl" width="400" height="36" />
    </a>
    <div class="flex justify-between labels">
      <span>0</span>
      <span>{{ paletteLength <= 0 ? maxIterations : paletteLength }}</span>
    </div>
  </div>
</template>

<style lang="css" scoped>
img {
  max-width: 100%;
}
.labels {
  font-size: 0.8rem;
}
</style>
