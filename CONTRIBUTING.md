# Contributing to gosparkpost

Transparency is one of our core values, and we encourage developers to contribute and become part of the SparkPost developer community.

The following is a set of guidelines for contributing to gosparkpost,
which is hosted in the [SparkPost Organization](https://github.com/sparkpost) on GitHub.
These are just guidelines, not rules, use your best judgment and feel free to
propose changes to this document in a pull request.

## Submitting Issues

* You can create an issue [here](https://github.com/sparkpost/gosparkpost/issues/new), but
  before doing that please read the notes below on debugging and submitting issues,
  and include as many details as possible with your report.
* Include the version of gosparkpost you are using.
* Perform a [cursory search](https://github.com/issues?utf8=%E2%9C%93&q=is%3Aissue+user%3Asparkpost+repo%3Agosparkpost)
  to see if a similar issue has already been submitted.

## Local development

* Fork this repository
* Clone the _original source_  using ``go get https://github.com/SparkPost/gosparkpost``
* `cd $GOPATH/src/github.com/SparkPost/gosparkpost`
* Add a remote pointing to your fork - ``git remote add fork https://github.com/YOUR_ID/gosparkpost``
* Create a branch - ``git checkout -b fix-or-feature``
* Write code!
* Commit code!
* Push code with - ``git push fork HEAD``

## Contribution Steps

### Guidelines

- Provide [documentation](http://blog.golang.org/godoc-documenting-go-code) for any newly added code.
- Provide [tests](https://golang.org/doc/code.html#Testing) for any newly added code.
- Run ``go fmt`` before submitting.

1. Create a new branch named after the issue you’ll be fixing (include the issue number as the branch name, example: Issue in GitHub is #8 then the branch name should be ISSUE-8))
2. Write corresponding tests and code (only what is needed to satisfy the issue and tests please)
    * Include your tests in the corresponding package directory in an appropriately-named test file
    * Write code to satisfy the tests
3. Ensure automated tests pass
4. Submit a new Pull Request applying your feature/fix branch to the develop branch

## Testing

Once you are set up for local development:

* Run ``go test`` to test against your current Go environment

