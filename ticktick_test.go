package ticktick

import "testing"

func TestNewClient(t *testing.T) {
	client, err := NewClient("username", "password", nil)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", client.getUserInfo())
}
