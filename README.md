# 🤖 Gabriella’s Resume Chatbot

A personalized, first-person chatbot that answers questions based on my resume, work experience, projects, and personality. Powered by React, Go, OpenAI, and a structured `resume.json` file.

---

## 🧠 What It Does

- Uses my structured resume data and personality traits
- Embeds that data for retrieval via semantic search (RAG-ready)
- Builds a prompt in my tone using `personaContext`
- Sends user input + resume context to OpenAI’s GPT model
- Returns chatbot-style answers like they're directly from me

---

## 🛡️ Architecture

```
🔐 .env
  └️ Configures secrets (OpenAI key, frontend origin, etc.)

┌────────────────────┐
│ 1. React Frontend  │
│ (Vite-based app)   │
└────────┬───────────┘
         │ POST /chat
         ▼
┌──────────────────────────────┐
│ 2. Go Backend                │
│ (API Server)                 │
│                              │
│ - Loads resume.json          │
│ - Builds prompt              │
│ - Generates embeddings       │
│ - Handles user history       │
│ - Routes via Chi             │
└────────┬─────────────────────┘
         │ API Request
         ▼
┌────────────────────┐
│ 3. OpenAI API      │
│ (GPT-4 / GPT-3.5)  │
└────────┬───────────┘
         │ Completion
         ▼
┌────────────────────┐
│ Back to Go Server  │
└────────┬───────────┘
         │ JSON Response
         ▼
┌────────────────────┐
│  React Frontend    │
│  displays message  │
└────────────────────┘
```

---

## 🔍 Key Features

- **Dynamic Prompting**: GPT responses use my voice, tone, and background
- **Resume-Driven**: Pulls structured data from `resume.json`
- **Embeddings**: Generates vector embeddings from resume content
- **Chat Memory**: User-specific thread history with support for future sessions
- **CORS Configurable**: Frontend origin is loaded from `.env`
- **Modular & Clean**: Frontend and backend separated into Go packages
- **Deployable**: Vercel (frontend), Railway/Render (backend)

---

## 📁 File Structure

```
resume-chatbot/
├── client/              # React frontend
│   └── src/
├── server/              # Go backend
│   ├── main.go
│   ├── routes.go
│   ├── open_ai.go       # Handles OpenAI chat completions
│   ├── embedding.go     # Handles resume embeddings
│   ├── resume.go        # Resume structs + loading logic
│   ├── chat_history.go  # Per-user session history
│   └── config/
│       └── env.go       # Centralized env var access
├── .env.example         # Environment variable sample
├── README.md
```

---

## 🚀 Running the Project

1. **Clone the repo**

```bash
git clone https://github.com/yourusername/resume-chatbot.git
```

2. **Set up your `.env` file**

```bash
cp .env.example .env
```

Edit `.env` with:

```env
OPENAI_API_KEY=your-key-here
OPENAI_API_URL=https://api.openai.com/v1/chat/completions
OPENAI_EMBEDDING_URL=https://api.openai.com/v1/embeddings
FRONTEND_ORIGIN=http://localhost:5173
REACT_APP_API_URL=http://localhost:8080/chat
PORT=8080
```

3. **Install frontend dependencies**

```bash
cd client && npm install && npm run dev
```

4. **Run backend (Go)**

```bash
cd ../server && go run main.go
```

5. **Chat live!** Frontend sends messages to Go backend which returns LLM responses.

---

## 📦 Future Enhancements

- Vector search with cosine similarity
- Streaming GPT output to UI
- Admin dashboard to manage/edit resume.json

### ✅ Resume File Safety

- `resume.json` is only read once on startup — no hot reloading or public exposure.
- No sensitive keys or tokens are stored in the JSON structure.
- Embedding and prompt construction happen securely on the server.
- CORS settings are loaded from `.env` for flexibility.
- Future: consider access control or rate-limiting on the `/chat` endpoint.
