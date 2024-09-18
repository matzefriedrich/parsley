package bootstrap

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RunParsleyApplication(t *testing.T) {
	// Arrange
	var run bool
	runFunc := func(ctx context.Context) error {
		run = true
		return nil
	}

	appFactory := func(resolver types.Resolver) bootstrap.Application {
		return &testApp{
			RunFunc: runFunc,
		}
	}
	// Act
	err := bootstrap.RunParsleyApplication(context.Background(), appFactory, func(registry types.ServiceRegistry) error {
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.True(t, run)
	assert.NotErrorIs(t, err, bootstrap.ErrCannotRegisterAppFactory)
}

type testApp struct {
	RunFunc ApplicationRunFunc
}

type ApplicationRunFunc func(ctx context.Context) error

func (t testApp) Run(ctx context.Context) error {
	return t.RunFunc(ctx)
}

var _ bootstrap.Application = (*testApp)(nil)
