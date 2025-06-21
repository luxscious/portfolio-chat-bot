# 📘 Lessons Learned – Resume Assistant (June 21, 2025)

## 🎯 Goal

Build a highly accurate and natural-sounding resume chatbot that answers questions about Gabriella's work and education history using OpenAI's tools and structured `resume.json` data.

---

## 🚧 Issues Encountered

### 1. Accuracy in Resume Context Retrieval

- **RAG Logic (Local Embedding)**:

  - Didn’t return specific or context-aware chunks.
  - Loosely matched vectors led to vague responses.

- **Passing Full Resume to GPT**:

  - Worked well in accuracy, but consumed far too many tokens.
  - Not scalable or cost-effective.

- **Assistant + Vector Store**:

  - Ideal in theory — lets OpenAI handle semantic retrieval.
  - Still struggled with relevance due to lack of metadata filtering (e.g. featured, tags).

**🔁 Lesson**: Accuracy depends heavily on retrieval quality. Smart filtering (by tags, priority, context) is essential.

---

### 2. Cost & Latency with Hosted Vector DB

- **OpenAI Assistant API** adds overhead:

  - `CreateThread → Run → Poll → FetchMessage` = long response time (\~4–6s avg).
  - Each step adds latency + token usage.

- **Token Costs**:

  - Assistants consume system prompt + vector match + all thread messages — even for short replies.

- **Trade-off**:

  - No infra to manage ✅
  - But less speed and more token cost ❌

**🔁 Lesson**: Self-hosted RAGFlow might be worth revisiting for speed, cost, and control.

---

### 3. Inconsistent Response Formatting

- Assistant was instructed to return:

```json
{
  "response": "text reply",
  "ids": ["val-t"]
}
```

- **Issues**:

  - Returned raw strings instead of JSON.
  - Escaped JSON string (double decoding needed).
  - Malformed escape characters / invalid JSON.

**Workaround**:

- Added fallback parsing if JSON fails.
- Improved system prompt to enforce JSON output.

**🔁 Lesson**: Need stricter system prompts or function-calling enforcement to ensure formatting consistency.

---

### 4. Relevance & Project Prioritization

- Assistant frequently repeated the same projects (e.g. `HIVE`).
- High school projects (e.g. `sailor-dash`) mentioned before award-winning work (`travel-buddy`).

**Root Cause**:

- Tags like `featured`, `award`, or `professional` weren’t being used during retrieval.
- Lack of explicit filtering logic inside assistant memory.

**🔁 Lesson**: Must leverage `tags` + `featured: true` in either vector metadata or system prompt to enforce context relevance.

---

## ✅ Overall Takeaways

- Start simple but plan to scale: local embedding -> OpenAI vector -> maybe back to RAGFlow with metadata control.
- Structure your resume with `tags`, `featured`, and `ids` early.
- System prompts are powerful, but they aren’t guaranteed — have fallback logic.
- Balance infra effort with token cost. OpenAI simplifies things, but every shortcut costs tokens and time.

---

Let’s revisit RAGFlow with metadata support + smart filters next.
