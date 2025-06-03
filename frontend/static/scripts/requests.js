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


// Session-token
let loginField = document.getElementById("login_form");

loginField.addEventListener('submit', function(event){
    const username = this.uname.value; // или this.elements['uname'].value
    const password = this.psw.value;
    const remember = this.remember.checked;

    fetch("?", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username: username,
            password: password
        })
    })
    .then(response => {
        const sessionToken = response.headers.get('Session');

        if(sessionToken) {
            localStorage.setItem('sessionToken', sessionToken);
        }
    })
})





// Page
