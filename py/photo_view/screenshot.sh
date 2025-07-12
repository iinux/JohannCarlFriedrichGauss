duration=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 $1)
echo "视频总时长: $duration 秒"
random_time=$(awk -v dur="$duration" 'BEGIN { srand(); print rand() * dur }')
echo "随机时间点: $random_time 秒"
ffmpeg -ss "$random_time" -i $1 -vframes 1 -q:v 2 $1random_screenshot.jpg
