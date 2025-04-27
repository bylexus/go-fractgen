<script setup lang="ts">
import type { FractalParams } from '@/lib/use-presets'
import Dialog from './Dialog.vue'
import { computed, reactive } from 'vue'
import { apiroot, queryStr } from '@/lib/url_helper'
import { screenSize } from '@/lib/element-info'
const props = defineProps<{
  fractParams: FractalParams
}>()

const { physicalWidth, physicalHeight } = screenSize()

const state = reactive({
  imgWidth: physicalWidth || props.fractParams.width,
  imgHeight: physicalHeight || props.fractParams.height,
})

const isSecureContext = computed(() => {
  return window.isSecureContext
})

function fractParamsAsQueryParams(inputObj: { [key: string]: any }) {
  return queryStr({ ...inputObj })
}

const jpegImageLink = computed(() => {
  const params = { ...props.fractParams }
  params.width = state.imgWidth
  params.height = state.imgHeight
  return `${apiroot()}/fractal-image/jpg?${fractParamsAsQueryParams(params)}`
})
const pngImageLink = computed(() => {
  const params = { ...props.fractParams }
  params.width = state.imgWidth
  params.height = state.imgHeight
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
        <div class="label-field">
          <label for="imgWidth">Image Width (px)</label>
          <input type="number" v-model="state.imgWidth" id="imgWidth" />
        </div>
        <div class="label-field">
          <label for="imgHeight">Image Height (px)</label>
          <input type="number" v-model="state.imgHeight" id="imgHeight" />
        </div>
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
