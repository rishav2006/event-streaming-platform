package models

type EventModel struct {
	Offset  int    `json:"offset"`
	Message string `json:"message"`
}

var Count int = 2;

var Events = []EventModel{
	{0, "Hello"},
	{1,"World"},
}
