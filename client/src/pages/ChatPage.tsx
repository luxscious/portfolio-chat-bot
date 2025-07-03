import ChatWindow from "@/components/ChatWindow";
import MessageInput from "@/components/MessageInput";
import AnimatedIcon from "@/components/ui/animated_icon";
import { Button } from "@/components/ui/button";
import { useChat } from "@/hooks/ChatPage";
import { Github, Linkedin, Minus, Monitor } from "lucide-react";

export default function ChatPage() {
  const { messages, sendMessage, isLoading, isTyping, setIsTyping } = useChat();

  return (
    <div className="flex h-screen">
      {/* Left Sidebar */}
      <div className="w-1/4 h-full bg-gray-900 p-4 flex flex-col items-center">
        {/* Profile Image */}
        <div className="flex h-1/3 justify-center">
          <div className="relative w-full max-w-xs aspect-square mx-auto">
            {/* Blurred Background */}
            <div
              className="absolute inset-0 bg-center bg-cover blur-lg scale-110 rounded"
              style={{
                backgroundImage: `url(${
                  import.meta.env.VITE_BLOB_BASE_URL
                }/PFP.jpg)`,
              }}
            />
            {/* Foreground Image */}
            <img
              src={`${import.meta.env.VITE_BLOB_BASE_URL}/PFP.jpg`}
              alt="Profile"
              className="relative z-10 w-full h-full object-cover object-[center_30%]"
            />
          </div>
        </div>

        {/* Intro Text and Links */}
        <div className="flex h-1/3 w-full flex-col px-12 py-6 text-white">
          <p className="font-['Inter'] pb-4 max-w-2xl">
            Hi there! This site is an interactive chatbot I built to showcase my
            experience and projects in a more engaging way. Instead of just
            reading a static resume, you can ask questions and explore my work
            dynamically. Feel free to start a chat or connect with me!
          </p>

          <div className="flex items-center justify-between">
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

            {/* - + Signature */}
            <div className="flex items-center">
              <Minus className="w-5 h-5 mt-4" strokeWidth={0.5} />
              <img
                className="h-16"
                src={`${import.meta.env.VITE_BLOB_BASE_URL}/signature.svg`}
                style={{ filter: "invert(1)" }}
                alt="Signature"
              />
            </div>
          </div>
        </div>
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

      {/* Right Panel */}
      <div className="w-1/4 bg-gray-900 p-4">
        {/* Right panel (maybe project cards?) */}
      </div>
    </div>
  );
}
