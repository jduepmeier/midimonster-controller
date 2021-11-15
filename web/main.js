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