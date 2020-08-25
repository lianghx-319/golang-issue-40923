package main

import (
	"golang-issues-40923/src/process"
	"syscall/js"
)

var (
	global = js.Global().Get("Go").Get("exports")
)

func log(msg string) {
	js.Global().Get("console").Call("log", msg)
}

func main() {
	ch := make(chan bool)

	log("Wasm Start")

	traceProcessor := process.GetTraceCompareController()

	test := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log(args[0].Get("0").String())
		return nil
	})
	defer test.Release()

	parse := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//log("parsing...")
		params, callback := args[0], args[1]
		var (
			groupType = params.Get("0").Int()
			isEoF     = params.Get("1").Bool()
			length    = params.Get("3").Int()
		)
		chunk := make([]byte, length)
		js.CopyBytesToGo(chunk, args[0].Get("2"))
		finish := traceProcessor.Parse(chunk, isEoF, groupType)
		//log("chunk parsed")
		callback.Invoke(js.Null(), finish)
		return nil
	})
	defer parse.Release()

	getCompareTable := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params, callback := args[0], args[1]
		var (
			functionName = params.Get("0").String()
			threshold    = params.Get("1").Int()
			changed      = params.Get("2").Bool()
			added        = params.Get("3").Bool()
			deleted      = params.Get("4").Bool()
			field        = params.Get("5").Int()
			direction    = params.Get("6").Bool()
			current      = params.Get("7").Int()
			pageSize     = params.Get("8").Int()
		)

		log("getCompareTable")
		string := traceProcessor.GetTableResult(functionName, threshold, changed, added, deleted, field, direction, current, pageSize)
		callback.Invoke(js.Null(), string)
		return nil
	})
	defer getCompareTable.Release()

	global.Set("hello", test)
	global.Set("parse", parse)
	global.Set("getJSON", getCompareTable)

	<-ch
}
