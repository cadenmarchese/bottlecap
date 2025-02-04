## About
bottlecap is a simple Go client for interacting with locally-run LLMs available through Podman Desktop's AI lab. The app is meant to provide a more streamlined experience for interacting with these LLMs in chatbot style via CLI. For more details on Podman Desktop's AI Lab, including installation instructions, please see [the official documentation](https://podman-desktop.io/docs/ai-lab). 

## Instructions
1. Install Podman Desktop, the AI Lab extension, and an LLM of your choice, making note of the host and port where your LLM is serving requests.

1. `git clone https://github.com/cadenmarchese/bottlecap`

1. Within the repo, modify the config file accordingly, where URL is the host and port where your local LLM is served, and instructions are the LLM's system role:
    ```bash
    {
        "url": "http://localhost:54277/v1/chat/completions",
        "method": "POST",
        "instructions": "You are a helpful assistant."
    }
    ```
1. Ask away.

    ```bash
    ./bottlecap --help
    
    NAME:
        ask - ask the local LLM a question, in quotes, as the first argument to the "ask" subcommand

    USAGE:
        ask [global options]

    GLOBAL OPTIONS:
        --help, -h  show help
    ```
    
    ```bash
    ./bottlecap ask "Why is the sky blue?"
    ```