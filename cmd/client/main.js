const net = require("net");
const readline = require("readline");

const client = new net.Socket();

client.connect(6380, "127.0.0.1", () => {
  console.log("Connected to Go server");
});

client.on("data", (data) => {
  console.log("Server:", data.toString());
});

client.on("close", () => {
  console.log("Connection closed");
});

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

rl.on("line", (input) => {
  client.write(input + "\n");
});
