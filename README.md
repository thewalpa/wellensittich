# Wellensittich DC
The most wellensittich discord bot

# Bot Features
- Voice send and receive, toggle voice receiving with the /listen command
- Guild-features: Toggle individual bot features on or off as needed
- GIF transcribing: If currently listening, the bot transcribes all voice channel audio to GIFs and posts them in the corresponding text channel 
- Play music from links (supports mainly youtube, you can try without any promise every page mentioned here https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md)

# How to Run
- add a config.json with at least a bot token
- for the GIF transcribing feature you have to have an instance of https://github.com/ahmetoner/whisper-asr-webservice running and add the address to the config + you have to add an API key for the tenor GIF search
- Note: for debugging purposes the bot will currently store sound files with random names in a ./tmp directory if it exists
- for playing music yt-dlp and ffmpeg need to be available in PATH 

### All configurations

|Name|Default|Description|
|---|---|---|
|token|-|Discord Bot token|
|dev_server|-|Discord server ID for development of slash commands|
|whisper_asr_webservice|-|Full URL:port of the Whisper ASR Webservice|
|tenor_key|-|Tenor API key|