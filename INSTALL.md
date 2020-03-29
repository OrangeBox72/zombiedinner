INSTALL NOTES
=============

**Installation**

1. install *Golang*
2. set your environment
   * `export GOPATH=~/Documents/go`   (or whereever you would like..)
   * `export GOAPPS=$GOPATH/bin`
   * `export GOROOT=_Whereever you install the GO app_`
   * `export PATH=$PATH:$GOAPPS`
   * `mkdir -p $GOPATH/{src,pkg,bin}`
   * `cd $GOPATH/src`
   * `go get github.com/fatih/color`  (for ANSI color pkg.  {needed} )
   * `go get github.com/OrangeBox72/zombiedice`
   * `cd $GOPATH/src/github.com/OrangeBox72/zombiedice`
   * `go build zombiedice.go`
3. eat **braaains**


