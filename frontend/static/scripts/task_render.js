const task_name = document.getElementById("task_name")
const task_des = document.getElementById("task_des")

function render_task(task) {
    console.log(task)
    task_name.innerHTML = task.id;
    task_des.innerHTML = task.title; 
}

const sub_btn = document.getElementById("sub_btn")
const test_btn = document.getElementById("test_btn")
sub_btn.addEventListener('click' , e => {
    const text = document.getElementById("task_text")
    const text_value = text.value;
    if(text_value == "") {} 
    else {
        console.log(text_value); //TODO: submit task
    }
})

test_btn.addEventListener('click' , e => {
    const text = document.getElementById("task_text")
    const text_value = text.value;
    if(text_value == "") {} 
    else {
        console.log(text_value); //TODO: test task
    }
})