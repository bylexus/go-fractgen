<script setup lang="ts">
import { useColorPresets } from '@/lib/use-presets'
import { watch } from 'vue'
const colorPresets = useColorPresets()
const colorPreset = defineModel<string>()
watch(colorPresets.presets, () => {
  if (!colorPreset.value) {
    colorPreset.value = colorPresets.presets.value[0].name
  }
})
</script>

<template>
  <select v-model="colorPreset">
    <option v-for="preset in colorPresets.presets.value" :value="preset.ident" :key="preset.ident">
      {{ preset.name }}
    </option>
  </select>
</template>
