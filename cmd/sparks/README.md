# sparks

`sparks` is a command-line tool for sending email using the [SparkPost API](https://developers.sparkpost.com/api/).

### Why does this exist?

I've found this tool useful for testing and troubleshooting, and for making sure that reported issues with `gosparkpost` get squashed. It also has handy working example code showing how to use the `gosparkpost` library to do various things, like using inline images, adding attachments, and managing cc/bcc recipients.

### Why is it called `sparks`?

It's similar in function to [swaks](http://www.jetmore.org/john/code/swaks/), which is a handy SMTP tool: `Swiss Army Knife for SMTP`, and `sparks` sounded better than `swaksp`, `swakapi` or `apisak`, etc.

### Installation

    $ go get git@github.com:SparkPost/gosparkpost
    $ cd $GOPATH/src/github.com/SparkPost/gosparkpost/cmd/sparks
    $ go build && go install

### Config

    $ export SPARKPOST_API_KEY=0000000000000000000000000000000000000000

### Usage Examples

HTML content with inline image, dumping request and response HTTP headers, and response body.

    $ sparks -from img@sp.example.com -html $HOME/test/image.html \
      -img image/jpg:hi:$HOME/test/hi.jpg -subject 'hello!' \
      -to me@example.com.sink.sparkpostmail.com -httpdump

HTML content with attachment.

    $ sparks -from att@sp.example.com -html 'Did you get that <i>thing</i> I sent you?' \
      -to you@example.com.sink.sparkpostmail.com -subject 'that thing' \
      -attach image/jpg:thing.jpg:$HOME/test/thing.jpg

Text content with cc and bcc, but don't send it.
The output that would be sent is pretty printed using the handy JSON tool `jq`.

    $ sparks -from cc@sp.example.com -text 'This is an ambiguous notification' \
      -subject 'ambiguous notification' \
      -to me@example.com.sink.sparkpostmail.com \
      -cc you@example.com.sink.sparkpostmail.com \
      -bcc thing1@example.com.sink.sparkpostmail.com \
      -bcc thing2@example.com.sink.sparkpostmail.com \
      -dry-run | jq .
