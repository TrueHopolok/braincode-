const PROBLEM_SET = document.getElementById('problem-set');
const PROBLEM_UPLOAD = document.getElementById('problem-upload');
const SECTION_UPPER = document.getElementById('section-upper');
const SECTION_CONTENT = document.getElementById('section-content')

tasks = [
    ["example1"],
    ["example2"],
    ["example3"]
]

function render_task() {
    let table = document.createElement('table');
    for(let i = 0; i < tasks.length; i++) {
        let table_row = document.createElement('tr');
        let table_col1 = document.createElement('th');  
        let table_col2 = document.createElement('th');
        let table_col3 = document.createElement('th');

        table_col1.innerHTML = i + 1;
        table_col2.innerHTML = tasks[i][0];

        let link = document.createElement('a');
        link.innerHTML = 'task ->';
        link.href = 'taskpage.html';
        table_col3.append(link);

        table_row.append(table_col1);
        table_row.append(table_col2);
        table_row.append(table_col3);
        table.append(table_row);
    }

    SECTION_CONTENT.append(table);
}

function hide() {
    SECTION_UPPER.innerHTML = '';
    SECTION_CONTENT.innerHTML = '';
}

render_task();

PROBLEM_SET.addEventListener("click", (e) => {
    hide();
    
    let input = document.createElement('input');
    input.className = 'section-search';
    input.type = 'text';
    
    let section_div = document.createElement('div');
    section_div.className = 'section-user';

    let div_checkbox = document.createElement('input');
    let div_text = document.createElement('div');
    div_checkbox.type = 'checkbox';
    div_text.innerHTML = "USER PROBLEM";
    section_div.append(div_checkbox);
    section_div.append(div_text);

    SECTION_UPPER.append(input);
    SECTION_UPPER.append(section_div);

    render_task();

});

PROBLEM_UPLOAD.addEventListener("click", (e) => {

    hide();
    let task_form = document.createElement('form');

    let taskname_label = document.createElement('label');
    taskname_label.innerHTML = 'Task Name';

    let taskname_input = document.createElement('input');
    taskname_input.type = 'text';

    task_form.append(taskname_label);
    task_form.append(taskname_input);
    

    SECTION_CONTENT.append(task_form)
});
