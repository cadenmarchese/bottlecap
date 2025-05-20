## About
bottlecap is a simple Go client for interacting with OpenAI compatible LLM servers run locally or elsewhere. The app is meant to provide a more streamlined experience for interacting with these LLMs in chatbot style via CLI.

## Instructions (local LLM)

1. `git clone https://github.com/cadenmarchese/bottlecap`

1. Within the repo, modify the config file accordingly, where URL is the host and port where your LLM is served, and instructions are the LLM's system role:
    ```json
    {
        "url": "https://example.com/v1/chat/completions",
        "bearerToken": "<your bearer token>",
        "model": "llama-scout-maas",
        "chatInstructions": "You are a helpful assistant.",
        "imageInstructions": "Describe this image."
    }
    ```
1. Ask away.

    ```bash
    NAME:
    bottlecap - Provide inputs to the bottelecap application using quotes

    USAGE:
    bottlecap [global options]

    GLOBAL OPTIONS:
    --help, -h  show help
    ```
    
    ```bash
    ./bottlecap ask "Why is the sky blue?"
    ```