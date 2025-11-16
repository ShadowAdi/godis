const net = require("net");
const readline = require("readline");

// Convert user input (like: set user adi) to RESP
function encodeRESP(input) {
  const parts = input.trim().split(" ");

  let resp = `*${parts.length}\r\n`;

  for (const part of parts) {
    resp += `$${part.length}\r\n${part}\r\n`;
  }

  return resp;
}

// Create TCP client
const client = new net.Socket();

client.connect(6380, "127.0.0.1", () => {
  console.log("Connected to Go server");
});

// When server sends back something
client.on("data", (data) => {
  console.log("Server:", data.toString());
});

// On close
client.on("close", () => {
  console.log("Connection closed");
});

// Read user input
const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

// On user entering text
rl.on("line", (input) => {
  const resp = encodeRESP(input);
  client.write(resp);
});
