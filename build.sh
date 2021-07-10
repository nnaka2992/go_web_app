#!/bin/sh
if [ "$#" -eq 0 ]; then
    exit 1
elif [ "$#" -eq 1 ]; then
    dir=`git rev-parse --show-prefix`
    program=$1
elif [ "$#" -eq 2 ]; then
    dir=$1
    program=$2
fi

docker-compose run --rm server sh -c "cd $dir; go build -o $program"
