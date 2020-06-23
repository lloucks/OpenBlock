#!/usr/bin/bash


#Please use this file to set your gopath before building the project!
#It will ensure all the imports don't have errors

OLDPATH=`go env GOPATH`

echo old GOPATH is $OLDPATH

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"


echo Setting GOPATH to: $DIR
echo "Run 'restore' to return the GOPATH to the previous value"


function restore(){
    export GOPATH=$OLDPATH
    echo reset GOPATH to $OLDPATH
}


export GOPATH=$DIR
