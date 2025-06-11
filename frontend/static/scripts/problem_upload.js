const go = new Go()

// parseMarkleft: (source: string, locale: string?) -> Schema
// Schema:
// {
//     source: string (source passed to parseMarkleft)
//     formatted_source: string (formatted markleft source)
//     html: string (rendered html)
//     locale: selected locale ("en" | "ru" | "")
//     errors: string[]
// }
const parseML = WebAssembly.instantiateStreaming(fetch("/static/wasm/markleft.wasm"), go.importObject).then(result => {
  go.run(result.instance);
  return parseMarkleft
})

const inputArea = document.getElementById("problem_statement")
const enRadio = document.getElementById("radio-en")
const ruRadio = document.getElementById("radio-ru")
const preview = document.getElementById("task-preview")
const previewErrors = document.getElementById("preview-errors")
const formatButton = document.getElementById("format-button")

async function renderPreview() {
  const parse = await parseML;
  const source = inputArea.value;
  const locale = enRadio.checked ? "en" : "ru";

  try {
    const res = parse(source, locale);
    preview.innerHTML = res.html;
    previewErrors.innerText = res.errors.join("\n");
  } catch (err) {
    preview.innerHTML = "<p><strong>Preview error</strong></p>";
    previewErrors.innerText = String(err);
  }
}

async function formatSource() {
  const parse = await parseML;
  const source = inputArea.value;
  const locale = enRadio.checked ? "en" : "ru";

  try {
    const res = parse(source, locale);
    preview.innerHTML = res.html;

    const e = res.errors.join("\n")
    previewErrors.innerText = e ? "Errors:\n" + e : "";
    inputArea.value = res.formatted_source;
  } catch (err) {
    preview.innerHTML = "<p><strong>Format error</strong></p>";
    previewErrors.innerText = String(err);
  }
}

let debounceTimer;
inputArea.addEventListener("input", () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(renderPreview, 300);
});

formatButton.addEventListener("click", formatSource);
enRadio.addEventListener("change", renderPreview);
ruRadio.addEventListener("change", renderPreview);
