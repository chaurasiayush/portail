# 🌐 Portail

> A blazing-fast, YAML-configurable port forwarder written in Go.  
> Supports TCP/UDP forwarding, TLS passthrough, and hot-reloading on config changes.

---

## 🚀 Features

- 🔁 Forward **TCP** and **UDP** ports to remote destinations
- 🔐 Optional **TLS passthrough** with skip verify support
- 🔧 Configurable via a single `config.yaml` file
- ♻️ **Hot reload** on config file changes
- 🛠️ Minimal dependencies, easy to run as a binary or Docker container

---

## 📦 Installation

### Option 1: Build from source

```bash
git clone https://github.com/your-username/portail.git
cd portail
go build -o portail
```

### Option 2: Run with Docker

Coming soon! (Let me know if you want the Dockerfile now.)

---

## ⚙️ Usage

```bash
./portail --config=config.yaml
```

> Or if using `go run`:
>
> ```bash
> go run . --config=config.yaml
> ```

---

## 🧾 Sample `config.yaml`

```yaml
forwards:
  - protocol: tcp
    listen: "127.0.0.1:8080"
    forward: "example.com:80"

  - protocol: udp
    listen: "127.0.0.1:5353"
    forward: "8.8.8.8:53"

  - protocol: tcp
    listen: "127.0.0.1:8443"
    forward: "secure.example.com:443"
    tls:
      enabled: true
      skip_verify: true
```

---

## 🔄 Hot Reload

Portail automatically watches your `config.yaml` file and applies changes on the fly — no restart needed.

---

## 📂 Project Structure

```
portail/
├── main.go         # Entry point and CLI flags
├── forwarder.go    # TCP/UDP forwarding logic
├── config.go       # YAML parsing and config structs
├── config.yaml     # Sample config
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