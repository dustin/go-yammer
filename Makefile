include $(GOROOT)/src/Make.inc

TARG=github.com/dustin/yammer.go
GOFILES=structs.go yammer.go

include $(GOROOT)/src/Make.pkg
