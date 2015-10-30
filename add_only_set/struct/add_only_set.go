package structs

type AddOnlySet interface {
	Add(v interface{}, r *string) error
	Show(v interface{}, r []interface{}) error
}
type AddOnlyImpl struct {
	vals map[interface{}]interface{}
}

var _ AddOnlySet = &AddOnlyImpl{}

func (a *AddOnlyImpl) Add(v interface{}, r *string) error {
	fmt.Print("Got to add\n")
	return nil
}
func (a *AddOnlyImpl) Show(v interface{}, r []interface{}) error {
	fmt.Print("Got to show\n")
	return nil
}
