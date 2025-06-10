const pageSize = 5;

const problemset_content = document.getElementById("problemset-content");
const problems = document.getElementById("problems");

const left_arrow = document.createElement("a")
left_arrow.textContent = "<- ";
left_arrow.href = `#`;
problemset_content.appendChild(left_arrow);

for(let i = 0; i < size / pageSize; i++) {
    const num_link = document.createElement("a");
    num_link.textContent = " " + (i + 1);
    num_link.href = `?page=${i+1}`;
    num_link.className = "page"
    problemset_content.appendChild(num_link);
}
const right_arrow = document.createElement("a");
right_arrow.textContent = " ->";
right_arrow.href = "#";
problemset_content.appendChild(right_arrow);




function showPage(tasks) {
    problems.innerHTML = '';
    for(let i = 0; i < tasks.length; i++) {
        const task = document.createElement("a");
        task.innerHTML += tasks[i].id + '.' + tasks[i].title + '<br>' + 'by:' + tasks[i].author;
        task.href = `taskpage.html?id=${tasks[i].id}`;
        task.className = "task";
        problems.appendChild(task);
        problems.appendChild(document.createElement("br"));
    }
}

problemset_content.addEventListener('click', e => {
    if (e.target.classList.contains("page")) {
        console.log(currentPage);
    }
});