<script setup lang="ts">
import OlMap from 'ol/Map'
import View from 'ol/View'
import TileLayer from 'ol/layer/Tile'
import WMTS from 'ol/source/WMTS'
import WMTSTileGrid from 'ol/tilegrid/WMTS'
import Projection from 'ol/proj/Projection'
import DoubleClickZoom from 'ol/interaction/DoubleClickZoom'
import MouseWheelZoom from 'ol/interaction/MouseWheelZoom'
import DragPan from 'ol/interaction/DragPan'
import DragZoom from 'ol/interaction/DragZoom'
import PinchZoom from 'ol/interaction/PinchZoom'
import { defaults as defaultControls } from 'ol/control/defaults'
import type { ImageTile, Tile } from 'ol'
import { onMounted, ref, watch, watchEffect, type ModelRef, type Ref } from 'vue'
import { apiroot, queryStr } from '@/lib/url_helper'
import type { FractalParams } from '@/lib/use-presets'
import { useElementResize, type ElementInfo } from '@/lib/element-info'

const fractalParams: ModelRef<FractalParams> = defineModel('fractalParams', { required: true })
const props = defineProps<{
  colorPreset: string
  showHud: boolean
}>()

const emit = defineEmits(['mapSingleClick'])

const map = ref<HTMLDivElement>()
const mapControls = ref<HTMLDivElement>()

const tileWidth = 256
const maxZoom = 46
const mapRedrawProps = [
  'iterFunc',
  'maxIterations',
  'colorPreset',
  'colorPaletteRepeat',
  'colorPaletteLength',
  'juliaKr',
  'juliaKi',
]

let olMap: OlMap
let fractalOlLayer: TileLayer | null = null
let zoomInProgress = false

onMounted(() => {
  const mapSizes = useElementResize(map.value!)
  watch(
    mapSizes.sizes,
    () => {
      changeFractalParams({
        width: mapSizes.width.value,
        height: mapSizes.height.value,
      })
    },
    {
      immediate: true,
    },
  )

  // minx, miny, maxx, maxy
  const mandelbrotExtent = [-1.7, -1, 0.3, 1] // Complex plane bounds
  // resolution means: unit per pixel (in our case: fractal pixel per screen pixel)
  // as a tile is 256 x 256 pixels, the
  const resolutions = Array.from({ length: maxZoom }, (_, z) => calcResolutionForZoom(z))

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
    format: 'image/jpeg',
    tileLoadFunction: (tile: Tile, src) => {
      const zoomLevel = tile.getTileCoord()[0]
      // const maxIterations = Math.ceil(50 * Math.pow(1.3, zoomLevel))
      const maxIterations = fractalParams.value.maxIterations
      const colorPaletteRepeat = fractalParams.value.colorPaletteRepeat
      const colorPaletteLength = fractalParams.value.colorPaletteLength

      const res = calcResolutionForZoom(zoomLevel)

      const imgTile = tile as ImageTile
      const urlParams = fractParamsAsQueryParams({
        TileMatrix: zoomLevel,
        maxIterations: maxIterations,
        colorPreset: props.colorPreset,
        colorPaletteRepeat: colorPaletteRepeat,
        colorPaletteLength: colorPaletteLength,
        iterFunc: fractalParams.value.iterFunc,
        tileWidthPixels: tileWidth,
        tileWidthFractal: res * tileWidth,
        juliaKr: fractalParams.value.juliaKr,
        juliaKi: fractalParams.value.juliaKi,
      })
      ;(imgTile.getImage() as HTMLImageElement).src = `${src}&${urlParams}`
    },
  })

  fractalOlLayer = new TileLayer({
    source: fractalSource,
  })

  const view = new View({
    center: [-0.7, 0],
    zoom: 1,
    constrainResolution: true,
    maxZoom: maxZoom,
    enableRotation: false,
    constrainOnlyCenter: true,
    projection: new Projection({
      code: 'MANDELBROT',
      extent: mandelbrotExtent,
    }),
  })

  olMap = new OlMap({
    target: map.value!,
    layers: [fractalOlLayer],
    view: view,
    interactions: [
      new DoubleClickZoom({ duration: 250 }),
      new MouseWheelZoom({ duration: 250 }),
      new DragPan(),
      new DragZoom({ duration: 250 }),
      new PinchZoom({ duration: 250 }),
    ],
    controls: defaultControls({
      rotate: false,
      zoom: true,
      zoomOptions: {
        target: mapControls.value,
        duration: 250,
      },
      attribution: true,
    }),
  })

  olMap.on('singleclick', () => emit('mapSingleClick'))

  let oldZoom = 0
  olMap.on('movestart', () => {
    oldZoom = view.getZoom() || 0
    zoomInProgress = true
  })
  olMap.on('moveend', () => {
    zoomInProgress = false
    console.log('end zoom', view.getZoom())
    // if (oldZoom !== view.getZoom()) {
    //   refreshTiles = true
    //   let zoomDiff = view.getZoom()! - oldZoom
    //   if (zoomDiff > 0) {
    //     fractalParams.value.maxIterations = Math.floor(
    //       fractalParams.value.maxIterations * Math.pow(1.25, zoomDiff || 1),
    //     )
    //   } else {
    //     zoomDiff = Math.abs(zoomDiff)
    //     fractalParams.value.maxIterations = Math.ceil(
    //       fractalParams.value.maxIterations / Math.pow(1.25, zoomDiff || 1),
    //     )
    //   }
    //   // fractalParams.value.maxIterations = Math.floor(50 * Math.pow(1.3, (view.getZoom() || 1) - 1))
    //   oldZoom = view.getZoom() || 0
    // }

    // Update fractal parameters to reflect actual map settings:
    const actExtent = olMap.getView().calculateExtent()
    const viewCenter = olMap.getView().getCenter()!
    changeFractalParams({
      centerCX: viewCenter[0],
      centerCY: viewCenter[1],
      diameterCX: actExtent[2] - actExtent[0],
    })
  })

  calcMap(fractalParams.value)
})

watch(fractalParams, (newVal, oldVal) => {
  // check for relevant changes that need a map reload: This is not the case when the user just
  // pans / zooms
  let refresh = false
  for (const prop of mapRedrawProps) {
    if (newVal[prop] !== oldVal[prop]) {
      console.log('prop changed: ', prop)
      refresh = true
      break
    }
  }

  calcMap(fractalParams.value, refresh)
})

function changeFractalParams(params: Partial<FractalParams>) {
  fractalParams.value = { ...fractalParams.value, ...params }
}

function calcMap(fractalParams: FractalParams, refresh = false) {
  // olMap.getView().setCenter([preset.centerCX, preset.centerCY])
  // olMap.getView().setZoom(1)
  console.log(
    'actual zoom for resolution: ',
    olMap.getView().getZoomForResolution(fractalParams.diameterCX / tileWidth),
  )
  // olMap.getView().setResolution(preset.diameterCX / tileWidth)
  const halfWidth = fractalParams.diameterCX / 2
  const halfHeight = fractalParams.diameterCX / 2
  olMap.getView().fit([
    // because the height is not part of the calculation (we only have the diameter in the X direction),
    // we make sure it is something very small, to fit into the map window anyway. This means that
    // only the width is taken into account for fitting:
    fractalParams.centerCX - halfWidth,
    fractalParams.centerCY - halfHeight * 0.01,
    fractalParams.centerCX + halfWidth,
    fractalParams.centerCY + halfHeight * 0.01,
  ])
  const imageParams = getActualFractParams()
  // TODO: this image link should be placed on a button or in a small menu to
  // choose the destination / output size, then generate an image from it:
  let imgLink = `${apiroot()}/fractal-image.jpg?${fractParamsAsQueryParams(imageParams)}`
  console.log(imgLink)
  // console.log(imgLink)
  // fractalOlLayer?.getSource()?.changed()
  if (refresh) {
    refreshTiles()
  }
}

function refreshTiles() {
  fractalOlLayer?.getSource()?.changed()
}

function getActualFractParams(): FractalParams {
  const actExtent = olMap.getView().calculateExtent()
  const viewCenter = olMap.getView().getCenter()!
  return {
    width: olMap.getSize()![0],
    height: olMap.getSize()![1],
    iterFunc: fractalParams.value.iterFunc,
    maxIterations: fractalParams.value.maxIterations,
    centerCX: viewCenter[0],
    centerCY: viewCenter[1],
    diameterCX: actExtent[2] - actExtent[0],
    colorPreset: props.colorPreset,
    colorPaletteRepeat: fractalParams.value.colorPaletteRepeat,
    juliaKr: fractalParams.value.juliaKr,
    juliaKi: fractalParams.value.juliaKi,
  }
}

function fractParamsAsQueryParams(inputObj: { [key: string]: any }) {
  return queryStr({ ...inputObj })
}

function calcResolutionForZoom(zoomLevel: number) {
  return 4 / (tileWidth * Math.pow(2, zoomLevel))
}
</script>

<template>
  <div class="map-container">
    <div ref="map" class="map"></div>
    <div ref="mapControls" :class="{ 'map-controls': true, hidden: !showHud }"></div>
  </div>
</template>

<style lang="scss" scoped>
@import 'ol/ol.css';

.map-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
  .map {
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .map-controls {
    display: inline-block;
    position: absolute;
    z-index: 1;
    top: 0;
    right: 0;
    margin: 0.5rem;
    transition:
      opacity 0.2s ease-in-out,
      top 0.2s ease-in-out;
    &.hidden {
      opacity: 0;
      top: -100%;
    }

    :deep(.ol-zoom) {
      display: flex;
      gap: 0.2rem;
      button {
        border-radius: 50%;
        font-size: 1.3rem;
        width: 2rem;
        height: 2rem;
        border: 1px solid #aaa;
        box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.4);
        opacity: 0.75;
        &:hover {
          opacity: 1;
          box-shadow: 1px 1px 5px rgba(0, 0, 0, 0.5);
        }
      }
    }
  }
}
</style>

<style lang="css">
.ol-dragzoom {
  border: 2px solid white !important;
}
</style>
