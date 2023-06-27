package controllers

// IRegisterable is interface to perform on-demand registrations.
type IRegisterable interface {
	// Register perform required registration steps.
	Register()
}
