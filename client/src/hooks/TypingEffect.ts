import { useEffect, useState } from "react";

interface TypingEffectResult {
  displayedText: string;
  isTyping: boolean;
}

export function useTypingEffect(
  text: string | undefined,
  baseSpeed: number = 20
): TypingEffectResult {
  const [displayedText, setDisplayedText] = useState("");
  const [isTyping, setIsTyping] = useState(false);

  useEffect(() => {
    if (!text || text.length === 0) {
      setDisplayedText("");
      setIsTyping(false);
      return;
    }

    setDisplayedText("");
    setIsTyping(true);

    let index = 0;
    let timeoutId: number;

    const typeChar = () => {
      const currentChar = text.charAt(index);
      setDisplayedText((prev) => prev + currentChar);
      index++;

      if (index < text.length) {
        const delay = [",", ".", "!", "?"].includes(currentChar)
          ? baseSpeed * 10
          : baseSpeed;

        timeoutId = window.setTimeout(typeChar, delay);
      } else {
        setIsTyping(false);
      }
    };
    typeChar();
    return () => clearTimeout(timeoutId);
  }, [text, baseSpeed]);

  return { displayedText, isTyping };
}
