package main

import (
	"syscall/js"
)

func main() {
	window := js.Global()
	doc := window.Get("document")
	doc.Set("title", "Welcome page")
	header := doc.Call("createElement", "h1")
	header.Set("innerText", "oh hi mark")
	body := doc.Get("body")
	body.Call("appendChild", header)
}
