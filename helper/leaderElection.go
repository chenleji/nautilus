package helper

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	logs "github.com/sirupsen/logrus"
	"time"
)

///////////////////////////////////
// role
//////////////////////////////////

func GetRoleInst() *Role {
	if role == nil {
		role = &Role{}
	}
	return role
}

var role *Role

type Role struct {
	isLeader bool
}

func (r *Role) IsLeader() bool {
	return r.isLeader
}

func (r *Role) SetRole(isLeader bool) {
	r.isLeader = isLeader
}

///////////////////////////////////
// leader election
//////////////////////////////////

const (
	LeaderElectionPathFmt = "lock/%s/leader"
)

type LeaderElection struct {
	Consul   *Consul
	TTL      time.Duration
	Callback func(leader bool)
	key      string
}

func (l *LeaderElection) Run() {
	var (
		identity    = Utils{}.GetMyIPAddr()
		retryPeriod = time.Second * 5
		maxAttempt  = 30
		attempt     = 0
		kv          = l.Consul.client.KV()
		session     = l.Consul.client.Session()
	)

	l.key = fmt.Sprintf(LeaderElectionPathFmt, Utils{}.GetAppName())
	callback := func(leader bool) {
		GetRoleInst().SetRole(leader)
		l.Callback(leader)
	}

	se := &consul.SessionEntry{
		Name:      identity,
		TTL:       "10s",
		LockDelay: time.Nanosecond,
	}

	for {
		// check retry times
		if attempt > maxAttempt {
			panic("Run retry times reach Max failed count.")
		}

		// create new session
		sessionId, _, err := session.CreateNoChecks(se, nil)
		if err != nil {

			logs.Info("create session err, retry after 5 second. ", err)
			time.Sleep(retryPeriod)
			attempt ++

			continue
		}

		logs.Info("session sessionId:", sessionId)
		p := &consul.KVPair{
			Key:     l.key,
			Value:   []byte(identity),
			Session: sessionId,
		}

		// try to acquire lock
		locked, _, err := kv.Acquire(p, nil)
		if err != nil {

			logs.Info("acquire err, retry after 5 second. ", err)
			time.Sleep(retryPeriod)
			attempt ++

			continue
		}

		// unlocked
		if !locked {

			callback(false)
			respChan := l.Consul.WatchKey(l.key, nil)

			select {
			case ret := <-respChan:
				if ret.Error != nil {
					logs.Info("watch key err, retry after 5 second. ", err)
					time.Sleep(retryPeriod)

				} else {
					logs.Info("leader released, it's time to election lock!")
				}
			}

			// locked
		} else {

			callback(true)
			session.RenewPeriodic(se.TTL, sessionId, nil, nil)
			//utils.Display("err:", err)

			callback(false)
			//time.Sleep(10 * time.Second)
		}

		attempt = 0
	}
}
