const pageSize = 2;

const problemset_content = document.getElementById("problemset-content");
const problems = document.getElementById("problems");
const checkbox = document.getElementById("user_problem");

function arrowRender() {
    const left_arrow = document.createElement("a")
    left_arrow.textContent = "<- ";
    left_arrow.href = "#";
    problemset_content.appendChild(left_arrow);

    for(let i = 0; i < tasks.length / 2; i++) {
        const num_link = document.createElement("a");
        num_link.textContent = " " + (i + 1);
        num_link.href = "#";
        num_link.className = "page"
        problemset_content.appendChild(num_link);
    }
    const right_arrow = document.createElement("a");
    right_arrow.textContent = " ->";
    right_arrow.href = "#";
    problemset_content.appendChild(right_arrow);
}

function showPage(page) {
    problems.innerHTML = '';
    const start = (page - 1) * pageSize;
    const end = start + pageSize;
    
    const currentTasks = tasks.slice(start, end);

    for(let i = 0; i < currentTasks.length; i++) {
        const task = document.createElement("a");
        task.innerHTML += currentTasks[i].id + '.' + currentTasks[i].title + '<br>' + 'by:' + currentTasks[i].author;
        task.href = `taskpage.html?id=${currentTasks[i].id}`;
        task.className = "task";
        problems.appendChild(task);
        problems.appendChild(document.createElement("br"));
    }
}

problemset_content.addEventListener('click', e => {
    if (e.target.classList.contains("page")) {
        const currentPage = e.target.textContent;
        showPage(currentPage);
    }
});

function problems_render() {
    showPage(1);
    arrowRender();
}

task_req();
problems_render();