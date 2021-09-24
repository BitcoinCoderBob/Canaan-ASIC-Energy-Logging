Log into your Canaan miner. Go to the log screen. Open the network tab (depends on your browser but ctrl+shift+j open console in chrome then click network tab). Double click on the first one. Not sure it matters which one. It will open in a new tab. Copy the URL. 

Head to terminal and run: `go run main.go -u <your-url-here> -s <desired-capture-duration-in-seconds>

The `-s` flag will determine how often the program pings your logs to take a reading. Note that the logs appear to be updating every 15 seconds, which I set as the default. If you would like to use this default, you need not pass this flag.

Note when I pasted it into terminal, it got automatically formatted from: `http://192.168.1.123/updatecglog.cgi?num=0.12345` to `http://192.168.1.123/updatecglog.cgi\?num\=0.12345` so make sure it is formatted this way.

I rebooted my computer, and without signing into the miner via the borwser as usual, just using my previously saved URL seemed to still work.

This works with my Canaan Miner AvalonMiner. Let me know if it does or does not work with your Canaan ASIC.