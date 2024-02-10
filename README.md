# benevolent
IRC Bot

An attempt at learning Golang and creating an IRC bot that can idle with the best of us.

__Disclaimer:__ This project is intended to be used by myself to help me learn Golang. It's likely going to have bugs. So just a heads up, there are other IRC bots out there on GitHub that are more battle-tested and are programmed by people who know what they are doing.

I do recommend the journey though. It's been good to learn how to handle network connections with Golang, then build up from there. I am glad TLS support is already inbuilt... that was going to be interesting if you need to roll your own support... I've heard it more than once that it's never a great idea to roll your own own PKI... I imagine that applies to enabling encryption on your network connections. :D

## Attributions

- ChatGPT has heavily supported me in this effort
- [Zek](https://github.com/miku/zek) saved me from having to spend hours figuring out mapping XML to Golang Structs.

# TODO

- ~~tell us about the weather~~
- save a message and send it to another user when they join the chat
- save interesting urls
- ~~use a secure connection~~
- giphy and or tenor integration (interesting to see if there are IRC clients that support loading images)
- ~~add allow list to control certain functionality~~
- ~~implement using a registered user~~
- Get trusted users from the database, rather than an env file.

# Build

```
docker-compose build
```

# Run

We use docker-compose in combination with a `.env` file with our values. See the `env.example` file for an example of the values you can use.

```
docker-compose up
```

To stop, you can hit `ctrl-c` 

More info on [docker-compose](ihttps://docs.docker.com/compose/).

