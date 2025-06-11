profile_req();

const sub_list = document.getElementById("sub_list");

function render_profile(data) {
    console.log(data);
    if (data === null || data.TotalAmount === 0 || data.Rows === null) {
        sub_list.innerHTML = 'No submissions found...';
        return;
    }
    sub_list.innerHTML = '';
    let isEnglish = document.LANG !== 'ru'
    data.Rows.forEach(sub => {
        let node = document.createElement("div");
        node.classList.add("submission");
        let id = sub.TaskId.Int64;
        if (!sub.TaskId.Valid) {
            id = "deleted";
        }
        const title = isEnglish ? sub.TitleEn.String : (sub.TitleRu.String || sub.TitleEn.String);
        node.innerHTML = `
        <div class="timestamp">${sub.Timestamp}</div>
        <div class="task-id">${id}. ${title}</div>
        <div class="score">Score: ${sub.Score}</div>
        `;
        sub_list.appendChild(node);
    });
}