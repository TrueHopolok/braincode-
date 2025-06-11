const upload_form = document.getElementById('task_upload');
const status_text = document.getElementById('status_text');

upload_form.addEventListener('submit', e => {
    e.preventDefault();
    const form = e.target;
    console.log(form.des.value);
    fetch('/upload/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: (form.des.value),
    })
    .then(resp => showSuccess())
    .catch(error => showError(error));
});

function showSuccess() {
    status_text.innerHTML = "";
    let node = document.createElement("p");
    node.innerHTML = "Success in uploading the task!";
    status_text.appendChild(node);
}

function showError(error) {
    status_text.innerHTML = "";
    let node = document.createElement("p");
    node.innerHTML = `Failed upload.<br>Error:${error}`;
    status_text.appendChild(node);
}