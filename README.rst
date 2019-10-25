.. image:: https://www.sparkpost.com/sites/default/files/attachments/SparkPost_Logo_2-Color_Gray-Orange_RGB.svg
    :target: https://www.sparkpost.com
    :width: 200px

`Sign up`_ for a SparkPost account and visit our `Developer Hub`_ for even more content.

.. _Sign up: https://app.sparkpost.com/join?plan=free-0817?src=Social%20Media&sfdcid=70160000000pqBb&pc=GitHubSignUp&utm_source=github&utm_medium=social-media&utm_campaign=github&utm_content=sign-up
.. _Developer Hub: https://developers.sparkpost.com

SparkPost Go API client
=======================

.. image:: https://travis-ci.org/SparkPost/gosparkpost.svg?branch=master
    :target: https://travis-ci.org/SparkPost/gosparkpost
    :alt: Build Status

.. image:: https://coveralls.io/repos/SparkPost/gosparkpost/badge.svg?branch=master&service=github
    :target: https://coveralls.io/github/SparkPost/gosparkpost?branch=master
    :alt: Code Coverage  
    
.. image:: https://img.shields.io/badge/godoc-gosparkpost-blue.svg
    :target: https://godoc.org/github.com/SparkPost/gosparkpost
    :alt: Go Doc


The official Go package for using the SparkPost API.

Installation
------------

Install from GitHub using `go get`_:

.. code-block:: bash

    $ go get github.com/SparkPost/gosparkpost

.. _go get: https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies

Get a key
---------

Go to `API & SMTP`_ in the SparkPost app and create an API key. We recommend using the ``SPARKPOST_API_KEY`` environment variable. The example code below shows how to set this up.

.. _API & SMTP: https://app.sparkpost.com/#/configuration/credentials

Send a message
--------------

Here at SparkPost, our "send some messages" api is called the `transmissions API`_ - let's use it to send a friendly test message:

.. code-block:: go

    package main

    import (
      "log"
      "os"

      sp "github.com/SparkPost/gosparkpost"
    )

    func main() {
      // Get our API key from the environment; configure.
      apiKey := os.Getenv("SPARKPOST_API_KEY")
      cfg := &sp.Config{
        BaseUrl:    "https://api.sparkpost.com",
        ApiKey:     apiKey,
        ApiVersion: 1,
      }
      var client sp.Client
      err := client.Init(cfg)
      if err != nil {
        log.Fatalf("SparkPost client init failed: %s\n", err)
      }

      // Create a Transmission using an inline Recipient List
      // and inline email Content.
      tx := &sp.Transmission{
        Recipients: []string{"someone@somedomain.com"},
        Content: sp.Content{
          HTML:    "<p>Hello world</p>",
          From:    "test@sparkpostbox.com",
          Subject: "Hello from gosparkpost",
        },
      }
      id, _, err := client.Send(tx)
      if err != nil {
        log.Fatal(err)
      }

      // The second value returned from Send
      // has more info about the HTTP response, in case
      // you'd like to see more than the Transmission id.
      log.Printf("Transmission sent with id [%s]\n", id)
    }

.. _transmissions API: https://www.sparkpost.com/api#/reference/transmissions

Documentation
-------------

* `SparkPost API Reference`_
* `Code samples`_
* `Command-line tool: sparks`_

.. _SparkPost API Reference: https://developers.sparkpost.com/api
.. _Code samples: examples/README.md
.. _Command-line tool\: sparks: cmd/sparks/README.md

Contribute
----------

TL;DR:

#. Check for open issues or open a fresh issue to start a discussion around a feature idea or a bug.
#. Fork `the repository`_.
#. Go get the original code - ``go get https://github.com/SparkPost/gosparkpost``
#. Add your fork as a remote - ``git remote add fork http://github.com/YOURID/gosparkpost``
#. Make your changes in a branch on your fork
#. Write a test which shows that the bug was fixed or that the feature works as expected.
#. Push your changes - ``git push fork HEAD``
#. Send a pull request. Make sure to add yourself to AUTHORS_.

More on the `contribution process`_

.. _`the repository`: https://github.com/SparkPost/gosparkpost
.. _AUTHORS: AUTHORS.rst
.. _`contribution process`: CONTRIBUTING.md

