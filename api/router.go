package api

type Server interface {
	Run(address string) error
}
