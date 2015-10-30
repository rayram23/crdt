package structs

type AddOnlySet interface {
	Add(v string, r *Result) error
	Show(v string, r *Result) error
}
