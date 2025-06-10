const task_name = document.getElementById("task_name")
const task_des = document.getElementById("task_des")
const sub = document.getElementById("des");
const urlParams = new URLSearchParams(window.location.search);
const id = urlParams.get('id');

task_information(id);

function render_task(task) {
    task_name.innerHTML = task.id;
    task_des.innerHTML = task.title; 
}

const sub_btn = document.getElementById("sub_btn")
const test_btn = document.getElementById("test_btn")
sub_btn.addEventListener('click', e => {
    const text = document.getElementById("task_text");
    const text_value = text.value.trim();

    if (text_value === "") {
        return;
    } else {
        console.log(text_value);
        const text_object = { text: text_value };
        console.log(text_object);
        submit_req(text_object, id);

    }
});

function sub_render(subs) {
    sub.innerHTML = '';
    for(let i = 0; i < subs.length; i++) {
        sub.innerHTML = subs[i].status;
        // ...
    }
}
