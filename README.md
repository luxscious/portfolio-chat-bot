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

- **Persona-Based Prompting**: Responses mirror my voice and personality
- **Semantic Chunk Embedding**: Resume data is chunked and embedded for similarity-based retrieval
- **UUID-Based Session Tracking**: Each user gets a unique ID saved in `localStorage` to persist history. Not too concerned about XSS attacks here or security worries.
- **Typewriter Animation**: Assistant replies animate one character at a time
- **Input Disable Logic**: Input is disabled while waiting for backend OR while the assistant is typing
- **Message Placeholder**: An empty assistant message renders animated dots (`<Ellipsis />`) before reply
- **Chat Memory**: MongoDB stores per-user chat threads

## 📁 File Structure

```
portfolio-chat-bot/
├── client/         # REACT + VITE FRONTEND
│   ├── src/
│   │   ├── components/
│   │   │   └── ui/
│   │   │       ├── button.tsx
│   │   │       ├── input.tsx
│   │   │       ├── scroll-area.tsx
│   │   │       ├── ChatWindow.tsx
│   │   │       ├── MessageBubble.tsx
│   │   │       └── MessageInput.tsx
│   │   ├── hooks/
│   │   │   └── TypingEffect.ts
│   │   ├── lib/
│   │   │   └── utils.ts
│   │   ├── pages/
│   │   │   └── ChatPage.tsx
│   │   ├── styles/
│   │   │   └── index.css
│   │   ├── types/
│   │   │   └── index.ts
│   │   ├── App.tsx
│   │   ├── ChatPage.ts
│   │   ├── main.tsx
│   │   └── vite-env.d.ts
│   ├── .env
│   ├── .env.example
│   ├── components.json
│   ├── eslint.config.js
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.app.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── server/                 # GO BACKEND
│   ├── config/
│   │   └── env.go
│   ├── db/
│   ├── openai/
│   ├── resume/
│   │   ├── resume.go
│   │   └── resume.json
│   ├── .env
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── routes.go
├── .gitignore
└── README.md

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

## Bugs to look into

- The backend logic is kinda funky right now.
- If a user loads the page and the connection to server is bad, but then is restored, need to be able to display that the bot isnt connected to the server.. Convo history also should be restored when connection resumes.

---

## 🧪 Ideas for Future Improvements

- [ ] Resume.json editor with live preview
- [ ] Need to send back json objects of resume data to display images, and more on projects mentioned
