package helper

import (
	logs "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestLeaderElection_RunElection(t *testing.T) {
	le := &LeaderElection{
		Consul: Consul{}.New(),
		TTL:    time.Second,
		Callback: func(leader bool) {
			logs.Info("unit test == role is leader:", leader)
		},
	}

	go le.Run()

	time.Sleep(10 * time.Second)

	t.Log(GetRoleInst().IsLeader())

}
