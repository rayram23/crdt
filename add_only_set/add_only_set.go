package main

type AddOnlySet interface {
	Add(v interface{}, r *string) error
	Show(v interface{}, r []interface{}) error
}
