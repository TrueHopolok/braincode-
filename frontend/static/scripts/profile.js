profile_req();

const sub_year = document.getElementById("sub_year");
const sub_month = document.getElementById("sub_month");
const sub_all = document.getElementById("sub_all");
const sub_list = document.getElementById("sub_list");
const sub_rate = document.getElementById("sub_rate");
const problems = document.getElementById("problems");
const rate = document.getElementById("sub_rateто");

function render_profile(data) {
    sub_month.innerHTML = `${data.month} problems`;
    sub_all.innerHTML = `${data.all} problems`;
    sub_year.innerHTML = `${data.year} problems`;

    sub_list.innerHTML = `${data.list} problems`;
    sub_rate.innerHTML = `${data.rate} problems`;
}

function render_usertask(data) {
    problems.innerHTML = '';
    const start = (page - 1) * pageSize;
    const end = start + pageSize;
    
    const currentTasks = tasks.slice(start, end);

    for(let i = 0; i < currentTasks.length; i++) {
        const task = document.createElement("a");
        task.innerHTML += data[i].id + '.' + data[i].title + '<br>' + 'by:' + data[i].author;
        task.href = `taskpage.html?id=${currentTasks[i].id}`;
        task.className = "task";
        problems.appendChild(task);
        problems.appendChild(document.createElement("br"));
    }
}