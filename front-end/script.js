function ShowDate() {
    var date = new Date();
    document.getElementById("dateparagraph").innerText = date;
}

document.addEventListener("DOMContentLoaded", function() {
    ShowDate();
    document.getElementById("datebutton").onclick = ShowDate;
});
