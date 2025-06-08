const login_form = document.getElementById("login_form");
const registration_form = document.getElementById("reg_form")

login_form.addEventListener('submit', e => {
    e.preventDefault();
    const form = e.target;
    login_req(form.uname.value, form.psw.value, form.remember.value);
})

registration_form.addEventListener('submit', e => {
    e.preventDefault();
    const form = e.target;
    reg_req(form.uname.value, form.psw.value);
})
