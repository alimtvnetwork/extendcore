#!/bin/bash

WorkDir="$PWD/.."
BinDir=$WorkDir/bin
Platforms="darwin linux windows"
Architectures="386 amd64"
CliName=main
CliPath=cmd/$CliName/*.go
OutputName="cli"

echo $Platforms
echo $Architectures

echo ""
echo " ---- [Start] deploy for platforms ($Platforms [$Architectures]) [Start]-----"
echo ""
echo "Work dir     : $WorkDir"
echo "Binaries dir : $BinDir"

docker run --rm -it -v "$WorkDir":/usr/src/myapp -v "$GOPATH":/go -w /usr/src/myapp golang:1.17.8 bash -c '
rm -rf bin && \
mkdir bin && \
for CliName in main; do
  for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do

    export GOOS=$GOOS
    export GOARCH=$GOARCH

    if [ $GOOS == "windows" ]
    then
       echo "Building $GOOS-$GOARCH.exe"
       go build -o bin/cli-$GOOS-$GOARCH.exe cmd/$CliName/*.go
    else
       echo "Building $GOOS-$GOARCH"
      go build -o bin/cli-$GOOS-$GOARCH cmd/$CliName/*.go
    fi
    done
  done
done' &&
  echo "Build complete" &&
  echo "" &&
  echo "Permission adding:" &&
  echo "chown -R root:root $BinDir" &&
  echo "chmod -R 777 $BinDir" &&
  echo "" &&
  chown -R root:root "$BinDir" &&
  chmod -R 777 "$BinDir" &&
  echo "" &&
  echo "ls -la $BinDir:" &&
  ls -la "$BinDir" &&
  echo "" &&
  echo "EnvPath" &&
  export PATH=$PATH:"$BinDir" &&
  echo "" &&
  echo $PATH &&
  echo "" &&
  echo "$\"Path\":" &&
  echo "export PATH=\$PATH:\"$BinDir\"" &&
  echo "running : ${BinDir}/cli-linux-amd64" &&
  "$BinDir"/cli-linux-amd64 &&
  echo "" &&
  echo " ---- [End] deploy for all platforms ($Platforms [$Architectures]) [end]-----" &&
  echo ""
