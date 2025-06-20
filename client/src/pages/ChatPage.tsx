import ChatWindow from "@/components/ChatWindow";
import MessageInput from "@/components/MessageInput";
import { useChat } from "@/hooks/ChatPage";

export default function ChatPage() {
  const { messages, sendMessage } = useChat();
  return (
    <div className="flex h-screen">
      {/* Left Sidebar */}
      <div className="w-1/4 bg-gray-900 p-4">{/* Put Log top left here */}</div>

      {/* Center Chatbot */}
      <div className="w-2/4 bg-gray-900 p-4 flex items-center justify-center">
        <div className="w-full max-w-md bg-gray-200 p-6 rounded shadow">
          <ChatWindow messages={messages} />
          <MessageInput onSend={sendMessage} />
        </div>
      </div>

      {/* Right Projects */}
      <div className="w-1/4  bg-gray-900 p-4">
        {/* Right panel (maybe project cards?) */}
      </div>
    </div>
  );
}
