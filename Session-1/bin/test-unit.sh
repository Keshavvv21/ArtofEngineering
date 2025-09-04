#!/bin/bash

set -e


THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $THIS_DIR/..

ginkgo -r -p race -keepGoing --trace --progress $THIS_DIR/../resources