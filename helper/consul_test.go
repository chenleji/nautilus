package helper

import (
	"testing"
	"time"
)

func TestRegistryConsul(t *testing.T) {
	consul := Consul{}.New()
	err := consul.RegistryService("nautilus", "8011", "/health")
	if err != nil {
		t.Error(err)
	}

	t.Log("registry consul success!")
}

func TestWatchKey(t *testing.T) {
	key := "test/nautilus/key"

	consul := Consul{}.New()

	if ! consul.CheckKey(key, nil) {
		t.Log("key not exist, create it.")

		if err := consul.SetKey(key, "1"); err != nil {
			t.Error(err)
		}
	}

	stopChan := make(chan bool)
	respCh := consul.WatchKey(key, stopChan)
	go func() {
		for {
			select {
			case <-stopChan:
				return

			case ret := <-respCh:
				if ret.Error != nil {
					t.Error(ret.Error)
				} else {
					t.Log("key changed to :", string(ret.Value))
					if (string)(ret.Value) != "2" {
						t.Error("updated key is err")
					}
				}
			}
		}
	}()

	time.Sleep(5 * time.Second)
	consul.SetKey(key, "2")
	time.Sleep(5 * time.Second)

	stopChan <- true

}
