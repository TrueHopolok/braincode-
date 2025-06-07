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
  { id: 1, title: "Hello" },
  { id: 2, title: "Test" },
  { id: 3, title: "1234" },
  { id: 4, title: "iiii" }
];

fetch("https://jsonplaceholder.typicode.com/todos/1", {
    method: 'GET',
    headers: {
        'Content-Type': 'application/json',
        'Session': localStorage.getItem('sessionToken'),
        'lang': 'ru' 
    }
})
.then(response => response.json())
.then(data => {
    problems_render();
});
//.then(data => tasks.push(data));


// Task
function task_information(id) {
    fetch("https://jsonplaceholder.typicode.com/todos/1", {
            method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Session': localStorage.getItem('sessionToken'),
            'lang': 'ru',
            'id': id
        }
    })
    .then(response => response.json())
    .then(data => {
        render_task(data);
    });
}

// Submit and test

function submit_req(data) {
    fetch('', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    }).then(response => response.json())
}

function test_req(data) {
    fetch('', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    }).then(response => response.json())
}