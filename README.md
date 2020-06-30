# Mattermost Icebreaker Plugin
This plugin adds the ability to ask random users Icebreaker questions in a channel.

```
Mike:           /icebreaker
IceBreaker Bot: Hey John! Emacs or Vim?
John:           VSCode! But with Vim bindings...
```

```
Mike:           /icebreaker add What's your favorite sports?
IceBreaker Bot: Thanks John! Added your proposal: 'What's your favorite sports?'. Total number of proposals: 1
Admin:          /icebreaker show proposals
IceBreaker Bot: Proposed questions:
                1. John: What's your favorite sports?
Admin:          /icebreaker approve 1
IceBreaker Bot: Question has been approved: What's your favorite sports?
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
* Let users propose new questions, admins need to approve them before they can get asked: `/icebreaker add <question>`, `/icebreaker show porposals` and `/icebreaker approve <question-index>`
* Questions are stored per channel, so each channel can have their own set of questions
* Fill in a bunch of default questions using `/icebreaker reset questions`

## Contribute
This plugin is based on the [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template). See there on how to set everything up and test the plugin.

## Attributions
The icecube logo is licensed under Creative Commons: `ice cube by 23 icons from the Noun Project`
