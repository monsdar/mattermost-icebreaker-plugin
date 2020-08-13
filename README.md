# Mattermost Icebreaker Plugin
This plugin adds the ability to ask random users Icebreaker questions in a channel.

```
Mike:           /icebreaker
IceBreaker Bot: Hey John! Emacs or Vim?
John:           VSCode! But with Vim bindings...
```

```
Mike:           /icebreaker add What's your favorite sports?
IceBreaker Bot: Thanks Mike! Added your question: 'What's your favorite sports?'. Total number of questions: 1
Mike:           /icebreaker
IceBreaker Bot: Hey John! What's your favorite sports?
John:           I love playing basketball! Anyone else here? We can meet up tomorrow and play some 3on3...
```

```
Mike:           /icebreaker
IceBreaker Bot: Hey John! What's your favorite superhero?
John:           Wtf Mike! Stop triggering this bot every 5 minutes!
```

## Why?
In COVID times it's hard to get to know your colleagues by casually chatting by the watercooler. This bot enables these type of random interactions between everyone.

## Features
* Everyone can trigger a new Icebreaker question using `/icebreaker`
* Everyone can add new questions: `/icebreaker add <question>`
* Global list of questions, bot can be triggered in any channel and it asks a random online user from that channel
* Fill in a bunch of default questions using `/icebreaker reset questions`

## Contribute
This plugin is based on the [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template). See there on how to set everything up and test the plugin.

Feel free to post any issues and features to the [issue tracker](https://github.com/monsdar/mattermost-icebreaker-plugin/issues).

## Attributions
The icecube logo is licensed under Creative Commons: `ice cube by 23 icons from the Noun Project`
