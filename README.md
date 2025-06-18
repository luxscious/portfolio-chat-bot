# 🤖 Gabriella’s Resume Chatbot

A personalized, first-person chatbot that answers questions based on my resume, work experience, projects, and personality. Powered by React, Go, OpenAI, and a structured `resume.json` file.

---

## 🧠 What It Does

- Uses my structured resume data and personality traits
- Embeds that data for retrieval via semantic search (RAG-ready)
- Builds a prompt in my tone using `personaContext`
- Sends user input + resume context to OpenAI’s GPT model
- Stores chat history per user with in-memory threads or MongoDB
- Returns chatbot-style answers like they're directly from me

---

## 🧱 Architecture

```
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
│   with OpenAI API            │
│ - Stores/retrieves chat      │
│   history from MongoDB       │
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
- **CORS Configurable**: Frontend origin is loaded from `.env`
- **MongoDB Storage**: Optional persistent chat history storage
- **Modular & Clean**: Frontend and backend separated
- **Deployable**: Vercel (frontend), Railway/Render (backend)

---

## 📁 File Structure

```
resume-chatbot/
├── client/            # React frontend
│   └── src/
├── server/            # Go backend
│   ├── main.go
│   ├── routes.go
│   ├── resume.go
│   ├── embedding.go
│   ├── chat_history.go
│   ├── call_openai.go
│   ├── config/
│   │   └── env.go
│   ├── db/
│   │   └── mongo.go
│   └── resume.json
├── .env.example       # OpenAI, Mongo, and frontend config
├── README.md
```

---

## 🚀 Running the Project

1. **Clone the repo**

```bash
git clone https://github.com/yourusername/resume-chatbot.git
```

2. **Set up your `.env` files**

```bash
cp .env.example .env
```

Make sure `.env` contains:

```env
OPENAI_API_KEY=your-key-here
FRONTEND_ORIGIN=http://localhost:5173
OPENAI_API_URL=https://api.openai.com/v1/chat/completions
OPENAI_EMBEDDING_URL=https://api.openai.com/v1/embeddings
MONGO_URI=mongodb://localhost:27017
MONGO_DB=resumeChatbot
MONGO_COLLECTION=messages
```

3. **Install frontend dependencies**

```bash
cd client && npm install && npm run dev
```

4. **Run backend (Go)**

```bash
cd ../server && go run .
```

5. **Chat live!** Frontend sends messages to Go backend which returns LLM responses.

---

## 📦 Future Enhancements

- Vector search with cosine similarity
- Streaming GPT output to UI
- Admin dashboard to manage/edit resume.json
- Authenticated sessions for persistent chat threads

### ✅ Resume File Safety

- `resume.json` is only read once on startup — no hot reloading or public exposure.
- No sensitive keys or tokens are stored in the JSON structure.
- Embedding and prompt construction happen securely on the server.
- CORS settings are loaded from `.env` for flexibility.
- MongoDB used to persist chat history per user.
- Future: consider access control or rate-limiting on the `/chat` endpoint.
