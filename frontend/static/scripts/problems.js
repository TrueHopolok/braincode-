let problems, search, filter;

document.addEventListener("DOMContentLoaded", (event) => {
	problems = document.getElementById("section-content");
	search = document.getElementById("search-input");
	filter = document.getElementById("filter-input");
	getProblemSet();
	search.addEventListener("oninput", getProblemSet()); // TODO(vadim) not updating without reload of the page
	filter.addEventListener("onclick", getProblemSet()); // TODO(vadim) not working
});

function getProblemSet() {
	let params = new URLSearchParams(document.location.search);
	let lang = params.get("lang");
	if (lang === null) {
		lang = "EN";
	}
	let isUserProblems = filter.checked;
	if (isUserProblems) {
		isUserProblems = "user-only";
	}
	fetch(`/?lang=${lang}`, {
		method: "GET",
		headers: {
			"Content-Type": "application/json",
			"Search": search.value,
			"Filter": isUserProblems,
		},
	})	.then(response => response.json().then(data => renderTasks(data)))
		.catch(error => failedRequest(error));
}
	
function createEmptyNode() {
	let node = document.createElement("p");
	node.innerHTML = "No task found...";
	return node;
}

function createTaskNode(task) {
	let upperNode = document.createElement("li");
	let node = document.createElement("a");
	node.innerHTML = task.Id + ". " + task.Title + '<br>' + 'by: ' + task.OwnerName;
	node.href = `/task/?id=${task.Id}`;
	node.className = "task";
	upperNode.appendChild(node);
	return upperNode;
}

function failedRequest(error) {
	console.error(error);
	problems.innerHTML = ""
	
	problems.appendChild(createEmptyNode());
}

function renderTasks(data) {
	problems.innerHTML = ""

	if (data.TotalAmount == 0) {
		problems.appendChild(createEmptyNode());
		return;
	} else if (data.Rows.length == 0) {
		problems.appendChild(createEmptyNode());
		return;
	}

	let problemlist = document.createElement("ul");
	problems.appendChild(problemlist);
	data.Rows.forEach(task => {
		problemlist.appendChild(createTaskNode(task));
		problemlist.appendChild(document.createElement("br"))
	});
	
	// TODO(vadim) add page selection render
	// TODO(vadim) add search and filter connection
}
