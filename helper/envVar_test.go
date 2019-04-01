package helper

import (
	"os"
	"testing"
)

func TestEnvVar_GetConsulURL(t *testing.T) {
	oldEnv := os.Getenv("CONSUL_ADDR")
	os.Setenv("CONSUL_ADDR", "localhost")

	url := EnvVar{}.GetConsulAddr()
	if url != "localhost" {
		t.Error("get consul url failed!")
	}

	os.Setenv("CONSUL_ADDR", oldEnv)
}

func TestEnvVar_GetConsulPort(t *testing.T) {
	os.Setenv("CONSUL_PORT", "8500")

	port := EnvVar{}.GetConsulPort()
	if port != "8500" {
		t.Error("get consul port failed!")
	}

}
