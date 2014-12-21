Masspong server in golang
=========================

ugly code for chaosdorf GameJam #0

also needed:

[masspong-web](https://github.com/SkaveRat/masspong-web)

optional:

[masspong-client](https://github.com/SkaveRat/masspong-client)


Protocol
--------

`echo "1 down" | nc localhost 1337`

first part is the player, second the direction


Ports
-----

server: runs on 8080
web: needs to run on localhost:8000
clients: connect to 1337
