package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiErr(t *testing.T) {
	multierr := &MultiErr{}
	require.True(t, multierr.IsNil())
	require.Nil(t, multierr.Err())

	multierr.Append(errors.New("I am err0"))
	multierr.Append(errors.New("I am err1"))
	require.NotNil(t, multierr.Err())
	require.Equal(t, "multierr: [``I am err0``  ``I am err1``]", multierr.Error())
}
