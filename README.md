# AUTO MESSAGE SENDER

In this project, we are developing a system that automatically sends 
unsent messages retrieved from the database every 2 minutes. 
The system will obtain the message content, recipient phone number, 
and sending status, ensuring each message is sent only once. 
Additionally, it will cache the messageId values and sending times of 
sent messages in Redis. The project includes APIs to 
control the automatic message sending and to retrieve sent messages.