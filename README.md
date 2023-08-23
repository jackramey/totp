# totp

A TOTP Client and Server inspired by [this](https://news.ycombinator.com/item?id=37211021) hacker news post about how
simple it is to write a TOTP client in Go. This was a great example of how to take the technical document
[RFC6238](https://datatracker.ietf.org/doc/html/rfc6238) and implement it. There were a few things about the post
that could improved upon for fun practice.

1. Make a Server that could produce a new secret and then subsequently verify the TOTP.
2. Make the TOTP implementation a library so that it could be called from both Client and Server.
3. Make the code testable.
