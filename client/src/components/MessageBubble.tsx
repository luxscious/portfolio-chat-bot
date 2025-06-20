import { Message } from "@/types";
import { useTypingEffect } from "@/hooks/TypingEffect";
import { Ellipsis } from "lucide-react";
import { useEffect } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";

interface Props {
  message: Message;
  onTypingDone?: () => void;
}

export default function MessageBubble({ message, onTypingDone }: Props) {
  const isUser = message.role === "user";
  const isEmptyAssistant =
    message.role === "assistant" && message.content === "";

  const { displayedText, isTyping } = useTypingEffect(
    isUser || isEmptyAssistant || message.disableAnimation
      ? ""
      : message.content,
    20
  );
  useEffect(() => {
    const hasContent = message.content.trim().length > 0;

    if (
      !isTyping &&
      message.role === "assistant" &&
      hasContent &&
      onTypingDone
    ) {
      onTypingDone();
    }
  }, [isTyping]);
  return (
    <div
      className={`flex items-start gap-2 ${
        isUser ? "justify-end" : "justify-start"
      }`}
    >
      {/* Bot Avatar */}
      {!isUser && (
        <Avatar className="w-18 h-18">
          <AvatarImage
            src={`${import.meta.env.VITE_BLOB_BASE_URL}/avatar.png`}
          />
          <AvatarFallback>GG</AvatarFallback>
        </Avatar>
      )}

      {/* Message Bubble */}
      <div
        className={`px-4 py-3 max-w-2xl text-sm rounded-md break-words whitespace-pre-wrap ${
          isUser ? "bg-[#444654] text-gray-100" : "text-white"
        }`}
      >
        {isEmptyAssistant ? (
          <Ellipsis className="animate-pulse text-white w-5 h-5" />
        ) : (
          <>
            {isUser || message.disableAnimation
              ? message.content
              : displayedText}
            {!isUser && isTyping && (
              <span className="animate-pulse inline-block w-[1ch]">|</span>
            )}
          </>
        )}
      </div>
    </div>
  );
}
