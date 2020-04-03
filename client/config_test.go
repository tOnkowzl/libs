package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigReturnHasBasicAuthWhenHasUsernameAndPassword(t *testing.T) {
	b := BasicAuth{
		Username: "username",
		Password: "password",
	}

	assert.Equal(t, b.HasBasicAuth(), true)
}

func TestConfigReturnNotHasBasicAuthWhenNotHasUsername(t *testing.T) {
	b := BasicAuth{
		Username: "",
		Password: "password",
	}

	assert.Equal(t, b.HasBasicAuth(), false)
}

func TestConfigReturnNotHasBasicAuthWhenNotHasPassword(t *testing.T) {
	b := BasicAuth{
		Username: "username",
		Password: "",
	}

	assert.Equal(t, b.HasBasicAuth(), false)
}
