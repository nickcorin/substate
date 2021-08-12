package testdata

import "database/sql"

//go:generate gensubstate

type Substate interface {
	Database() *sql.DB
	FooClient() FooClient
}

type Foo struct{}

type FooClient interface {
	Lookup() (*Foo, error)
}
