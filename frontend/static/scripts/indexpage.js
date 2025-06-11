{
    function render_task_list(elementId, data) {
        let isEnglish = document.LANG !== 'ru'

        console.log("renderTaskList", elementId, data, isEnglish)
        const container = document.getElementById(elementId);
        if (!container) {
            console.error(`Element with ID "${elementId}" not found`);
            return;
        }

        container.innerHTML = '';
        container.classList.add('task-list-container');

        const list = document.createElement('ul');
        list.classList.add('task-list');
        container.appendChild(list);

        if (data.Rows.length === 0) {
            const listItem = document.createElement('li');
            listItem.textContent = isEnglish ? "No tasks found" : "Задачи не найдены"
        } else {
            data.Rows.forEach(task => {
                const listItem = document.createElement('li');
                listItem.classList.add('task-item');

                const title = isEnglish ? task.TitleEn : (task.TitleRu || task.TitleEn);
                const scoreDisplay = task.Score.Valid ? task.Score.Float64 : (isEnglish ? 'Not solved yet' : 'Ещё не решено');

                // Create link element
                const taskLink = document.createElement('a');
                taskLink.href = `/task/?id=${task.Id}`;
                taskLink.classList.add('task-link');

                taskLink.innerHTML = `
                <div class="task-info">${task.Id}. <b>${title}</b></div>
                <div class="task-secondaries">
                <div class="task-author">${isEnglish ? 'By' : 'От'} ${task.OwnerName || (isEnglish ? "[DELETED]" : "[УДАЛЁН]")}</div>
                ${document.AUTH ? `<div class="task-score">${isEnglish ? "Score" : "Результат"}: <i>${scoreDisplay}</i></div>` : ""}
                </div>
            `;


                listItem.appendChild(taskLink);
                const canDelete = (document.ISADMIN || task.OwnerName === document.USERNAME) && document.AUTH
                if (canDelete) {
                    deleteButton = document.createElement("button");
                    const text = isEnglish ? "Delete task" : "Удалить задачу";
                    deleteButton.innerHTML = task.OwnerName === document.USERNAME ? text : "[A] " + text;
                    deleteButton.classList.add("task-delete");
                    deleteButton.addEventListener("click", () => delete_task(task.Id));
                    listItem.appendChild(deleteButton);
                }
                list.appendChild(listItem);
            });
        }
    }

    let currentPage = 0;

    async function get_tasks() {
        let page = currentPage;
        let current_user_only = document.getElementById("current-user-only-checkbox")?.checked ?? false;
        let query = document.getElementById('tasks-search')?.value ?? "";

        let url = `/api/tasks/?page=${page}`;
        if (query) {
            url += `&query=${encodeURIComponent(query)}`
        }
        if (current_user_only) {
            url += "&current-only=1";
        }

        return (await fetch(url)).json()
    }

    async function delete_task(id) {
        await fetch(
            "/",
            {
                method: "DELETE",
                headers: {
                    "Id": id,
                },
            },
        )
        tasks_search()
    }

    function tasks_next() {
        if (currentPage < data.TotalPages - 1) {
            currentPage++;
            get_tasks().then(d => {
                data = d;
                render_task_list("section-content", d);
            });
        }
    }

    function tasks_prev() {
        if (currentPage > 0) {
            currentPage--;
            get_tasks().then(d => {
                data = d;
                render_task_list("section-content", d);
            });
        }
    }

    function tasks_search() {
        const searchInput = document.getElementById('tasks-search').value;
        const query = searchInput.value;
        get_tasks(0, false, query).then(d => {
            data = d;
            currentPage = 0;
            render_task_list("section-content", d);
        });
    }

    get_tasks().then(d => {
        data = d
        render_task_list("section-content", d)
    })

    document.getElementById("tasks_next").addEventListener("click", tasks_next)
    document.getElementById("tasks_prev").addEventListener("click", tasks_prev)
    document.getElementById("current-user-only-checkbox")?.addEventListener("click", tasks_search)
    let debounceTimer;
    document.getElementById("tasks-search").addEventListener("input", function () {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            tasks_search()
        }, 300);
    });
}