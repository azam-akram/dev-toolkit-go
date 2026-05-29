package main

type Message struct {
	UUID    string `json:"uuid"`
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}
