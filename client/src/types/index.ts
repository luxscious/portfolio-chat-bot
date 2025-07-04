export type Message = {
  role: "user" | "assistant";
  content: string;
  disableAnimation?: boolean;
};
