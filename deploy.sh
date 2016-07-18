#!/bin/bash

go get github.com/laher/goxc
git checkout master
goxc bump
git add .goxc.json
git commit -m "Updated version."
git push
goxc
