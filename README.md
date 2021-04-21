bgp-zerowindow-test
===

This is a test to reproduce undesirable behaviour in almost all (not OpenBGPd) BGP deamons by causing a zero window TCP window to happen.

This code is heavily modified https://github.com/cloudflare/fgbgp , the main binary is in the [trigger-bug](/trigger-bug) directory. 

