package defaultserviceregistry

import "github.com/riomhaire/lightauthuserapi/usecases"

type DefaultServiceRegistry struct {
	registry *usecases.Registry
}

func NewDefaultServiceRegistry(registry *usecases.Registry) *DefaultServiceRegistry {
	r := DefaultServiceRegistry{}
	r.registry = registry

	return &r
}

func (r *DefaultServiceRegistry) Register() error {
	return nil

}

func (r *DefaultServiceRegistry) Deregister() error {
	return nil

}
