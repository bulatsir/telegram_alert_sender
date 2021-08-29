# telegram_alert_sender
Alert sender that receives a message from Alertmanager and sends a message to Telegram's bot depending on the label with chat_id on the Kubernetes namespace.


1. Create a telegram bot
2. Add the bot to group/channel
3. Label namespace chat_id=<group's chat_id>

4. Add bot's key to config.yaml
5. Run (go build main.go && ./main)

Messages is listened at  <ip>:9270/alerts
