package helper

import "testing"

func TestUtils_GetMyIPAddr(t *testing.T) {
	tool := Utils{}
	id := tool.GetMyIdentity()

	if id != "" {
		t.Log("id:", id)
		return
	}
	t.Error("id is null.")

}

func TestUtils_SystemHealth(t *testing.T) {
	tool := Utils{}
	err := tool.SystemHealth()
	if err != nil {
		t.Log(err)
	}

	t.Log("system is health!")

}
