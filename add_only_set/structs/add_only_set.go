package structs

type AddOnlySet interface {
	Add(v interface{}, r *Result) error
	Show(v interface{}, r *Result) error
}
