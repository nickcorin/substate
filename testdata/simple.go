package testdata

//go:generate gensubstate

type Substate interface {
	FooClient() FooClient
	BarClient() BarClient
}

type FooClient interface{}

type BarClient interface{}
