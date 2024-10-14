#!/bin/bash


osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 5 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 10 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 15 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 20 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 25 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 30 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 35 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 40 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 45 && nc localhost 8989"'
osascript -e 'tell app "Terminal" to activate' -e 'tell app "Terminal" to do script "sleep 50 && nc localhost 8989"'

go run main.go