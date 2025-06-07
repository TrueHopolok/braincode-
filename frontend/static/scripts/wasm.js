const go = new Go()

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(result => {
    go.run(result.instance);
    const parseResult = parseMarkleft(".task = ASDDS");

    // schema:
    // {
    //     source: string (source passed to parseMarkleft)
    //     formatted_source: string (formatted markleft source)
    //     html: string (rendered html)
    //     locale: selected locale ("en" | "ru" | "")
    //     errors: string[]
    // }

    console.log(parseResult)
    document.getElementById("doc").innerHTML = parseResult.html
})
