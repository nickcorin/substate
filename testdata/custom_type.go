package testdata

//go:generate gensubstate -typeName=ServiceLocator

type ServiceLocator interface {
	FooClient() FooClient
	BarClient() BarClient
}

type FooClient interface{}

type BarClient interface{}
