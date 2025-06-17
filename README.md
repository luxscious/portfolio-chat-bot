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
│ - Generates full resume      │
│   embedding with OpenAI API  │
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
- **Embeddings**: Generates a full resume vector embedding
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
│   ├── resume.go
│   ├── embedding.go
│   └── resume.json
├── .env.example       # OpenAI API key and server config
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
- Future: consider access control or rate-limiting on the `/chat` endpoint.
