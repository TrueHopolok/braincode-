const upload_form = document.getElementById('task_upload');

upload_form.addEventListener('submit', e => {
    e.preventDefault();
    const form = e.target;

    fetch('/upload/', {
        method: 'POST',
        body: form.des.value,
    }).then(resp => {
        if (resp.redirected) {
            window.location.hred = resp.url
        }
    })
})