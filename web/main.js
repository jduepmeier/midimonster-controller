import { CodeJar } from './codejar.js';
import { withLineNumbers } from './linenumbers.js';
import hljs from './hljs/core.js';
import ini from './hljs/ini.min.js';

class App {
  constructor() {
    this.logsOldest = 0
    this.messagePromise = undefined;
    this.dom = {
      config: document.querySelector("#config"),
      logs: document.querySelector("#logs"),
      status: document.querySelector("#status"),
      autoScroll: document.querySelector("#autoscroll"),
      toast: document.querySelector("#toast"),
    }
    document.querySelector('#write').addEventListener("click", () => {
      this.writeConfig();
    });
    document.querySelector('#restart').addEventListener("click", () => {
      this.restart();
    });
    document.querySelector('#get').addEventListener("click", () => {
      this.get();
    });
  }

  async init() {
    hljs.registerLanguage('ini', ini)
    this.code = CodeJar(this.dom.config, withLineNumbers((editor) => {
      const code = editor.textContent;
      editor.innerHTML = hljs.highlight(code, { language: 'ini' }).value
    }, {
      color: '#222',
    }))

    this.connectWebSocket(this);

    this.get();
  }

  handleError(err, statusMessage) {
    if (err !== undefined) {
      this.dom.logs.value += "ERROR: " + myJson["Error"] + "\n";
    } else {
      this.showMessage(statusMessage);
    }
  }

  async writeConfig() {
    fetch('/api/write', {
      method: 'POST',
      body: JSON.stringify({
        "Content": this.code.toString(),
      }),
      headers: {
        'Content-Type': 'application/json'
      }
    }).then((response) => {
      response.json().then((myJson) => {
        this.handleError(myJson["Error"], "config written");
      });
    }).catch((err) => {
      this.handleError(err);
    });
  }

  async restart() {
    fetch('/api/reload', {
      method: 'POST',
      body: JSON.stringify({}),
      headers: {
        'Content-Type': 'application/json',
      }
    }).then((response) => {
      response.json().then((myJson) => {
        this.handleError(myJson["Error"], "restarted");
      });
    }).catch((err) => {
      this.handleError(err);
    });
  }

  async connectWebSocket() {
    let host = location.host;
    let protocol = "ws";
    if (location.protocol == "https") {
      protocol = "wss";
    }
    this.ws = new WebSocket(`${protocol}://${host}/api/ws`);
    this.ws.onopen = () => {
      console.log("websocket open");
      this.ws.send(JSON.stringify({
        "command": "getLogs"
      }))
      this.ws.send(JSON.stringify({
        "command": "getStatus"
      }))
    };
    this.ws.onclose = () => {
      console.log("websocket closed. Try to reopen.");
      this.ws = null;
      setTimeout(() => {
        this.connectWebSocket()
      }, 5000);
    };
    this.ws.onmessage = (evt) => {
      console.log("received message", evt);

      let data = JSON.parse(evt.data);
      switch (data.type) {
        case "logs":
          console.log("got logs message");
          this.appendLogs(data.lines)
          break;
        case "status":
          console.log("got status message")
          this.setStatus(data.status)
          break;
      }
    }
    this.ws.onerror = (evt) => {
      console.log("got error from websocket", evt);
    }
  }

  async get() {
    const response = await fetch('/api/config');
    const myJson = await response.json();
    this.code.updateCode(myJson.Content);
    this.showMessage("config fetched");
  }

  async setStatus(status) {
    let statusBlock = this.dom.status
    if (status.code === 0) {
      statusBlock.classList.add("green")
      statusBlock.classList.remove("red")
    } else {
      statusBlock.classList.add("red")
      statusBlock.classList.remove("green")
    }
    statusBlock.innerText = status.text;
  }

  async showMessage(msg) {
    Toastify({
      text: msg,
      duration: 3000,
      newWindow: true,
      close: true,
      gravity: "bottom",
      position: "right",
      onClick: () => { },
    }).showToast();
  }

  async appendLogs(logs) {
    let logsBlock = this.dom.logs
    logsBlock.value += logs.join("\n");
    logsBlock.value += "\n";
    if (this.dom.autoScroll.checked) {
      logsBlock.scrollTop = logsBlock.scrollHeight;
    }
  }
}

let app = new App();
app.init()
export default app;
