package sshclient_test

import (
	"context"
	"encoding/json"

	"htdvisser.dev/exp/ssh/sshclient"
)

func Example() {
	configJSON := `{
  "address": "localhost:2222",
  "host_key": {
    "source": "known_hosts",
    "known_hosts": { "file": "testdata/known_hosts" }
  },
  "username": "testuser",
  "auth_methods": [
    {
      "method": "private_keys",
      "private_keys": [
        { "file": "testdata/id_ed25519" },
        { "file": "testdata/id_ecdsa" },
        { "file": "testdata/id_rsa" }
      ]
    },
    {
      "method": "password",
      "password": "testpassword"
    }
  ]
}`
	var config sshclient.ConnectConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		// handle error
		return
	}

	_, client, err := config.Dial(context.Background())
	if err != nil {
		// handle error
		return
	}

	client.Close()
}
