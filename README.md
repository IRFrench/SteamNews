# Steam News

This project had one goal, give a quick rundown of the steam news for my owned games in Discord.

## Background

Where I work, we use slack, and in slack we have a lot of automation tools that produce messages for us to consume in the application

I figured if I do the same thing in discord, I can look at what I am interesting in and read further into it. And the one thing I want to keep up with but keep forgetting about is game news.

Thus, this project was born. It uses a bot to send DM messages to the configured users with the last 24 hours of recorded steam news.

## Use

Set up the configuration and run the service. The configuration is searched for using the `ETC` environment variable.

An example configuration would look like this

```
discord:
  bot_token: "<BOT TOKEN HERE>"
steam:
  key: "<STEAM WEBAPI KEY HERE>"
start_time:
  hour: 19
  minute: 0
users:
  - name: "test user" // Only used for logging
    discord_id: <DISCORD USER ID>
    steam:
      id: <STEAM USER ID>
      removed:
        - 1592110
        - 1440670
        - 602320
      added:
        - name: "fake_game"
          id: 340
      played_only: true
```

### Discord

All you need for the service is a bot token and user ID.

You can learn about making discord bots here: https://discord.com/developers/docs/quick-start/getting-started

When you have created the bot, you will need to go to installation, set authorization to user install and use the install link to allow the bot to send you a private message

### Steam

For this, you will need a web API token. You can create one here: https://steamcommunity.com/dev/apikey

When you have this, you are mostly set up. Once you have your steam ID in the configuration, you are good to go.

From here you just need to add any configuration you would like for your news.
- `removed` is a list of removed games, games you don't want news for.
- `added` is a list of extra games you wanted added, like any games you don't own
- `played_only` will only collect news for played games. Added games will work even if not played

### Start time

I didn't want to wait until 7 to start the service, so I added a waiter to wait till 7 for me.

## Future work

I want this to be a platform for a full golang discord bot, this was a cool first step. I would like to have the option to add and remove games from the news list directly in the app. Its not as easy to change being in a configuration file.

## Features and Requests

I am most likey not going to work on this project consitently and there are other discord bots out there you can use, but if you like this and want to add features raise an issue or a pull request.

I can't say I will accept your change or features but if it is sensible I'll probably use it.