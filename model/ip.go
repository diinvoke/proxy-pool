package model

type IP struct {
	Address  string `json: "address"`
	Port     string `json: "port"`
	Protocol string `json: "protocal"`
}
