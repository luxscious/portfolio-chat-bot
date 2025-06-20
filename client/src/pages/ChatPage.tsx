import ChatWindow from "@/components/ChatWindow";
import MessageInput from "@/components/MessageInput";
import { Button } from "@/components/ui/button";
import { useChat } from "@/hooks/ChatPage";

export default function ChatPage() {
  const { messages, sendMessage, isLoading, isTyping, setIsTyping } = useChat();
  return (
    <div className="flex h-screen">
      {/* Left Sidebar */}
      <div className="w-1/4 bg-gray-900 p-4">
        {/* Put Log top left here */}
        {/* <img
          src="https://media.licdn.com/dms/image/v2/D4D03AQEWeLpaykX60g/profile-displayphoto-shrink_200_200/B4DZdf9DGXHwAc-/0/1749661529229?e=1755734400&v=beta&t=Kc4jS-8iO4Irp0YfMAzr74sWwIsQcrgxSL8aGOueNb8"
          alt="Bot Avatar"
          className="w-20 h-20 rounded-full object-cover"
        /> */}
      </div>

      {/* Center Chatbot */}
      <div className="w-1/2 bg-gray-900 p-6 flex justify-center pb-12">
        <div className="w-full space-y-4">
          <ChatWindow
            messages={messages}
            onTypingDone={() => setIsTyping(false)}
          />
          <MessageInput
            onSend={sendMessage}
            isDisabled={isLoading || isTyping}
          />

          <p className="text-sm text-gray-400 text-center font-light">
            This chat agent can make mistakes. Please verify information through
            resume.
          </p>

          <div className="flex justify-center">
            <Button variant="secondary">
              <a
                className="font-normal"
                href={`${
                  import.meta.env.VITE_BLOB_BASE_URL
                }/Gabriella_Gerges_Resume.pdf`}
              >
                Resumé
              </a>
            </Button>
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
