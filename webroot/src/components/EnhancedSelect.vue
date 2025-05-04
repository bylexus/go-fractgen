<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'

const model = defineModel<any>('value')
const filter = defineModel<any>('filter')

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
    if (newVal && !oldVal) {
      filter.value = ''
      filterInput.value?.focus()
    }
  },
)

const valueItem = computed(() => {
  return props.items.find((item: any) => item[props.valueProperty] === model.value)
})

const valueDisplay = computed(() => {
  return valueItem.value ? valueItem.value[props.displayProperty] : ''
})

function onSelect(item: any) {
  console.log('item selected: ', item)
  model.value = item[props.valueProperty]
  state.selectorOpen = false
}
</script>

<template>
  <div class="container" @click="state.selectorOpen = true">
    <slot name="display" :item="valueItem">
      <span class="value">{{ valueDisplay }}</span>
    </slot>
    <div :class="{ 'selector-overlay': true, open: state.selectorOpen }">
      <div class="search-field">
        <input
          type="text"
          :placeholder="filterPlaceholder || undefined"
          v-model="filter"
          ref="filterInput"
        />
        <button type="button" @click.stop="state.selectorOpen = false">🅧</button>
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
  position: relative;
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
    position: absolute;
    width: 400px;
    max-width: 100vw;
    height: 400px;
    max-height: 400px;
    opacity: 0;
    left: -100vw;
    bottom: 0;
    z-index: 2;
    background-color: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(10px);
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
    border: 1px solid #aaa;
    border-radius: 3px;
    overflow: none;

    display: flex;
    flex-direction: column;
    align-items: start;
    gap: 0.2rem;
    padding: 0.5rem;

    transition: opacity 0.2s ease-in-out, left 0.2s ease-in-out;

    &.open {
      opacity: 1;
      left: 0;
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
