profile_req();

const sub_year = document.getElementById("sub_year");
const sub_month = document.getElementById("sub_month");
const sub_all = document.getElementById("sub_all");
const sub_list = document.getElementById("sub_list");
const sub_rate = document.getElementById("sub_rate");

function render_profile(data) {
    sub_month.innerHTML = `${data.month} problems`;
    sub_all.innerHTML = `${data.all} problems`;
    sub_year.innerHTML = `${data.year} problems`;

    sub_list.innerHTML = `${data.list} problems`;
    sub_rate.innerHTML = `${data.rate} problems`;
}