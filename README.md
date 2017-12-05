![](https://github.com/go-ignite/ignite/raw/master/snapshots/ignite.png)
---

[![Release](http://github-release-version.herokuapp.com/github/go-ignite/ignite/release.svg?style=flat)](https://github.com/go-ignite/ignite/releases/latest)
[![Build Status](https://travis-ci.org/go-ignite/ignite.svg?branch=master)](https://travis-ci.org/go-ignite/ignite)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-ignite/ignite)](https://goreportcard.com/report/github.com/go-ignite/ignite)
[![Chat Telegram](https://img.shields.io/badge/chat-telegram-brightgreen.svg)](https://t.me/joinchat/FuddvBLqni3u5wJmNwnc0w)
[![Microbadger Image](https://images.microbadger.com/badges/image/goignite/ignite.svg)](https://microbadger.com/images/goignite/ignite)
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/go-ignite/ignite/blob/master/LICENSE)

A SS panel for managing multiple users, powered by Go &amp; Docker.

* [Overview](#overview)
* [Features](#features)
* [Snapshots](#snapshots)
* [Install](#install)
* [Build](#build)
* [Contributing](#contributing)
* [License](#license)

# Overview

![https://github.com/go-ignite/ignite/raw/master/snapshots/create.gif](https://github.com/go-ignite/ignite/raw/master/snapshots/create.gif)

# Features

* __Simple to use__ - Easy to create and activate server by one click.
* __Fast__ - Create & destroy server container via docker in seconds.
* __Safe__ - Users are isolated by different containers, it is easy to suspend account by stopping related docker container.
* __Automatic__ - User's account will be suspended automatically by background job if it exceeds the max bandwidth limit or the expired date.
* __Easy to deploy__ - ignite is powered by Go, you can copy and run it by very easy steps, you can also deploy ignite by docker.

# Snapshots

![https://github.com/go-ignite/ignite/raw/master/snapshots/1.png](https://github.com/go-ignite/ignite/raw/master/snapshots/1.png)

![https://github.com/go-ignite/ignite/raw/master/snapshots/2.png](https://github.com/go-ignite/ignite/raw/master/snapshots/2.png)

![https://github.com/go-ignite/ignite/raw/master/snapshots/3.png](https://github.com/go-ignite/ignite/raw/master/snapshots/3.png)

# Install

Please refer to [《ignite中文安装指南》](https://github.com/go-ignite/ignite/wiki)

# Build

To build ignite, you need to prepare your Go development environment first, then follow the steps:

* clone ignite to your go workspace
* ```go build```

# Contributing

Pull request is welcome!

* Fork ignite
* Clone it: ```git clone https://github.com/yourname/ignite```
* Create your feature branch: ```git checkout -b my-new-feature```
* Make changes and add them: ```git add .```
* Commit changes: ```git commit -m "Add some feature"```
* Push your commits: ```git push origin my-new-feature```
* Create pull request

# License
[MIT License](https://github.com/go-ignite/ignite/blob/master/LICENSE)
