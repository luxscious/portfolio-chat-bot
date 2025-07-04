import ChatWindow from "@/components/ChatWindow";
import MessageInput from "@/components/MessageInput";
import AnimatedIcon from "@/components/ui/animated_icon";
import { Button } from "@/components/ui/button";
import { useChat } from "@/hooks/ChatPage";
import { Github, Linkedin, Minus, Monitor } from "lucide-react";
import ReactMarkdown from "react-markdown";
import rehypeHighlight from "rehype-highlight";
import "highlight.js/styles/github-dark.css";

export default function ChatPage() {
  const { messages, sendMessage, isLoading, isTyping, setIsTyping } = useChat();

  const exampleMarkdown = `
\`\`\`bash
👤 What skills do you have?
🤖 I'm experienced in React, Node.js, and building secure IoT applications.

👤 Can you show me your latest projects?
🤖 Sure! I've recently built an AI resume chatbot and a smart charging system simulator.
\`\`\`
`;

  return (
    <div className="flex flex-col gap-4 lg:flex-row lg:gap-0 h-auto lg:h-screen">
      {/* Left Sidebar */}
      <div className="w-full lg:w-1/4 bg-gray-900 p-4 flex flex-col items-center">
        {/* Profile Image */}
        <div className="flex h-1/3 justify-center">
          <div className="relative w-full max-w-xs aspect-square mx-auto pt-6">
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
        <div className="flex h-1/3 w-full flex-col px-4 py-6 text-white">
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
          </div>
        </div>

        {/* Example Markdown */}
        <div className="mt-8 w-full text-xs leading-relaxed">
          {/* Label */}
          <div className="bg-gray-800 text-gray-300 text-[10px] px-2 py-1 rounded-t-lg w-max">
            example text
          </div>
          {/* Code block */}
          <div className="w-full rounded-b-lg">
            <ReactMarkdown
              rehypePlugins={[rehypeHighlight]}
              components={{
                code({ node, inline, className, children, ...props }) {
                  return !inline ? (
                    <pre className="m-0 whitespace-pre-wrap break-words rounded-b-lg p-2 overflow-x-auto">
                      <code className={className} {...props}>
                        {children}
                      </code>
                    </pre>
                  ) : (
                    <code className={className} {...props}>
                      {children}
                    </code>
                  );
                },
              }}
            >
              {exampleMarkdown}
            </ReactMarkdown>
          </div>
        </div>
      </div>

      {/* Center Chatbot */}
      <div className="w-full lg:w-1/2 bg-gray-900 p-6 flex justify-center pb-12">
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
                target="_blank"
                rel="noopener noreferrer"
              >
                Resumé
              </a>
            </Button>
          </div>
        </div>
      </div>

      {/* Right Panel */}
      <div className="hidden lg:block w-full lg:w-1/4 bg-gray-900 p-4">
        {/* Right panel (maybe project cards?) */}
      </div>
    </div>
  );
}
