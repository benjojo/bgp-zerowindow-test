bgp-zerowindow-test
===

This is a test to reproduce undesirable behaviour in almost all (not OpenBGPd) BGP deamons by causing a zero window TCP window to happen.

This code is heavily modified https://github.com/cloudflare/fgbgp , the main binary is in the [trigger-bug](/trigger-bug) directory. 

This was written for the investigation of: https://blog.benjojo.co.uk/post/bgp-stuck-routes-tcp-zero-window and the internet draft https://datatracker.ietf.org/doc/html/draft-spaghetti-idr-bgp-sendholdtimer
