echo $1
ffprobe -i $1 -show_entries format=duration -v quiet -of csv="p=0"
