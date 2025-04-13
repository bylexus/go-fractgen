<script lang="ts" setup>
import { useElementResize } from '@/lib/element-info'
import { apiroot, queryStr } from '@/lib/url_helper'
import { onMounted, reactive, type Ref, ref, watch, watchEffect } from 'vue'
import ColorPresetsSelect from './ColorPresetsSelect.vue'
import { useLocalStorageVariable } from '@/lib/use-local-storage'
import FractalPresetsSelect from './FractalPresetsSelect.vue'
import { fractalPresetByName, useFractalPresets, type FractalPreset } from '@/lib/use-presets'
import OlMap from 'ol/Map'
import View from 'ol/View'
import TileLayer from 'ol/layer/Tile'
import WMTS from 'ol/source/WMTS'
import WMTSTileGrid from 'ol/tilegrid/WMTS'
import Projection from 'ol/proj/Projection'
import type { ImageTile, Tile } from 'ol'

const loading = ref(false)
const colorPreset = ref('')
const fractalPreset = ref('')
const map = ref<HTMLDivElement | null>(null)
const tileWidth = 256
const maxZoom = 64

let disableIterRecalcOnZoom = false

let olMap: OlMap | null = null
let fractalOlLayer: TileLayer | null = null

const fractalParams: Ref<FractalPreset & { width: number; height: number }> =
  useLocalStorageVariable(
    'fractalParams',
    Object.assign(
      {
        width: 0, // will be updated on mount, based on the container size
        height: 0,
      },
      // we use the first preset as default:
      { ...useFractalPresets().presets.value[0] },
    ),
  )

onMounted(() => {
  // minx, miny, maxx, maxy
  const mandelbrotExtent = [-1.7, -1, 0.3, 1] // Complex plane bounds
  // resolution means: unit per pixel (in our case: fractal pixel per screen pixel)
  // as a tile is 256 x 256 pixels, the
  const resolutions = Array.from({ length: maxZoom }, (_, z) => 1 / (tileWidth * Math.pow(2, z)))

  const tileGrid = new WMTSTileGrid({
    origin: [mandelbrotExtent[0], mandelbrotExtent[1]],
    resolutions: resolutions,
    tileSize: tileWidth,
    matrixIds: resolutions.map((_, i) => i.toString()),
  })
  // console.log(tileGrid)
  const fractalSource = new WMTS({
    layer: 'fractal',
    style: 'fractal',
    matrixSet: 'default',
    url: `${apiroot()}/wmts`,
    tileGrid: tileGrid,
    format: 'image/png',
    tileLoadFunction: (tile: Tile, src) => {
      const zoomLevel = tile.getTileCoord()[0]
      // console.log(tile)
      // console.log(view.getZoom())
      // console.log(tile.getTileCoord())
      // const maxIterations = Math.ceil(50 * Math.pow(1.3, zoomLevel))
      const maxIterations = fractalParams.value.maxIterations
      const colorPaletteRepeat = fractalParams.value.colorPaletteRepeat

      const imgTile = tile as ImageTile
      const urlParams = fractParamsAsQueryParams({
        TileMatrix: zoomLevel,
        maxIterations: maxIterations,
        colorPaletteRepeat: colorPaletteRepeat,
        tileWidthPixels: tileWidth,
        tileWidthFractal: (view.getResolution() || 0) * tileWidth,
      })
      ;(imgTile.getImage() as HTMLImageElement).src = `${src}&${urlParams}`
    },
  })
  // console.log(tileGrid)
  fractalOlLayer = new TileLayer({
    source: fractalSource,
  })
  const view = new View({
    center: [-0.7, 0],
    zoom: 1,
    constrainResolution: true,
    maxZoom: maxZoom,
    projection: new Projection({
      code: 'MANDELBROT',
      extent: mandelbrotExtent,
    }),
  })
  view.on('change:resolution', () => {
    // console.log(view.getZoom())
  })
  olMap = new OlMap({
    target: map.value!,
    layers: [fractalOlLayer],
    view: view,
  })
  let oldZoom = 0
  olMap.on('movestart', () => {
    oldZoom = view.getZoom() || 0
    console.log('start zoom', view.getZoom())
  })
  olMap.on('moveend', () => {
    if (disableIterRecalcOnZoom) {
      disableIterRecalcOnZoom = false
      return
    }
    console.log('end zoom', view.getZoom())
    if (oldZoom !== view.getZoom()) {
      let zoomDiff = view.getZoom()! - oldZoom
      if (zoomDiff > 0) {
        fractalParams.value.maxIterations = Math.floor(
          fractalParams.value.maxIterations * Math.pow(1.25, zoomDiff || 1),
        )
      } else {
        zoomDiff = Math.abs(zoomDiff)
        fractalParams.value.maxIterations = Math.floor(
          fractalParams.value.maxIterations / Math.pow(1.25, zoomDiff || 1),
        )
      }
      // fractalParams.value.maxIterations = Math.floor(50 * Math.pow(1.3, (view.getZoom() || 1) - 1))
      oldZoom = view.getZoom() || 0
    }
  })

  // Initial values:
  colorPreset.value = fractalParams.value.colorPreset || ''
  fractalPreset.value = fractalParams.value.name || ''
})

watch(colorPreset, () => {
  fractalParams.value.colorPreset = colorPreset.value || ''
})

watch(fractalPreset, () => {
  const preset = fractalPresetByName(fractalPreset.value)
  if (preset) {
    disableIterRecalcOnZoom = true
    fractalParams.value.iterFunc = preset.iterFunc
    fractalParams.value.maxIterations = preset.maxIterations
    fractalParams.value.centerCX = preset.centerCX
    fractalParams.value.centerCY = preset.centerCY
    fractalParams.value.diameterCX = preset.diameterCX
    fractalParams.value.colorPreset = preset.colorPreset
    fractalParams.value.colorPaletteRepeat = preset.colorPaletteRepeat
    fractalParams.value.name = preset.name || ''
    colorPreset.value = preset.colorPreset

    if (olMap) {
      // olMap.getView().setCenter([preset.centerCX, preset.centerCY])
      // olMap.getView().setZoom(1)
      // console.log(olMap.getView().getZoomForResolution(preset.diameterCX / tileWidth))
      // olMap.getView().setResolution(preset.diameterCX / tileWidth)
      const halfWidth = preset.diameterCX / 2
      const halfHeight = preset.diameterCX / 2
      olMap
        .getView()
        .fit([
          preset.centerCX - halfWidth,
          preset.centerCY - halfHeight,
          preset.centerCX + halfWidth,
          preset.centerCY + halfHeight,
        ])
      // fractalParams.value.maxIterations = preset.maxIterations
    }
  }
})

watch(
  fractalParams,
  () => {
    calcImage(fractalParams.value)
  },
  { deep: true },
)

function fractParamsAsQueryParams(inputObj: { [key: string]: any }) {
  return queryStr({ ...inputObj, colorPreset: colorPreset.value })
}

function recalcIterations(diameterCX: number) {
  // approximation of the number of iterations, based on the following formula,
  // maxIterations := int(40 * math.Pow(1.3, float64(zoomLevel)))

  // which seems to work well:
  // maxIter = 50 * (log10(scale))^1.25
  // where scale is pixelWidth/complexPlaneWidth e.g. 1280/5
  return Math.ceil(50 * Math.pow(Math.log10(fractalParams.value.width / diameterCX), 1.25))
}

function calcImage(fractalParams: any) {
  console.log('recalc:', fractalParams)
  fractalOlLayer?.getSource()?.changed()
}
</script>

<template>
  <div class="display-container">
    <div ref="map" class="img-map"></div>
    <div class="settings-overlay">
      <div class="label-field">
        <label>Preset:</label>
        <FractalPresetsSelect v-model="fractalPreset"></FractalPresetsSelect>
      </div>
      <div class="label-field">
        <label>Color Palette:</label>
        <ColorPresetsSelect v-model="colorPreset"></ColorPresetsSelect>
      </div>
      <div class="label-field">
        <label for="iterations">Max. Iterations</label>
        <input type="number" v-model.lazy="fractalParams.maxIterations" id="iterations" />
      </div>
      <div class="label-field">
        <label for="paletteRepeat">Palette Repeat:</label>
        <input type="number" v-model.lazy="fractalParams.colorPaletteRepeat" id="paletteRepeat" />
      </div>
    </div>
    <div v-if="loading" class="loading-overlay">Calculating...</div>
  </div>
</template>

<style scoped>
@import 'ol/ol.css';

.display-container {
  width: 100%;
  height: 100%;
}

.img-map {
  width: 100%;
  height: 100%;
  background-color: blue;
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

.settings-overlay {
  position: absolute;
  padding: 0.2rem;
  z-index: 1;
  bottom: 0;
  left: 0;
  width: 100%;
  background-color: rgba(255, 255, 255, 0.3);
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
