# 🚀 Redis-Go Rate Limiting API with Nginx Load Balancing

This project demonstrates a **high-performance API** built using **Golang (Fiber)**, **Redis** for rate limiting, and **Nginx** for reverse proxy and load balancing. The app simulates a scalable, production-ready backend architecture — ideal for real-world traffic management and deployment scenarios.

---

## 📦 Tech Stack

- **Golang** with [Fiber](https://github.com/gofiber/fiber) – ultra-fast web framework
- **Redis** – in-memory store for rate-limiting control
- **Nginx** – reverse proxy and load balancer
- **Docker & Docker Compose** – containerized multi-service setup
- **Logger Middleware** – logs all incoming requests and blocked IPs

---

## ⚙️ Features

✅ Redis-based IP-level rate limiting (5 requests/minute)  
✅ Blocked IP logging to `blocked_ips.log`  
✅ Nginx load balancing across multiple Fiber app instances  
✅ Dedicated `/health` route bypassing rate limiting  
✅ Real-world deployment simulation using Docker Compose  

---

## 🗂️ Project Structure

```bash
.
├── nginx/
│   └── default.conf        # Nginx config (load balancing + proxy)
├── tmp/                    # Temporary dev files
├── main.go                 # Golang Fiber app
├── blocked_ips.log         # IPs blocked due to rate limiting
├── Dockerfile              # Container image for Fiber app
├── docker-compose.yml      # Multi-service orchestration
├── go.mod / go.sum         # Go module dependencies
├── .air.toml               # Optional: Air (live reload config)
└── .gitignore              # Git tracking exclusions



---

## 🔁 How Rate Limiting Works

Each client IP is tracked via Redis using the pattern `rate<IP>`. If the request count exceeds **5 requests per minute**, the IP is blocked and logged to `blocked_ips.log`. A graceful `429` is returned.

Health route (`/health`) skips rate limiting completely.

---

## 🧪 Load Balancing in Action

We use **Nginx** to distribute requests to **3 Fiber app containers**. Each `/test` call simulates a heavy computation task and returns the container hostname.

### Example curl:

```bash
curl http://localhost/test
