package main

import (
	"strconv"
	"syscall/js"
)

type callback func(js.Value, []js.Value) interface{}
type operation func(int, int) int

func buildCallback(op operation) callback {
	return func(this js.Value, i []js.Value) interface{} {
		res := op(getValue(i[0]), getValue(i[1]))
		js.Global().Get("document").Call("getElementById", i[2].String()).Set("value", res)
		return nil
	}
}

func getValue(obj js.Value) int {
	value := js.Global().Get("document").Call("getElementById", obj.String()).Get("value").String()
	result, _ := strconv.Atoi(value)
	return result
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

func divide(a, b int) int {
	return a / b
}

func multiply(a, b int) int {
	return a * b
}

func registerCallbacks() {
	js.Global().Set("add", js.FuncOf(buildCallback(add)))
	js.Global().Set("subtract", js.FuncOf(buildCallback(subtract)))
	js.Global().Set("divide", js.FuncOf(buildCallback(divide)))
	js.Global().Set("multiply", js.FuncOf(buildCallback(multiply)))
}

func main() {
	c := make(chan struct{}, 0)
	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
