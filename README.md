# Teamcity CLI

This project provides a handful of specific interactions with the TeamCity
9.x API. It is extremely rough, but satisfies a need I have right now.

Later we may add additional functionality and solve some of its obvious
shortcomings.

**WARNING: it sends passwords as plaintext, as part of the URL.**

Right now its only use case is for automatically publishing meta runners.

## To install this command.

First you will need to set up Go on your machine. The Go website has good instructions
for doing this at https://golang.org/doc/install

Then run `go get github.com/opentable/teamcity-cli/cmd/teamcity` to install.

## To publish a meta runner.

You will need the URL of the TeamCity server, as well as your username and
password.

Create a file named `$FILENAME.xml` replacing $FILENAME with a real file name
in your current working directory. First, you will need to upload this file
manually as a meta runner in the root project, via the TeamCity web UI.

Then run this command to overwrite it with the latest version of your meta
runner at `$FILENAME.xml`.

```shell
teamcity -user $YOUR_USERNAME -password $YOUR_PASSWORD  -baseurl $TEAMCITY_BASE_URL -action set-meta-runner -data $FILENAME
```

