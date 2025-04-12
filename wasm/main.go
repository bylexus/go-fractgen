//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/bylexus/go-fract/lib"
)

func updateCounter() {
	actVal := 0
	go func() {
		actVal += 1
		js.Global().Call("setCounterValue", actVal)
	}()
}

func main() {
	fmt.Println("Hello, WASM!")

	var firstFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("firstFn")
		counter := 0
		go (func() {
			for {
				time.Sleep(time.Millisecond * 1000)
				counter += 1
				fmt.Println(counter)
			}
		})()
		return 42
	})

	var calcFract = js.FuncOf(func(this js.Value, calcFractArgs []js.Value) interface{} {
		fp := calcFractArgs[0].String()
		fractPreset := lib.FractalPreset{}
		json.Unmarshal([]byte(fp), &fractPreset)
		fmt.Printf("%#v\n", fractPreset)

		cp := calcFractArgs[1].String()
		colorPreset := lib.ColorPreset{}
		json.Unmarshal([]byte(cp), &colorPreset)
		fmt.Printf("%#v\n", colorPreset)

		// jsArr := calcFractArgs[1]
		// arr := make([]byte, jsArr.Length())
		// js.CopyBytesToGo(arr, jsArr)
		// arr[0] = 42
		// js.CopyBytesToJS(jsArr, arr)

		// Handler for the Promise: this is a JS function
		// It receives two arguments, which are JS functions themselves: resolve and reject
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			// Commented out because this Promise never fails
			//reject := args[1]

			// Now that we have a way to return the response to JS, spawn a goroutine
			// This way, we don't block the event loop and avoid a deadlock
			go func() {
				// colorPreset := lib.ColorPreset{
				// 	Name: "test",
				// 	Palette: lib.ColorPalette{
				// 		color.RGBA{R: 255, G: 0, B: 0, A: 255},
				// 		color.RGBA{R: 0, G: 255, B: 0, A: 255},
				// 		color.RGBA{R: 0, G: 0, B: 255, A: 255},
				// 	},
				// }

				f, _ := lib.NewFractalFromPresets(colorPreset, fractPreset)
				img := f.CalcFractalImage(nil)
				// fmt.Printf("%#v\n", *img)
				fmt.Printf("image widthxheight: %dx%d\n", fractPreset.ImageWidth, fractPreset.ImageHeight)
				pix := img.Pix
				fmt.Printf("image size: %d\n", len(img.Pix))
				jsArr := js.Global().Get("Uint8Array").New(fractPreset.ImageWidth * fractPreset.ImageHeight * 4)
				fmt.Printf("%#v\n", jsArr)
				js.CopyBytesToJS(jsArr, pix)
				// fmt.Printf("%#v\n", jsArr)

				// Resolve the Promise, passing anything back to JavaScript
				// This is done by invoking the "resolve" function passed to the handler
				resolve.Invoke(jsArr)
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)

	})

	js.Global().Set("firstFn", firstFn)
	js.Global().Set("calcFract", calcFract)

	<-make(chan bool)
}
