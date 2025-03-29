<script setup lang="ts">
import { useFractalPresets } from '@/lib/use-presets'
import { watch } from 'vue'
const fractalPresets = useFractalPresets()
const fractalPreset = defineModel<string>()

watch(fractalPresets.presets, () => {
  if (!fractalPreset.value) {
    fractalPreset.value = fractalPresets.presets.value[0].name
  }
})
</script>

<template>
  <select v-model="fractalPreset">
    <option v-for="preset in fractalPresets.presets.value" :value="preset.name" :key="preset.name">
      {{ preset.name }}
    </option>
  </select>
</template>
