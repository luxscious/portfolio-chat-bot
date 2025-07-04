import { useState } from "react";
import { Button } from "@/components/ui/button";
import { ArrowUp } from "lucide-react";
import { Input } from "@/components/ui/input";

interface MessageInputProps {
  onSend: (message: string) => void;
  isDisabled: boolean;
}

export default function MessageInput({
  onSend,
  isDisabled,
}: MessageInputProps) {
  const [input, setInput] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = input.trim();
    if (!trimmed || isDisabled) return;

    onSend(trimmed);
    setInput("");
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="flex justify-between items-center h-max p-4 bg-[#40414f] w-9/10 mx-auto rounded-3xl"
    >
      <Input
        type="text"
        placeholder="Ask me anything..."
        value={input}
        onChange={(e) => setInput(e.target.value)}
        disabled={isDisabled}
        className="flex-1 mr-3 text-white bg-transparent border-none focus:outline-none focus:ring-0 focus-visible:ring-0 shadow-none disabled:opacity-50"
      />

      <Button
        type="submit"
        variant="secondary"
        disabled={isDisabled}
        className="size-8 p-0 bg-white hover:bg-gray-200 rounded-4xl disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <ArrowUp strokeWidth={3} className="text-black" />
      </Button>
    </form>
  );
}
