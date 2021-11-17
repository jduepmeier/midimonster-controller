import {CodeJar} from './codejar.js';
import {withLineNumbers} from './linenumbers.js';
import hljs from './hljs/core.js';
import ini from './hljs/ini.min.js'

class App {
  constructor() {
    this.logsOldest = 0
    this.dom = {
      config: document.querySelector("#config"),
      logs: document.querySelector("#logs"),
      status: document.querySelector("#status"),
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

    this.getStatus();
    setInterval(() => {
        this.getStatus();
        this.getLogs();
    }, 5000)
    this.get();
  }

  handleError(err) {
    if (err !== undefined) {
      this.dom.logs.value += "ERROR: " + myJson["Error"] + "\n";
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
      this.handleError(myJson["Error"]);
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
      this.handleError(myJson["Error"]);
    });
    }).catch((err) => {
      this.handleError(err);
    });
  }

  async get() {
    const response = await fetch('/api/config');
    const myJson = await response.json();
    this.code.updateCode(myJson.Content);
  }

  async getStatus() {
    const response = await fetch('/api/status');
    const myJson = await response.json();
    let statusBlock = this.dom.status
      if (myJson.Code === 0) {
          statusBlock.classList.add("green")
          statusBlock.classList.remove("red")
      } else {
          statusBlock.classList.add("red")
          statusBlock.classList.remove("green")
      }
    statusBlock.innerText = myJson.Text;
  }

  async getLogs() {
    const response = await fetch(`/api/logs?oldest=${this.logsOldest}`);
    const myJson = await response.json();
    let logsBlock = this.dom.logs
    if (myJson.Logs.length > 0) {
      logsBlock.value += myJson.Logs.join("\n")
      logsBlock.value += "\n"
    }
    this.logsOldest = myJson.Newest
    logsBlock.scrollTop = logsBlock.scrollHeight
  }
}

let app = new App();
app.init()
export default app;