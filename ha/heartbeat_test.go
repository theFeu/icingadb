// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package ha

import (
	"github.com/Icinga/icingadb/config/testbackends"
	"github.com/Icinga/icingadb/connection"
	"github.com/intel-go/fastjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var icingastate = "{\"IcingaApplication\":" +
	"{\"status\": " +
	"{\"icingaapplication\":" +
	"{\"app\":{" +
	"\"environment\": \"\"," +
	"\"node_name\": \"master1.icinga.test.com\"" +
	"}}}}, \"config_dump_in_progress\": false}"

func TestIcingaHeartbeatListener(t *testing.T) {
	rdb := connection.NewRDBWrapper(testbackends.RedisTestAddr, 64)
	assert.True(t, rdb.CheckConnection(false), "This test needs a working Redis connection")

	chEnv := make(chan *Environment)

	go func() {
		chErr := make(chan error)
		IcingaHeartbeatListener(rdb, chEnv, chErr)
		require.NoError(t, <-chErr, "redis connection error")
	}()

	time.Sleep(time.Second * 2)

	var uj interface{} = nil
	if err := fastjson.Unmarshal([]byte(icingastate), &uj); err != nil {
		assert.Nil(t, err)
	}

	rdb.Rdb.Publish("icinga:stats", icingastate)

	env := <-chEnv

	assert.NotNil(t, env.ID, "no valid env received")
}
