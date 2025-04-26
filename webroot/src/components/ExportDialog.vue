<script setup lang="ts">
import type { FractalParams } from '@/lib/use-presets'
import Dialog from './Dialog.vue'
import { computed } from 'vue'
import { apiroot, queryStr } from '@/lib/url_helper'
const props = defineProps<{
  fractParams: FractalParams
}>()

function fractParamsAsQueryParams(inputObj: { [key: string]: any }) {
  return queryStr({ ...inputObj })
}

const imageLink = computed(() => {
  return `${apiroot()}/fractal-image.jpg?${fractParamsAsQueryParams(props.fractParams)}`
})
</script>

<template>
  <Dialog>
    <div class="export-dialog">
      <fieldset>
        <legend>Export Image</legend>
        Image URL:
        <a :href="imageLink" target="_blank" rel="noopener noreferrer">{{ imageLink }} </a>
      </fieldset>
      <fieldset>
        <legend>Export Fractal params as JSON</legend>
      </fieldset>
    </div>
  </Dialog>
</template>

<style lang="scss" scoped>
.export-dialog {
  color: white;
  a {
    color: white;
  }
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
</style>
