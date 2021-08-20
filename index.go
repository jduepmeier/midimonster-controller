package midimonster

const indexHTML = `
<html>
	<head>
		<title>Midimonster Controller</title>
		<link rel="stylesheet" href="main.css" />
	</head>
	<body>
		<textarea id="config" col="80" rows="40"></textarea>
		<button id="write" onclick="writeConfig()">Write Config</button>
		<button id="restart" onclick="restart()">Restart Midimonster</button>
		<button id="get" onclick="get()">Get newest config</button>

		<br />
		<div>
			<textarea id="log">
			</textarea>
		</div>

		<script src="main.js"></script>
	</body>
</html>
`

const mainCSS = `
html {
	padding: 0;
	margin: 0;
}

body: {
	width: 80%;
	margin: 0 auto;
}

button {
	margin: 10 auto;
	background-color: #222222;
	color: #fefefe;
}

textarea {
	width: 100%;
	clear: both;
}
`

const mainJS = `

function handleError(err) {
	if (err !== undefined) {
		log = document.querySelector("#log");
		log.value += "\nERROR: " + myJson["Error"];
	}
}

async function writeConfig() {
  config = document.querySelector("#config")
  fetch('api/write', {
    method: 'POST',
    body: JSON.stringify({
		"Content": config.value,
	}),
    headers: {
      'Content-Type': 'application/json'
    }
  }).then((response) => {
	response.json().then((myJson) => {
		handleError(myJson["Error"]);
	});
  }).catch((err) => {
	  handleError(err);
  });
}

async function restart() {
  fetch('api/reload', {
    method: 'POST',
    body: JSON.stringify({}),
    headers: {
      'Content-Type': 'application/json',
    }
  }).then((response) => {
	response.json().then((myJson) => {
		handleError(myJson["Error"]);
	});
  }).catch((err) => {
	  handleError(err);
  });
}

async function get() {
	const response = await fetch('api/config');
	const myJson = await response.json();
	textArea = document.querySelector("#config");
	textArea.value = myJson.Content;
}

async function init() {
	log = document.querySelector("#log");
	config = document.querySelector("#config");
	get();
}

init();
`
