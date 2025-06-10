function logout() {
    fetch('/login/', {
        method: 'DELETE'
    });
}