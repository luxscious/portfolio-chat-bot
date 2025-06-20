import { useState } from "react";
import { Button } from "@/components/ui/button";
import { ArrowUp } from "lucide-react";
import { Input } from "@/components/ui/input";
interface MessageInputProps {
  onSend: (message: string) => void;
}

export default function MessageInput({ onSend }: MessageInputProps) {
  const [input, setInput] = useState("");
  const [submit, setSubmit] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    // Disable sending logic, so disbale the button
    setSubmit(true);
    e.preventDefault();
    const trimmed = input.trim();
    if (!trimmed) return;
    onSend(trimmed);
    setInput("");

    setSubmit(false);
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="flex justify-between content-center h-max p-4 bg-[#40414f] w-9/10 mx-auto rounded-3xl"
    >
      <Input
        type="text"
        placeholder="Ask me anything..."
        value={input}
        onChange={(e) => setInput(e.target.value)}
        className="text-white bg-transparent border-none focus:outline-none focus:ring-0 focus-visible:ring-0 shadow-none"
      />
      <Button
        disabled={submit}
        type="submit"
        variant={"secondary"}
        className="size-8 bg-white hover:bg-gray-200 rounded-4xl"
      >
        <ArrowUp strokeWidth={3} />
      </Button>
    </form>
  );
}
