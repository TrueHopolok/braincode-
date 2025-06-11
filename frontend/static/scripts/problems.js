document.addEventListener("DOMContentLoaded", (event) => {
	let params = new URLSearchParams(document.location.search);
	let lang = params.get("lang");
	if (lang === null) {
		lang = "EN";
	}
	fetch(`/?lang=${lang}`, {
		method: "GET",
		headers: {
			"Content-Type": "application/json",
		},
	})	.then(response => response.json().then(data => renderTasks(data)))
		.catch(error => failedRequest(error));
	});
	
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
	const problems = document.getElementById("section-content");
	console.log(problems); 	// TODO: remove console logs probably
	problems.appendChild(createEmptyNode());
}

function renderTasks(data) {
	const problems = document.getElementById("section-content");
	const search = document.getElementById("search-input");
	const checkbox = document.getElementById("filter-input");
	console.log(problems); 	// TODO: remove console logs probably
	console.log(search); 	// TODO: remove console logs probably
	console.log(checkbox); 	// TODO: remove console logs probably
	console.log(data.Rows); // TODO: remove console logs probably

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
	
}
