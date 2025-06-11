const upload_form = document.getElementById('task_upload');

upload_form.addEventListener('submit', e => {
    e.preventDefault();
    const form = e.target;
    task_upload(form.des.value);
})