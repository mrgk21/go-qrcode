#!/bin/bash

if [ -z "$1" ]; then
    echo "Available commands: install, test, coverage"
    exit 1
fi

# Main logic
case "$1" in
  test)
    cd app_test/
    go test -v
    cd ..
    ;;

  coverage)
    go test -v -coverprofile temp.out && go tool cover -html=temp.out -o temp.html
    echo "open temp.html to see the code coverage"
    ;;

  install)
    go install github.com/vektra/mockery/v3@v3.6.1
    ;;

  mock)
    # generate mocks
    mockery
    ;;

  *)
    echo "❌ Unknown command: $1"
    echo "Available commands: test, coverage, install, mock"
    exit 1
    ;;
esac
