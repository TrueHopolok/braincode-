function logout(event) {
    event.preventDefault()
    fetch('/logout/', {
        method: 'POST'
    }).then(v => {
        window.location.replace("/");
    });
}