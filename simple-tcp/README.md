# Simple TCP Server

A very simple game server, for the purposes of testing a TCP based server on Carrier.

## Server

Starts a server on port `7654` by default. Can be overwritten by `PORT` env var or `port` flag.

When it receives a text message ending with a newline, it will send back "ACK:<text content>" as an echo. 

If it receives the text "EXIT", then it will `sys.Exit(0)`