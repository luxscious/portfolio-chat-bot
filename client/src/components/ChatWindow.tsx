import { useEffect, useRef } from "react";
import MessageBubble from "./MessageBubble";
import { Message } from "@/types";
import { ScrollArea } from "./ui/scroll-area";
import { useBooleanFlagValue } from "@openfeature/react-sdk";
import { Bot, Github, Linkedin, Monitor } from "lucide-react";
import AnimatedIcon from "@/components/ui/animated_icon";

interface ChatWindowProps {
  messages: Message[];
  onTypingDone: () => void;
}
export default function ChatWindow({
  messages,
  onTypingDone,
}: ChatWindowProps) {
  const bottomRef = useRef<HTMLDivElement>(null);
  const serverDown = useBooleanFlagValue("server-down", false);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return serverDown ? (
    <div className="flex-1 h-10/12 px-4 py-8 space-y-4 flex items-center justify-center">
      <div className="flex flex-col items-center justify-center w-full text-center rounded">
        <Bot size={36} className="animate-float text-gray-400" />
        <p className="mb-4 text-sm pt-4 ">
          Sorry! Backend is down for Maintenance or due to Server Costs
        </p>
        <p className="mb-4 text-sm ">
          Feel free to check out my Github, LinkedIn, or personal site.
        </p>
        <div className="flex items-center justify-between text-gray-400">
          {/* Icons */}
          <div className="flex gap-4 h-full">
            <AnimatedIcon href="https://github.com/luxscious">
              <Github />
            </AnimatedIcon>
            <AnimatedIcon href="https://www.linkedin.com/in/gabriella-gerges/">
              <Linkedin />
            </AnimatedIcon>
            <AnimatedIcon href="https://luxscious.github.io/">
              <Monitor />
            </AnimatedIcon>
          </div>
        </div>
      </div>
    </div>
  ) : (
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
