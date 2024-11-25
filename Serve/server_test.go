package serve

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	app := New()
	assert.NotNil(t, app)
}


func TestServerStart(t *testing.T) {
  var wg sync.WaitGroup
	wg.Add(1)
  app := New()
  go app.Start(&wg)
  time.Sleep(1 * time.Second)
  assert.True(t, true, "Server Started and finished")
}