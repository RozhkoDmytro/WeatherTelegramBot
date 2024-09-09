## Weather Telegram-bot

## Description
This is a simple Telegram bot that can provide information about current holidays and weather. 
You can also subscribe and unsubscribe to the bot's notifications.

## Examples:
Set secret environment variables to .env file (firstly You sude rename ".sample.env" file) 

/start - provides a reply keyboard with flags of 6 countries. When the user presses the flag it should provide what holiday is today in this country.

/weather - get the current weather for your location

/subcsribe -   subscription to the weather report,

/unsubcsribe - unsubscription to the weather report,

/about – provides short info about you

/links – provides a list of your social links (GitHub, Linkedin, etc)

/help – shows a list of commands with reply markup

## Badges
![alt text](/storage/img/example.png)
![alt text](/storage/img/flags.png)
![alt text](/storage/img/weather.png)

## Getting Started
- Docker (for containerized setup)
- Go 1.16+ (for local setup)
- MongodB
- Telegram token

## Running with Docker
Ensure you have Docker installed.
Run the following command to start up the containers:
```
docker-compose up --build
```

## License

This project is licensed under the MIT License
