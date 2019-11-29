// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package ha

import (
	"crypto/sha1"
	"encoding/json"
	"github.com/Icinga/icingadb/connection"
	log "github.com/sirupsen/logrus"
)

type Environment struct {
	ID       []byte
	Name     string
	NodeName string
}

// Compute SHA1
func Sha1bytes(bytes []byte) []byte {
	hash := sha1.New()
	hash.Write(bytes)
	return hash.Sum(nil)
}

func IcingaHeartbeatListener(rdb *connection.RDBWrapper, chEnv chan *Environment, chErr chan error) {
	log.Info("Starting heartbeat listener")

	subscription := rdb.Subscribe()
	defer subscription.Close()
	if err := subscription.Subscribe("icinga:stats"); err != nil {
		chErr <- err
		return
	}

	for {
		msg, err := subscription.ReceiveMessage()
		if err != nil {
			chErr <- err
			return
		}

		log.Debug("Got heartbeat")

		var icingaStats struct {
			IcingaApplication struct {
				Status struct {
					IcingaApplication struct {
						App struct {
							Environment string `json:"environment"`
							NodeName    string `json:"node_name"`
						} `json:"app"`
					} `json:"icingaapplication"`
				} `json:"status"`
			} `json:"IcingaApplication"`
		}

		if err = json.Unmarshal([]byte(msg.Payload), &icingaStats); err != nil {
			chErr <- err
			return
		}

		app := &icingaStats.IcingaApplication.Status.IcingaApplication.App
		env := &Environment{Name: app.Environment, ID: Sha1bytes([]byte(app.Environment)), NodeName: app.NodeName}
		chEnv <- env
	}
}
