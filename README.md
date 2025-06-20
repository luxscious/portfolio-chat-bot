Here’s your updated `README.md` to reflect the current state of your chatbot, including the removal of project-related matching logic and clarifying the use of semantic search for chunked data only:

---

# 🤖 Gabriella’s Resume Chatbot

A personalized, first-person chatbot that answers questions based on my resume, work experience, projects, and personality. Powered by React, Go, OpenAI, and a structured `resume.json` file.

---

## 🧠 What It Does

- Uses my structured resume data and personality traits
- Embeds that data for retrieval via semantic search (RAG-ready)
- Builds a prompt in my tone using `personaContext`
- Sends user input + relevant context to OpenAI’s GPT model
- Returns responses that sound like me
- Stores conversation history per user (via MongoDB)

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
│ - Splits into semantic chunks│
│ - Generates embeddings       │
│ - Ranks top 3 relevant chunks│
│ - Builds prompt w/ persona   │
│ - Sends to OpenAI API        │
│ - Stores chat history        │
└────────┬─────────────────────┘
         │
         ▼
┌────────────────────┐
│ 3. OpenAI API      │
│      (GPT-4)       │
└────────┬───────────┘
         ▼
┌────────────────────┐
│ Back to Go Server  │
└────────┬───────────┘
         ▼
┌────────────────────┐
│ React Frontend     │
│ displays response  │
└────────────────────┘
```

---

## 🔍 Key Features

- **Persona-Based Prompting**: Uses a tone and voice based on my personality
- **Semantic Chunk Embedding**: Only resume `chunks` are embedded and compared for similarity using cosine distance
- **Chat Memory**: MongoDB persistence of chat messages based on user
- **Simple Project Referencing**: GPT responses may mention projects naturally — no backend project lookups or joins
- **Secure by Default**: No dynamic `resume.json` exposure or external writes

---

## 📁 File Structure

```
resume-chatbot/
├── client/                # React frontend
│   └── src/
├── server/                # Go backend
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── openai/
│   │   ├── embedding.go
│   │   └── open_ai.go
│   ├── resume/
│   │   ├── resume.go
│   │   └── resume.json
│   ├── config/
│   │   └── env.go
│   ├── db/
│   │   ├── models.go
│   │   └── mongo.go
├── README.md
├── .env.example          # API keys and config template

```

---

## 🚀 Running the Project

1. **Clone the repo**

```bash
git clone https://github.com/yourusername/resume-chatbot.git
```

2. **Set up your `.env`**

```bash
cp .env.example .env
```

Update the following:

```env
OPENAI_API_KEY=your-key-here
FRONTEND_ORIGIN=http://localhost:5173
OPENAI_API_URL=https://api.openai.com/v1/chat/completions
OPENAI_EMBEDDING_URL=https://api.openai.com/v1/embeddings
PORT=8080
MONGO_URI=
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

---

## ✅ Notes & Limitations

- Projects and awards are embedded as text chunks — no lookup by ID
- All responses are generated from the retrieved chunks only

---

## 🧪 Ideas for Future Improvements

- [ ] Stream GPT output to frontend
- [ ] Resume.json editor with live preview
- [ ] GPT function-calling for structured answers (e.g. job search tools)
- [ ] Need to send back json objects of resume data to display images, and more on projects mentioned
