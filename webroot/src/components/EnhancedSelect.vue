<script setup lang="ts">
import { GlobalEvent, onEvent } from '@/lib/event-bus'
import { computed, onMounted, reactive, ref, watch } from 'vue'

const model = defineModel<any>('value')
const filter = defineModel<any>('filter')
const container = ref<HTMLDivElement>()
const selectorOverlay = ref<HTMLDivElement>()

const props = defineProps<{
  items: any[]
  valueProperty: string
  displayProperty: string
  filterPlaceholder?: string
}>()
const filterInput = ref<HTMLInputElement>()

const state = reactive({
  selectorOpen: false,
})

watch(
  () => state.selectorOpen,
  (newVal, oldVal) => {
    // on open:
    if (newVal && !oldVal) {
      filter.value = ''
      filterInput.value?.focus()

      const containerBounds = container.value?.getBoundingClientRect()
      if (containerBounds) {
        if (containerBounds.left + 400 > window.innerWidth) {
          selectorOverlay.value!.style.left = 'unset'
          selectorOverlay.value!.style.right = '0'
        } else {
          selectorOverlay.value!.style.left = `${containerBounds.left}px`
          selectorOverlay.value!.style.right = 'unset'
        }
        if (containerBounds.bottom - 400 < 0) {
          selectorOverlay.value!.style.top = '0'
          selectorOverlay.value!.style.bottom = 'unset'
        } else {
          selectorOverlay.value!.style.bottom = `${window.innerHeight - containerBounds.bottom}px`
          selectorOverlay.value!.style.top = 'unset'
        }
      }
      console.log(container.value?.getBoundingClientRect())
    } else {
      // selectorOverlay.value!.style.left = `-200vw`
      selectorOverlay.value!.style.right = 'unset'
    }
  },
)

const valueItem = computed(() => {
  return props.items.find((item: any) => item[props.valueProperty] === model.value)
})

const valueDisplay = computed(() => {
  return valueItem.value ? valueItem.value[props.displayProperty] : ''
})

onMounted(() => {
  onEvent(GlobalEvent.hideHud, () => {
    state.selectorOpen = false
  })
})

function onSelect(item: any) {
  model.value = item[props.valueProperty]
  state.selectorOpen = false
}
</script>

<template>
  <div ref="container" class="container" @click="state.selectorOpen = true">
    <slot name="display" :item="valueItem">
      <span class="value">{{ valueDisplay }}</span>
    </slot>
    <div ref="selectorOverlay" :class="{ 'selector-overlay': true, open: state.selectorOpen }">
      <div class="search-field">
        <input
          type="text"
          :placeholder="filterPlaceholder || undefined"
          inputmode="none"
          v-model="filter"
          ref="filterInput"
        />
        <button type="button" @click.stop="state.selectorOpen = false">🅧</button>
      </div>
      <div class="settings">
        <slot name="settings"></slot>
      </div>
      <div class="items-list">
        <div
          v-for="item in items"
          :key="item[props.valueProperty]"
          class="selector-item"
          :class="{ selected: item[props.valueProperty] === model }"
          @click.stop="onSelect(item)"
        >
          <slot name="item" :item="item" :selected="item[props.valueProperty] === model"></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="css" scoped>
.container {
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
    position: fixed;
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
    overflow: hidden;

    display: flex;
    flex-direction: column;
    align-items: start;
    gap: 0.2rem;
    padding: 0.5rem;
    transform: translateY(150%);

    transition: opacity 0.2s ease-in-out, transform 0.2s ease-in-out;

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
    }
  }
}
</style>
