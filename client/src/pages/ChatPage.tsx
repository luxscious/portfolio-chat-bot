import ChatWindow from "@/components/ChatWindow";
import MessageInput from "@/components/MessageInput";
import { Button } from "@/components/ui/button";
import { useChat } from "@/hooks/ChatPage";

export default function ChatPage() {
  const { messages, sendMessage } = useChat();
  return (
    <div className="flex h-screen">
      {/* Left Sidebar */}
      <div className="w-1/4 bg-gray-900 p-4">{/* Put Log top left here */}</div>

      {/* Center Chatbot */}
      <div className="w-1/2 bg-gray-900 p-6 flex justify-center pb-12">
        <div className="w-full space-y-4">
          <ChatWindow messages={messages} />
          <MessageInput onSend={sendMessage} />

          <p className="text-sm text-gray-400 text-center font-light">
            This chat agent can make mistakes. Please verify information through
            resume.
          </p>

          <div className="flex justify-center">
            <Button variant="secondary">Resume</Button>
          </div>
        </div>
      </div>

      {/* Right Projects */}
      <div className="w-1/4  bg-gray-900 p-4">
        {/* Right panel (maybe project cards?) */}
      </div>
    </div>
  );
}
