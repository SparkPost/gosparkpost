# Example Code

Short snippets of code showing how to do various things.
Feel free to submit your own examples!

### Cc and Bcc

Mail clients usually set up the details of `Cc` and `Bcc` for you, so thinking about it in terms of individual emails to be sent can be a bit of an adjustment. Here's a snippet that shows how it's done. See also [sparks](cmd/sparks/sparks.go) for an example that will send mail, instead of just printing out JSON.

[Cc and Bcc example](cc/cc.go)

### Overview

Replicating mail clients' `Cc` and `Bcc` behavior with SparkPost is easier to reason about if you focus on the individual recipients - everyone that gets an email when you click send - and what the message needs to look like for them. Everyone who gets a message must currently have their own `Recipient` within the `Transmission` that's sent to SparkPost.

Here's the output of the example linked above, which generates a message with 2 recipients in the `To`, 2 in the `Cc`, and 2 in the `Bcc`. Notice we're setting `header_to` in each `Recipient`, and it's always the same. We're also setting the `Cc` header w/in `content`, which is also always the same.

    {
      "recipients": [
        {
          "address": {
            "email": "to1@test.com.sink.sparkpostmail.com",
            "header_to": "to1@test.com.sink.sparkpostmail.com,to2@test.com.sink.sparkpostmail.com"
          }
        },
        {
          "address": {
            "email": "to2@test.com.sink.sparkpostmail.com",
            "header_to": "to1@test.com.sink.sparkpostmail.com,to2@test.com.sink.sparkpostmail.com"
          }
        },
        {
          "address": {
            "email": "cc1@test.com.sink.sparkpostmail.com",
            "header_to": "to1@test.com.sink.sparkpostmail.com,to2@test.com.sink.sparkpostmail.com"
          }
        },
        {
          "address": {
            "email": "cc2@test.com.sink.sparkpostmail.com",
            "header_to": "to1@test.com.sink.sparkpostmail.com,to2@test.com.sink.sparkpostmail.com"
          }
        },
        {
          "address": {
            "email": "bcc1@test.com.sink.sparkpostmail.com",
            "header_to": "to1@test.com.sink.sparkpostmail.com,to2@test.com.sink.sparkpostmail.com"
          }
        },
        {
          "address": {
            "email": "bcc2@test.com.sink.sparkpostmail.com",
            "header_to": "to1@test.com.sink.sparkpostmail.com,to2@test.com.sink.sparkpostmail.com"
          }
        }
      ],
      "content": {
        "text": "This is a cc/bcc example",
        "subject": "cc/bcc example message",
        "from": "test@example.com",
        "headers": {
          "cc": "cc1@test.com.sink.sparkpostmail.com,cc2@test.com.sink.sparkpostmail.com"
        }
      }
    }

