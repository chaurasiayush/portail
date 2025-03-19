# Portail 🌀

**Portail** is a simple, production-ready, and hot-reloadable TCP/UDP port forwarder written in Go. It uses a declarative YAML configuration to define forwarding rules and supports live reloading using `fsnotify`.

---

## 🚀 Features

- 🔁 Port forwarding for both **TCP** and **UDP**
- 📁 Simple YAML-based configuration
- ♻️ **Hot reload** support using `fsnotify`
- 🛠️ Graceful shutdown via OS signals
- 🧱 Production-ready structure

---

## 📦 Project Structure

```
portail/
├── cmd/
│   └── portail/
│       └── main.go           # Main entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Loads YAML config
│   └── forwarder/
│       └── forwarder.go      # TCP and UDP forwarders
├── config.yaml              # Sample config file
├── Dockerfile               # Docker image definition
├── go.mod
├── go.sum
└── README.md
```

---

## 🧪 Example `config.yaml`

```yaml
forwards:
  - protocol: tcp
    listen: ":8080"
    target: "localhost:3000"

  - protocol: udp
    listen: ":5353"
    target: "localhost:5354"
```

---

## 🛠️ Usage

### 1. Build the Binary

```bash
go build -o portail ./cmd/portail
```

### 2. Run with Config

```bash
./portail --config=config.yaml
```

### 3. Modify Config
Just edit and save `config.yaml`. `portail` will auto-detect changes and reload the file.

---

## 🧠 How It Works

- `fsnotify.NewWatcher()` monitors the config file
- On file **Write** events, it re-parses and reloads the config
- Graceful signal handling with `SIGINT`/`SIGTERM`

---

## 🐳 Docker

### Build Docker Image

```bash
docker build -t portail .
```

### Run with Volume-Mounted Config

```bash
docker run --rm \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -p 8080:8080 \
  -p 5353:5353/udp \
  portail --config=/app/config.yaml
```
---
## 👨‍💻 Contributing

Got ideas or feature requests? PRs and issues welcome!

---
## 📄 License

MIT License. See `LICENSE` file for details.

---

## ✨ Credits

Built with ❤️ in Go.

---