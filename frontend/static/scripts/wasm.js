let task_object;
function getObject(object) {
    task_object = object;
}

const go = new Go()

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(result => {
    go.run(result.instance);
    const parseResult = parseMarkleft(`.task = ${task_object}`);

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
