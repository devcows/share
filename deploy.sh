#!/bin/bash

go get github.com/laher/goxc
goxc bump
git add .goxc.json
git commit -m "Updated version."
goxc
