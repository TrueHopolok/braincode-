const problems = document.getElementById("section-content");
const search = document.getElementById("section-search");
const checkbox = document.getElementById("user_problem");

document.addEventListener("DOMContentLoaded", (event) => {
	let data = fetch("/", {
		method: "GET",
		headers: {
			"Content-Type": "application/json",
		},
	}).then(response => response.json());
	console.log(data);
	console.log(problems);
	console.log(search);
	console.log(checkbox);
	// TODO: add render
});

