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
     │   Ollama API   │
     └──────┬─────────┘
            │
            ▼
     ┌────────────────┐
     │   Neo4j DB     │
     └──────┬─────────┘
            │
            ▼
     ┌────────────────┐
     │  OpenAI GPT    │
     └────────────────┘
             ▼
┌────────────────────────────┐
│  React: renders reply      │
└────────────────────────────┘
└────────────────────────────┘
```

---

## 🔍 Key Features

- **Graph-Powered Prompting**: Semantic plans via Ollama → structured context from Neo4j
- **Persona-Aware**: All responses are written in my own tone with first-person voice
- **Mongo-Backed Memory**: Conversation stored per-user for persistence
- **Small Talk Handling**: Recognizes vague input like “hi!” and skips AI calls

---

## 📁 File Structure

```
portfolio-chat-bot/
├── client/              # React + Vite frontend
│   ├── src/
│   │   ├── components/  # UI + chat components
│   │   ├── hooks/       # TypingEffect.ts etc.
│   │   ├── lib/         # Utility functions
│   │   ├── pages/       # ChatPage, App.tsx
│   │   └── styles/      # Tailwind + custom CSS
├── server/              # Go backend
│   ├── config/          # .env loading
│   ├── db/              # Neo4j + Mongo logic
│   ├── openai/          # GPT + SmartQuery logic
│   ├── ollama/          # Intent planning via LLaMA3
│   ├── resume/          # resume.json and chunking
│   ├── main.go          # Route binding
│   └── routes.go        # Route definitions
```

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
