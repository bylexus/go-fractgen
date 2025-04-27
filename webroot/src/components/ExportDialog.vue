<script setup lang="ts">
import type { FractalParams } from '@/lib/use-presets'
import Dialog from './Dialog.vue'
import { computed } from 'vue'
import { apiroot, queryStr } from '@/lib/url_helper'
const props = defineProps<{
  fractParams: FractalParams
}>()

const isSecureContext = computed(() => {
  return window.isSecureContext
})

function fractParamsAsQueryParams(inputObj: { [key: string]: any }) {
  return queryStr({ ...inputObj })
}

const jpegImageLink = computed(() => {
  return `${apiroot()}/fractal-image/jpg?${fractParamsAsQueryParams(props.fractParams)}`
})
const pngImageLink = computed(() => {
  return `${apiroot()}/fractal-image/png?${fractParamsAsQueryParams(props.fractParams)}`
})

function copyParamsAsJson() {
  window.navigator.clipboard.writeText(JSON.stringify(props.fractParams, null, 2))
}
</script>

<template>
  <Dialog>
    <div class="export-dialog">
      <fieldset>
        <legend>Export Image</legend>
        <div>
          Image URLs:
          <ul>
            <li>
              <a :href="jpegImageLink" target="_blank" rel="noopener noreferrer">jpeg-Link </a
              ><br />
            </li>
            <li>
              <a :href="pngImageLink" target="_blank" rel="noopener noreferrer">png-Link </a><br />
            </li>
          </ul>
        </div>
      </fieldset>
      <fieldset>
        <legend>Export Fractal params as JSON</legend>
        <div>
          <div class="json-code">{{ JSON.stringify(props.fractParams, null, 2) }}</div>
          <button v-if="isSecureContext" type="button" @click="copyParamsAsJson()">copy</button>
        </div>
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
  overflow: scroll;
  height: 100%;
}

.json-code {
  background-color: rgba(255, 255, 255, 0.1);
  padding: 0.5rem;
  font-family: monospace;
  white-space: pre;
}
</style>
