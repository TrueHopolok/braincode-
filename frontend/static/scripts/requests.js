/* 
The possible requests that frontend need to handle: In header:

    ANY - session-token as "Session" parameter
    GET - content-type as "Content-Type" parameter In URL:
    GET(json) - page as "page" parameter
    GET(task) - task-id as "id" parameter
    DELETE(task) - task-id as "id" parameter
    ANY - language as "lang" parameter

The responses gurantee:

    Status code
    Content-type
    Session (if authorized)
*/

const tasks = [
    { id: 1, title: "Hello", author: "1" },
    { id: 2, title: "Test", author: "2" },
    { id: 3, title: "1234", author: "3" },
    { id: 4, title: "iiii", author: "4" }
];

function task_req() {
    fetch("https://jsonplaceholder.typicode.com/todos/1", {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Session': session,
            //        'lang': lang 
        }
    })
        .then(response => response.json())
        .then(data => {
            // tasks = data;
        });
}
//.then(data => tasks.push(data));

let session = localStorage.getItem('sessionToken');

// Task
function task_information(id) {
    fetch('/task/', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Session': session,
            'lang': lang,
            'id': id
        }
    })
        .then(response => response.json())
        .then(data => {
            render_task(data);
        });
}

// Upload Task



// Submit and test
function submit_req(data, id) {
    fetch('/task/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Session': session,
            'id': id
        },
        body: JSON.stringify(data)
    }).then(response => response.json())
        .then(data => {
            const answer = document.getElementById("task_answer");
            answer.innerHTML = data;
        })
}

/*function test_req(data) {
    fetch('', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Session': session,
            'lang': lang,
            'id': id
        },
        body: JSON.stringify(data)
    }).then(response => response.json())
}*/

// Profile

function profile_req() {
    fetch('/stats/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
    })
        .then(response => response.json())
        .then(data => {
            render_profile(data);
        });
}
