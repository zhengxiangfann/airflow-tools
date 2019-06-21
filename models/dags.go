package models

type Label struct {
	Label string `json:"label"`
}

type Node struct {
	Id string `json:"id"`
	Value Label `json:"value"`
}
type Link struct {
	U string  `json:"u"`
	V string   `json:"v"`
	Value Label `json:"value"`
}

type DAG struct {
	Name string `json:"name"`
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
 }



