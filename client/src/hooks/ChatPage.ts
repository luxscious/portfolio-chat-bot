import { useState } from "react";
import { Message } from "@/types";

export function useChat() {
  const [messages, setMessages] = useState<Message[]>([]);

  const sendMessage = (msg: string) => {
    const newMessages: Message[] = [
      {
        role: "user",
        content: msg,
      },
    ];

    setMessages(newMessages);

    // TODO: Replace with actual fetch to backend
    // Loading logic prob should be here
    setTimeout(() => {
      setMessages([
        ...newMessages,
        { role: "assistant", content: "This is a stubbed response." },
      ]);
    }, 500);
  };

  return { messages, sendMessage };
}
