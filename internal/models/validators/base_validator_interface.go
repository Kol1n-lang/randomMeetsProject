package validators

type ValidatorUserInterface interface {
	Validate() error
}
