#Lightning Telegram Bot

***Still in Development***

Features
* Commands: wallet_balance, channel_balance, get_address, list_chain_txns, send_coins, open_channel, generate_invoice, list_invoices, lookup_invoice, send_payment, list_payments, close_channel, list_channels
* Bot responds only to the username specified in the argument
* OpenChannel commands uses streaming api that provides update in realtime about channel status

TODO
* Refactor code
* Implement help command
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