# 🤖 Gabriella’s Resume Chatbot

A personalized, first-person chatbot that answers questions based on my resume, work experience, skills, hobbies, and overall personality. Built with React + Vite, a Go backend, OpenAI, Neo4j for structured querying, and Ollama for intent parsing.

---

## 🧠 What It Does

- Embeds my structured resume data for retrieval (no raw text scraping)
- Uses Ollama (LLaMA3) to identify **target node types** and **filters**
- Runs Cypher queries on Neo4j to build **context**
- Builds prompts in my voice using `personaContext`
- Sends to OpenAI (GPT-3.5-Turbo)
- Returns chat responses that feel personal and natural
- Stores conversation history per user via MongoDB

---

## 🧱 Architecture

```
┌────────────────────────────┐
│  React Frontend (Vite)     │
└────────────┬───────────────┘
             │
             ▼
┌────────────────────────────┐
│   Go Backend (API Server)  │
│                            │
│ 1. Accepts input           │
│ 2. Plans query via Ollama  │
│ 3. Retrieves data via Neo4j│
│ 4. Builds context prompt   │
│ 5. Calls OpenAI            │
│ 6. Stores message in Mongo │
└────────────┬───────────────┘
             │
             ▼
     ┌────────────────┐
     │   Ollama API   │  (DigitalOcean droplet)
     └──────┬─────────┘
            │
            ▼
     ┌────────────────┐
     │ Neo4j Aura DB  │  (Managed Cloud)
     └──────┬─────────┘
            │
            ▼
     ┌────────────────┐
     │  OpenAI GPT    │
     └────────────────┘
             ▲
             │
     ┌────────────────┐
     │   MongoDB Atlas│  (Cloud)
     └────────────────┘
```

- **Caddy Reverse Proxy** (on the host) provides HTTPS and routes traffic to the backend container.

---

## 📁 File Structure

```
portfolio-chat-bot/
├── client/                  # React + Vite frontend
│   ├── src/
│   │   ├── components/      # UI + chat components
│   │   ├── hooks/           # Custom React hooks
│   │   ├── lib/             # Utility functions
│   │   ├── pages/           # ChatPage, App.tsx, etc.
│   │   └── styles/          # Tailwind + custom CSS
│   ├── public/              # Static assets
│   ├── .env.example         # Frontend env example
│   ├── vite.config.ts       # Vite config
│   └── ...                  # Other frontend files
├── server/                  # Go backend
│   ├── config/              # .env loading and config
│   ├── db/                  # Neo4j + Mongo logic
│   ├── openai/              # GPT + context builder logic
│   ├── ollama/              # Intent planning via LLaMA3
│   ├── resume/              # resume.json and chunking
│   ├── main.go              # Entry point
│   ├── routes.go            # Route definitions
│   ├── Dockerfile           # Backend Dockerfile
│   └── .env.example         # Backend env example
├── docker-compose.yml       # Orchestrates backend + Ollama
├── .github/
│   └── workflows/           # GitHub Actions CI/CD
│       ├── deploy-backend.yml
│       └── deploy.yml
└── README.md
```

**Notes:**

- Neo4j is now managed in the cloud (Aura), not as a local container.
- Ollama runs on a dedicated DigitalOcean droplet, not locally.
- Caddy runs on the host to provide HTTPS and reverse proxy.
- MongoDB is typically managed via MongoDB Atlas (cloud), but can be local for dev.

---

## ⚙️ Running the Project

### 1. Clone the repo

```bash
git clone https://github.com/luxscious/portfolio-chat-bot.git
cd portfolio-chat-bot
```

### 2. Set up env files

```bash
cp client/.env.example client/.env
cp server/.env.example server/.env
```

Then update:

```env
# OpenAI
OPENAI_API_KEY=sk-...
OPENAI_API_URL=https://api.openai.com/v1/chat/completions
OPENAI_EMBEDDING_URL=https://api.openai.com/v1/embeddings
OPENAI_THREAD_URL=https://api.openai.com/v1/threads
OPENAI_ASSISTANT_ID=

# Ollama
OLLAMA_URI=http://localhost:11434

# Server config
PORT=8080
FRONTEND_ORIGIN=http://localhost:5173

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=resumeChatbot
MONGO_COLLECTION=messages

# Neo4j
NEO4J_URI=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASS=password
```

### 3. Install frontend dependencies

```bash
cd client && npm install && npm run dev
```

### 4. Run backend (Go)

```bash
cd ../server && go run .
```

---

## 🐞 Known Bugs

- If the backend disconnects temporarily, the frontend doesn’t retry or recover well
- Need a better indicator for GPT loading vs. server unavailability
- Intent parsing is not 100% fool proof.

---

## 🌱 Future Ideas

- Look into costs + latency advantages of combining intent parsing and response into one LLM
- Caching to save on DB lookup

---

## 🐳 Docker Setup

This project is also configured to run all services via Docker Compose:

- Neo4j (graph database)
- Ollama (local LLM server)
- Go backend API

### 🛠️ Prerequisites

- Docker and Docker Compose.

Sure—here’s that as a **clean Markdown canvas snippet** you can drop right into your README:

---

## 📥 Pull Ollama Model

Before running the containers, **make sure you have downloaded the model you want Ollama to serve.**

You need to **enter the Ollama container shell** to pull the model.

For example:

```bash
docker compose run --rm ollama bash
```

Then inside the container shell:

```bash
ollama pull llama3
```

### 📝 Create `.env`

- Put all variables in project root `.env` file

### ⚡ Load environment variables (Need to do this for Neo4J Auth)

```bash
export $(cat .env | xargs)
```

### 🐳 Build and Run

```bash
docker compose up --build
```

This starts Neo4j, Ollama, and your backend.

### 🛑 Stopping

```bash
docker compose down
```

### ✅ Notes

- The frontend runs separately with Vite.

---

## 🚀 CI/CD

This project uses GitHub Actions for continuous integration and deployment:

### Deployment

- **Backend Deployment:**  
  On every push to `main` that affects the backend (`server/**`), `docker-compose.yml`, or the backend deploy workflow, the backend is automatically deployed to the production server via SSH.  
  The workflow:

  1. Checks out the latest code.
  2. Sets up SSH keys and known hosts.
  3. Connects to the server, pulls the latest code, and runs `docker compose up --build -d`.

- **Frontend Deployment:**  
  On every push to `main` that affects the frontend (`client/**`), the frontend is built and deployed to GitHub Pages.  
  The workflow:
  1. Checks out the latest code.
  2. Sets up Node.js and environment variables.
  3. Installs dependencies and builds the Vite project.
  4. Publishes the build output to GitHub Pages.

You can find the workflow files in `.github/workflows/`:

- `deploy-backend.yml` — Deploys the backend to the server.
- `deploy.yml` — Deploys the frontend to GitHub Pages.

## 🌐 Deployment Notes

This project is deployed on a DigitalOcean server to ensure reliable performance and scalability.

Why DigitalOcean?

Resource Requirements for Llama 3:Running the Ollama Llama 3 model requires substantial memory and CPU resources, which are difficult to allocate reliably on local machines or small VMs.

Production Stability:DigitalOcean provides predictable resource availability and uptime for hosting the backend API, Ollama server, and Caddy reverse proxy.

Simplicity:The deployment uses Docker Compose to manage all services consistently across environments.

How it Works

Ollama ServiceHosts the Llama 3 model server (ollama serve), exposing the API on port 11434.

Backend ServiceA Go API that proxies requests to Ollama and handles application logic.

Caddy Reverse ProxyProvides automatic HTTPS with Let’s Encrypt certificates for secure public access.

✅ Tip:If you want to run this stack locally, you’ll need a machine with enough RAM (8–16 GB recommended) to run the Llama 3 model without OOM errors.
