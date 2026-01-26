<div align="center">

<img src="https://raw.githubusercontent.com/zenkiet/window-service-watcher/refs/heads/main/frontend/src/icons/logo.png" width="250" height="250" alt="App Logo">

<h1>Window Service Watcher</h1>

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00?style=flat-square&logo=svelte)](https://svelte.dev/)
[![Wails](https://img.shields.io/badge/Wails-2.9-E3202E?style=flat-square&logo=wails)](https://wails.io/)
[![License](https://img.shields.io/github/license/zenkiet/window-service-watcher?style=flat-square)](LICENSE)

**A lightweight, modern desktop application for real-time monitoring and management of Windows services.**


</div>

---

## ✨ Overview

Provide a seamless dashboard for developers and system admins to monitor status, track resources (CPU/RAM), and control critical services. Default configuration supports monitoring POS services, but can be adapted for any Windows service.

### Key Features

| Feature | Description |
| :--- | :--- |
| **⚡ Real-time Monitoring** | Live status updates (Running, Stopped, Paused) with sub-second latency. |
| **📊 Resource Tracking** | Visual graphs for CPU usage (%) and Memory consumption (MB). |
| **🎮 Service Control** | Start, Stop, and Restart services directly from the UI. |
| **📦 Smart Installation** | Install and update service executables/files with one click. |
| **⚙️ Configurable** | Simple YAML configuration for managing multiple services. |
| **🔔 System Tray** | (Coming Soon) Background operation and status notifications. |

---

## 📸 Screenshots

<div align="center">

| **Real-time Monitoring** | **Configuration & Settings** |
|:---:|:---:|
| <img src="https://raw.githubusercontent.com/zenkiet/window-service-watcher/refs/heads/main/assets/screenshot/1.png" alt="Main Dashboard" width="100%"> | <img src="https://raw.githubusercontent.com/zenkiet/window-service-watcher/refs/heads/main/assets/screenshot/2.png" alt="Settings" width="100%"> |

| **Service Installer** | **OTA Updates** |
|:---:|:---:|
| <img src="https://raw.githubusercontent.com/zenkiet/window-service-watcher/refs/heads/main/assets/screenshot/3.png" alt="Install Service" width="100%"> | <img src="https://raw.githubusercontent.com/zenkiet/window-service-watcher/refs/heads/main/assets/screenshot/4.png" alt="Update OTA" width="100%"> |

</div>

---

## 🏗️ Architecture

Built on the **Clean Architecture** principle, ensuring separation of concerns and maintainability.

- **Frontend**: [Svelte 5](https://svelte.dev) + [TailwindCSS](https://tailwindcss.com) for a high-performance, reactive UI.
- **Backend**: [Go](https://go.dev) + [Wails](https://wails.io) for native Windows API interaction.
- **Bridge**: Async communication via Wails runtime events (`services-update`).

```bash
window-service-watcher/
├── internal/
│   ├── app/        # App lifecycle & Watcher logic
│   ├── config/     # Config loader (YAML)
│   ├── domain/     # Interfaces & DTOs
│   └── service/    # Windows SCM implementation
├── frontend/src/   # Svelte 5 UI Components
├── main.go         # Entry point
└── wails.json      # Wails config
```

---

## ⚙️ Installation and Setup

### Prerequisites
*   **Go** (v1.25+)
*   **Node.js** (LTS)
*   **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### 🛠️ Steps

1.  **Clone & Install Dependencies**
    ```bash
    git clone https://github.com/zenkiet/window-service-watcher.git
    cd window-service-watcher

    # Install backend & frontend deps
    make deps
    cd frontend && pnpm install
    ```

2.  **Run in Development Mode**
    ```bash
    wails dev
    ```

3.  **Build for Production**
    ```bash
    wails build
    # Output: ./build/bin/window-service-watcher.exe
    ```

---

## 🚀 Usage Guide

### Configuration (`config.yaml`)

The app auto-generates this file on first run.

```yaml
services:
  - id: "550e8400-e29b-41d4-a716-446655440000" # Unique ID
    name: "Payment Gateway Service"           # Display Name
    serviceName: "PaymentSvc"                 # Actual Windows Service Name
    path: "C:\\Services\\Payment\\app.exe"    # Executable Path
    installable: true                         # Allow file updates
```

### Dashboard Controls
*   **Status Cards**:  (🟢 Running, 🔴 Stopped, 🟡 Pending).
*   **Play/Stop**: Toggle service state.
*   **Metrics**: Hover over cards to see detailed CPU/RAM history.

---

## 🔌 API / Key Functions

### Backend (Go)

| Function | Description |
| :--- | :--- |
| `Startup(ctx)` | Initializes Service Manager connection. |
| `Start/StopService(id)` | Controls specific service state. |
| `InstallService(id, files)` | Updates binary files and restarts service. |
| `GetConfig()` | Returns the current loaded configuration. |

### Frontend Events

| Event | Payload | Description |
| :--- | :--- | :--- |
| `services-update` | `map[string]Status` | Real-time status stream (1Hz). |

---

## 🐞 Troubleshooting

<details>
<summary><strong>❌ Service Not Found / Error Fetching State</strong></summary>

*   **Cause**: `serviceName` in `config.yaml` doesn't match the actual Windows Service name.
*   **Fix**: Check `services.msc` for the exact Service Name (not Display Name) and update config.
</details>

<details>
<summary><strong>🚫 Permission Denied</strong></summary>

*   **Cause**: Controlling Windows services requires Administrator privileges.
*   **Fix**: Run the application or terminal as **Administrator**.
</details>

<details>
<summary><strong>📉 Metrics Not Showing</strong></summary>

*   **Cause**: The service might be running but the process ID (PID) cannot be retrieved, or monitoring is blocked.
*   **Fix**: Ensure the service is actually running and you have permissions to query its process information.
</details>

---

## 🤝 Contributing

Contributions are welcome! Please open issues or pull requests for bug fixes, features, or improvements.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

### **If this project helped you, please consider supporting its development**
```
Star the repository ⭐
Share it with your network 📢
```

<sub>Made with ❤️ by <a href="https://github.com/zenkiet">ZenKiet</a></sub>

</div>