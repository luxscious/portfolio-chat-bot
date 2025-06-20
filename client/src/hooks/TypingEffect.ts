import { useEffect, useState } from "react";

interface TypingEffectResult {
  displayedText: string;
  isTyping: boolean;
}

export function useTypingEffect(
  text: string,
  baseSpeed: number = 20
): TypingEffectResult {
  const [displayedText, setDisplayedText] = useState("");
  const [isTyping, setIsTyping] = useState(false);

  useEffect(() => {
    setDisplayedText("");
    setIsTyping(true);

    let index = 0;

    const typeChar = () => {
      setDisplayedText((prev) => prev + text.charAt(index));
      index++;

      if (index < text.length) {
        const char = text.charAt(index - 1);
        const delay = [",", ".", "!", "?"].includes(char)
          ? baseSpeed * 10
          : baseSpeed;
        setTimeout(typeChar, delay);
      } else {
        setIsTyping(false);
      }
    };

    const timeout = setTimeout(typeChar, baseSpeed);
    return () => clearTimeout(timeout);
  }, [text, baseSpeed]);

  return { displayedText, isTyping };
}
