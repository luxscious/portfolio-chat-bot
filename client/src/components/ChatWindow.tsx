import { useEffect, useRef } from "react";
import MessageBubble from "./MessageBubble";
import { Message } from "@/types";
import { ScrollArea } from "./ui/scroll-area";

interface ChatWindowProps {
  messages: Message[];
  onTypingDone: () => void;
}
export default function ChatWindow({
  messages,
  onTypingDone,
}: ChatWindowProps) {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <ScrollArea className="flex-1 h-10/12 px-4 py-8 space-y-4">
      {messages.map((msg, i) => (
        <MessageBubble
          key={i}
          message={msg}
          onTypingDone={
            i === messages.length - 1 && msg.role === "assistant"
              ? onTypingDone
              : undefined
          }
        />
      ))}
      <div ref={bottomRef} />
    </ScrollArea>
  );
}
