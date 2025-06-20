import { useState } from "react";
import { Message } from "@/types";
import { v4 as uuidv4 } from "uuid";

export function useChat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isTyping, setIsTyping] = useState(false);

  const getUserId = () => {
    let userId = localStorage.getItem("userId");
    if (!userId) {
      userId = uuidv4();
      localStorage.setItem("userId", userId);
    }
    return userId;
  };

  const userId = getUserId();

  const sendMessage = async (msg: string) => {
    const userMessage: Message = { role: "user", content: msg };
    setMessages((prev) => [...prev, userMessage]);

    // Add empty assistant message for loading
    const loadingPlaceholder: Message = { role: "assistant", content: "" };
    setMessages((prev) => [...prev, loadingPlaceholder]);
    setIsTyping(true);
    setIsLoading(true);

    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/chat`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userId: userId, content: msg }),
      });

      const data = await res.json();
      const botMessage: Message = { role: "assistant", content: data.content };

      // Replace the last message (placeholder) with the actual response
      setMessages((prev) => [...prev.slice(0, -1), botMessage]);
    } catch (err) {
      // Replace placeholder with error message
      setMessages((prev) => [
        ...prev.slice(0, -1),
        { role: "assistant", content: "Sorry, something went wrong." },
      ]);
    } finally {
      setIsLoading(false);
    }
  };
  return { messages, sendMessage, isLoading, isTyping, setIsTyping };
}
