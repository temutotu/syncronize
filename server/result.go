package main

const (
	FAILD = iota
	SUCCESS
)

type ResultJoin struct {
	result    byte
	ishost    byte
	sessionID byte
}

const (
	GUSET = iota
	HOST
)
