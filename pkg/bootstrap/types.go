package bootstrap

import "context"

type Application interface {
	Run(context.Context) error
}
