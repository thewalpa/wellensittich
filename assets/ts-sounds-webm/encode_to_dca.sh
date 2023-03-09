

for i in *.webm; do ffmpeg -i "$i" -f s16le -ar 48000 -ac 2 pipe:1 | dca > "../ts-sounds-dca/${i%.*}.dca"; done