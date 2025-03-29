import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import { loadPresets } from './lib/use-presets'

const boot = async () => {
  // Pre-load presets: They need to be present on app start:
  await loadPresets()
  createApp(App).mount('#app')
}

boot()
