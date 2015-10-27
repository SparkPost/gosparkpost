SparkPost Go API client
=======================

The official Go package for using the SparkPost API.

Installation
------------

Install from GitHub using `go get`_:

.. code-block:: bash

    $ go get https://github.com/SparkPost/go-sparkpost

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

      "github.com/SparkPost/go-sparkpost/api"
      te_api "github.com/SparkPost/go-sparkpost/api/templates"
      tr_api "github.com/SparkPost/go-sparkpost/api/transmissions"
    )

    func main() {
      // Get our API key from the environment; configure.
      apiKey := os.Getenv("SPARKPOST_API_KEY")
      TrAPI, err := tr_api.New(api.Config{
        BaseUrl:    "https://api.sparkpost.com",
        ApiKey:     apiKey,
        ApiVersion: 1,
      })
      if err != nil {
        log.Fatalf("Transmissions API init failed: %s\n", err)
      }

      // Create a Transmission using an inline Recipient List
      // and inline email Content.
      id, err := TrAPI.Create(&tr_api.Transmission{
        Recipients: []string{"someone@somedomain.com"},
        Content:    te_api.Content{
          HTML:    "<p>Hello world</p>",
          From:    "test@sparkpostbox.com",
          Subject: "Hello from go-sparkpost",
        },
      })
      if err != nil {
        log.Fatal(err)
      }

      // tr_api.Response has more details, in case you'd
      // like to see more than the Transmission id.
      log.Printf("Transmission sent with id [%s]\n", id)
    }

.. _transmissions API: https://www.sparkpost.com/api#/reference/transmissions

Documentation
-------------

* `SparkPost API Reference`_

.. _SparkPost API Reference: https://www.sparkpost.com/api

Contribute
----------

TL;DR:

#. Check for open issues or open a fresh issue to start a discussion around a feature idea or a bug.
#. Fork `the repository`_.
#. Go get the original code - ``go get https://github.com/SparkPost/go-sparkpost``
#. Add your fork as a remote - ``git remote add fork http://github.com/YOURID/go-sparkpost``
#. Make your changes in a branch on your fork
#. Write a test which shows that the bug was fixed or that the feature works as expected.
#. Push your changes - git push fork
#. Send a pull request. Make sure to add yourself to AUTHORS_.

More on the `contribution process`_

.. _`the repository`: https://github.com/SparkPost/go-sparkpost
.. _AUTHORS: https://github.com/SparkPost/go-sparkpost/blob/master/AUTHORS.rst
.. _`contribution process`: https://github.com/SparkPost/blob/master/go-sparkpost/CONTRIBUTING.md

