# telegram_alert_sender
Alert sender that receives a message from Alertmanager and sends a message to telegram depending on the label with chat_id on the Kubernetes namespace.


Create telegram bot
Add the bot to group/channel
Label namespace chat_id=<group's chat_id>

Add bot's key to config.yaml
Run
Message is listened at  <ip>:9270/alerts
