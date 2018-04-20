#! /bin/bash
docker build --build-arg VERSION=`git describe` -t goignite/ignite .
