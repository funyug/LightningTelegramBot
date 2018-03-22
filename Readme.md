#Lightning Telegram Bot

***Still in Development***

TODO
* Refactor code
* Implement close_channel, help command
* Implement 2FA for sending
* Option to disable certain commands


Setup
* Clone the project
* Run "dep ensure" to install the dependencies
* Create a telegram bot by visiting https://t.me/botfather and get the token for the bot
* Modify the token for bot in server.go
* Build binary using Go build

Start the bot by running

`LightningTelegramBot --username=YOUR_TELEGRAM_USERNAME --token=YOUR_TELEGRAM_BOT_TOKEN`

Note: Make sure you have lnd running in background and the wallet is unlocked