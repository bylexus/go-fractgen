<script lang="ts" setup>
import { type ElementInfo, useElementResize } from '@/lib/element-info'
import { queryStr } from '@/lib/url_helper'
import { onMounted, Ref, ref, watch, watchEffect } from 'vue'

const imageUrl = ref(null)
const loading = ref(true)
const imgContainer: Ref<HTMLImageElement | null> = ref(null)
const windowSizes = useElementResize(window.document.body, 1000)

const fractalParams = {
  width: 0,
  height: 0,
  iterFunc: 'Mandelbrot',
  maxIterations: 40,
  centerCY: 0,
  centerCX: -0.7,
  diameterCX: 4,
}

onMounted(() => {
  loading.value = true
  fractalParams.width = imgContainer.value?.clientWidth || 0
  fractalParams.height = imgContainer.value?.clientHeight || 0
  calcImage(fractalParams)
})

watch(windowSizes.sizes, ({ width, height }) => {
  fractalParams.width = width
  fractalParams.height = height
  calcImage(fractalParams)
})

function calcImage(fractalParams: any) {
  loading.value = true
  imageUrl.value = `http://localhost:8000/fractal-image.png?${queryStr(fractalParams)}`
}
</script>

<template>
  <div class="img-container" ref="imgContainer">
    <img
      :src="imageUrl"
      @load="loading = false"
      alt="Fractal Image"
    />
    <div v-if="loading" class="loading-overlay">Calculating...</div>
  </div>
</template>

<style scoped>
.img-container {
  width: 100%;
  height: 100%;
}

.loading-overlay {
  position: absolute;
  z-index: 2;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  color: white;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.5);
}
</style>
