package repository

type QueryType int

const (
	QueryTypeSelect QueryType = iota
	QueryTypeCount
)
