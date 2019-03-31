package model

type ContentType struct {
	DataID    int
	Published int
	Modified  int
	Title     string
	RemoteID  string
}

//The helpers that generated entity will use to reduce template logic.
