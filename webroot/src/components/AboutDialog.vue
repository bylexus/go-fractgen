<script setup lang="ts">
import Dialog from './Dialog.vue'

</script>

<template>
  <Dialog>
    <div class="about-dialog">
      <h1>About fractgen</h1>
      <p class="small">&copy; 2012-2025 <a href="https://alexi.ch">Alex Schenkel</a></p>
      <p>fractgen is a Mandelbrot / Julia fractal generator <a href="https://github.com/bylexus/go-fractgen">written in Go</a>.</p>
      <p>This is a hobby project of mine, and not intended for accuracy or scientific use. I just very like the idea of
        mathematically-based image generation.
      </p>
      <h3>Some features:</h3>
      <p>
        <ul>
          <li>generate Mandelbrot and Julia fractals</li>
          <li>Written in <a href="https://golang.org">Go</a></li>
          <li>Web Server / interactive Web Interface, CLI</li>
          <li>Start a local web server instance</li>
          <li>precision: float64. The smallest feasible diameter is about 1e-13, limited by the floating point precision in Go (higher precision is in work)</li>
          <li>Use (own) Fractal and Color presets</li>
          <li>Create a "flight" through a fractal</li>
        </ul>
      </p>
      <p>You can find the source code on <a href="https://github.com/bylexus/go-fractgen">GitHub</a>. The code is released under the MIT license.</p>
      <p>Have Fun!</p>
      <p><img src="https://fractgen.alexi.ch/fractal-image/jpg?width=640&height=480&name=Mandelbrot%20Lexus%27%20Secret%20Spot&iterFunc=Mandelbrot&diameterCX=0.0000514984130859375&centerCX=-0.2379717670706691&centerCY=-0.8451734844686328&colorPreset=sedum&juliaKi=0.6&juliaKr=-0.6&maxIterations=1000&colorPaletteLength=256&colorPaletteRepeat=2&colorPaletteReverse=true"></p>
    </div>
  </Dialog>
</template>

<style lang="scss" scoped>
.about-dialog {
  color: white;
  a {
    color: white;
  }
  padding: 0.5rem;
  overflow: scroll;
  height: 100%;

  img {
    max-width: 100%;
    border-radius: 10px;
  }
}

</style>
