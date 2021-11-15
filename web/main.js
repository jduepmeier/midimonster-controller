function handleError(err) {
	if (err !== undefined) {
		log = document.querySelector("#log");
		log.value += "\nERROR: " + myJson["Error"];
	}
}

async function writeConfig() {
  config = document.querySelector("#config")
  fetch('/api/write', {
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
  fetch('/api/reload', {
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
	const response = await fetch('/api/config');
	const myJson = await response.json();
	textArea = document.querySelector("#config");
	textArea.value = myJson.Content;
}

async function getStatus() {
	const response = await fetch('/api/status');
	const myJson = await response.json();
	statusBlock = document.querySelector("#status");
    if (myJson.Code === 0) {
        statusBlock.classList.add("green")
        statusBlock.classList.remove("red")
    } else {
        statusBlock.classList.add("red")
        statusBlock.classList.remove("green")
    }
	statusBlock.innerText = myJson.Text;
}

async function init() {
	log = document.querySelector("#log");
	config = document.querySelector("#config");
    getStatus();
    setInterval(() => {
        getStatus();
    }, 5000)
	get();
}

init();