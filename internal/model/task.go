package model

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type IdStrusct struct {
	Id int64 `json:"id"`
}

type passStruct struct {
	Password string `json:"password"`
}
type token struct {
	Token string `json:"token"`
}
