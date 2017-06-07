## fblgen

Testing your response to FBL reports doesn't have to involve waiting for an angry/lazy recipient to click "This is Spam".
Here's how to send an FBL report in response to a message sent via SparkPost, and saved (with full headers) to a local file:

    $ ./fblgen --file ./test.eml --verbose
    Got domain [sparkpostmail.com] from Return-Path
    Got MX [smtp.sparkpostmail.com.] for [sparkpostmail.com]
    Would send FBL from [test@sp.example.com] to [fbl@sparkpostmail.com] via [smtp.sparkpostmail.com.:smtp]

Note that this command (once you've added the `--send` flag) will attempt to connect from your local machine to the MX listed above.
It's entirely possible that there will be something blocking that port, for example a firewall, or your residential ISP.
Here are [two](http://nc110.sourceforge.net/) [ways](https://nmap.org/ncat/) to check whether that's the case.
Whichever command you run should return in under a second.
If there's a successful connection, you're good to go.

    $ nc -vz -w 3 smtp.sparkpostmail.com 25
    $ </dev/null ncat -vw 3s --send-only smtp.sparkpostmail.com 25

If you get a timeout, there are a couple solutions. The easiest is to `ssh` somewhere that allows outbound connections on port 25. Searching for "free ssh" will give you quite a few options, if you don't happen to have that sort of access set up already. My nostalgic favorite is [SDF](http://sdf.lonestar.org/).

Another option is to route your connections over a VPN, which is more involved, and out of the scope of this document.

Happy testing!
