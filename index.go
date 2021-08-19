package midimonster

const indexHTML = `
<html>
	<head>
		<title>Midimonster Controller</title>
		<link rel="stylesheet" href="main.css" />
	</head>
	<body>
		<textarea col="80" rows="40"></textarea>
		<button onscript="writeConfig()">Write Config</button>
		<button onscript="restart()">Restart Midimonster</button>
		<button onscript="get()">Get newest config</button>

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
async function writeConfig() {}

async function restart() {}

async function get() {
	const response = await fetch('api/config');
	const myJson = await response.json();
	textArea = document.querySelector("textarea");
	textArea.value = myJson.content;
}
get();
`
