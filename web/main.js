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
	let textArea = document.querySelector("#config");
	textArea.value = myJson.Content;
}

async function getStatus() {
	const response = await fetch('/api/status');
	const myJson = await response.json();
	let statusBlock = document.querySelector("#status");
    if (myJson.Code === 0) {
        statusBlock.classList.add("green")
        statusBlock.classList.remove("red")
    } else {
        statusBlock.classList.add("red")
        statusBlock.classList.remove("green")
    }
	statusBlock.innerText = myJson.Text;
}

async function getLogs() {
	const response = await fetch(`/api/logs?oldest=${logsOldest}`);
	const myJson = await response.json();
	let logsBlock = document.querySelector("#logs");
  if (myJson.Logs.length > 0) {
    logsBlock.value += myJson.Logs.join("\n")
    logsBlock.value += "\n"
  }
  logsOldest = myJson.Newest
  logsBlock.scrollTop = logsBlock.scrollHeight
}

async function init() {
  getStatus();
  logsOldest = 0
  setInterval(() => {
      getStatus();
      getLogs();
  }, 5000)
	get();
}

init();