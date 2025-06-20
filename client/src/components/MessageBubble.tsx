import { Message } from "@/types";
import { useTypingEffect } from "@/hooks/TypingEffect";

interface Props {
  message: Message;
}

export default function MessageBubble({ message }: Props) {
  const isUser = message.role === "user";
  const { displayedText, isTyping } = useTypingEffect(message.content, 20);

  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
      <div
        className={`px-4 py-3 max-w-2xl text-sm rounded-md break-words whitespace-pre-wrap ${
          isUser ? "bg-[#444654] text-gray-100" : "text-white"
        }`}
      >
        <span>
          {displayedText}
          {isTyping && (
            <span className="animate-pulse inline-block w-[1ch]">|</span>
          )}
        </span>
      </div>
    </div>
  );
}
