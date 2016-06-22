#!/bin/bash

go get github.com/laher/goxc
goxc bump
git add .
git commit -m "Updated version."
goxc
