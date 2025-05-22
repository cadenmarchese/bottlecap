## About
bottlecap is a simple Go client for interacting with OpenAI compatible LLM servers run locally or elsewhere. The app is meant to provide a more streamlined experience for interacting with these LLMs in chatbot style via CLI.

## Instructions (local LLM)

1. `git clone https://github.com/cadenmarchese/bottlecap`

1. Build the binary with `make build` or dowload it from GitHub Releases

1. Within the repo, modify the config file accordingly, or create your own config.json file in the same directory as the binary, with the following fields:
    ```bash
    {
        "url": "https://example.com", # url of the LLM (including port if necessary)
        "bearerToken": "<your bearer token>", # bearer token, if required
        "model": "<optionally-specify-model>", # name / id of model (such as llama-scout-maas)
        "chatInstructions": "You are a helpful assistant.", # system prompt for chat
        "imageInstructions": "Describe this image." # system prompt for "image" subcommand
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
    
    ```bash
    ./bottlecap image "https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg"
    ```

    ```bash
    ./bottlecap generate "Create an image of a dog in sunglasses."
    ```