## Telegram bot messages
Send messages to telegram users through telegram bot.</br>
</br>
Usage example in **tm_text.go**. Edit **tm.json** initialization file first, set **token** with your bot token ID and **chat_id** with your chat ID.</br>
Function **ApiRequestJson()** returns pure Telegram structure response and error. Error is not nil if any error occurs or HTTP response code is not 200.
