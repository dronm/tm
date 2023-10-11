## Telegram bot messages
Send messages to telegram users through telegram bot.

- Usage example in **tm_text.go**.
- Initialize environment variables:
	- **TM_TEST_BOT_TOKEN** with your bot token ID
	- **TM_TEST_CHAT_ID** with your chat ID
- Function **ApiRequestJson()** returns pure Telegram structure response and error. Error is not nil if any error occurs or HTTP response code is not 200.
