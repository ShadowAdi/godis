# 🚀 GoDis — A Redis-Inspired In-Memory Data Store in Golang

GoDis is a lightweight, high-performance, Redis-style in-memory key–value database built **entirely from scratch** in Golang.  
It supports core Redis-like commands, multi-client TCP connections, thread-safe operations, and a modular architecture that can be easily extended.

This project was created to understand how systems like Redis work internally — including networking, concurrency, command parsing, and low-level storage design.

---

## 🔥 Features

### ✔ Core Commands Implemented
- `SET key value` – Store a string value  
- `GET key` – Retrieve a string value  
- `HSET hash field value` – Set a field in a hash  
- `HGET hash field` – Retrieve a field from a hash  
- `HGETALL hash` – Retrieve entire hash object  

### ✔ Engine Capabilities
- **Thread-safe in-memory data store** using `sync.RWMutex`
- **Multi-client support** via TCP
- **Custom command parser** (similar to RESP)
- **Lightweight, minimal, and fast Go codebase**
- Clean architecture → Easy to add more commands

---

## 🏗 Architecture Overview

The internal workflow of GoDis:

Client → TCP Server → Command Parser → Command Router → In-Memory Store → Response


### Components
- **TCP Listener** – Accepts incoming connections  
- **Parser** – Reads raw input & breaks commands into tokens  
- **Command Handlers** – Implements SET, GET, HSET, HGET, HGETALL  
- **In-Memory Store**  
  - Strings: `map[string]string`
  - Hashes: `map[string]map[string]string`
- **RWMutex** – Ensures thread-safe operations

---

## 📦 Installation & Running the Server

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/godis.git
cd godis
```

### 2. Run the Go server
```bash
go run main.go
```

The server will start listening on:
```bash
localhost:6379
```

### 💡 Usage Example
Open any terminal and use nc or any TCP client.
SET name Aditya
GET name
