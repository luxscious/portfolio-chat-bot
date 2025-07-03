FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    curl \
    git \
    libatomic1

RUN curl -fsSL https://ollama.com/install.sh | sh

EXPOSE 11434

CMD ["ollama", "serve"]
