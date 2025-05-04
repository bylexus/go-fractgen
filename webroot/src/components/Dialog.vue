<script setup lang="ts">
const show = defineModel<boolean>()
</script>

<template>
  <div :class="{ dialog: true, hidden: !show }">
    <button class="close-btn" type="button" @click="show = false">X</button>
    <div class="dialog-content">
      <slot></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.dialog {
  position: fixed;
  bottom: 0;
  left: 0;
  z-index: 2;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(10px);
  transition:
    opacity 0.2s ease-in-out,
    bottom 0.2s ease-in-out;
  &.hidden {
    opacity: 0;
    bottom: -100%;
  }

  .close-btn {
    position: absolute;
    top: 0.2rem;
    right: 0.2rem;
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
  .dialog-content {
    padding-top: 2rem;
    height: 100%;
  }
}
</style>
