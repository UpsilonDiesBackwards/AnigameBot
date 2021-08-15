# AnigameBot
A bot for automating the Discord bot Anigame.

## Information

This bot uses a Discord user account (ie. not a bot account) and will read the contents of the last Anigame message, and will automatically send a response back to it.

Currently the bot can:
1. Automattically battle and advance to the next floor and location
2. (thats it.)

## DISCLAIMER

***THIS BOT VIOLATES TOS!!!***
Please ***DO NOT*** use your main discord account to do this! This could result you in having you account terminated. Please use an alternative discord account with an age of over 3 months (Anigame requires an account to be over 3 months old to be used)

***(ADDITIONALLY, I AM NOT RESPONSIBLE FOR ANY DAMAGE CAUSED BY THIS BOT! IF YOUR ACCOUNT GETS TERMINATED, (L) I WILL NOT BE HELD RESPONSIBLE OR LIABLE!!)***  

## Building Anigame Bot

The building steps of this application are as follows:

1. `go build totallynotillegaldiscordbot`
2. Then: `./totallynotillegaldiscordbot`, (if on Linux)

## Configuring the Bot

On the initial use of the Bot, you will be required to create a config file. To do this you must:

1. Create a file called `"config.json"`
2. Inside that file, use the following as a guide:

```JSON
  {
    "token": "PUT USER ACCOUNT TOKEN HERE",
    "allowed_guilds": [
      "ID OF SERVER TO OPERATE IN"
    ],
    "bot_user_id": "571027211407196161"
  }
```

## Using the Bot

Using the bot is easy, after you configure the bot and run it.
You will have to initate a battle *(.battle)* and the bot will take over for you.
