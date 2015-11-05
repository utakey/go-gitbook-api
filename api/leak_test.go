package api

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"testing"

	"github.com/GitbookIO/go-gitbook-api/client"
)

func TestLeaks(t *testing.T) {
	// Get current count
	c1 := openDescriptors()

	wg := &sync.WaitGroup{}

	// Create one client
	c := client.NewClient(client.ClientOptions{
		Host:     "http://localhost:5000/api/",
		Username: "james",
		Password: "730e0de8-ca53-42f9-9ad3-49c547b0cdc0",
	})

	// Do some work
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			c2 := c.Fork(client.ClientOptions{})
			b := Book{c2}

			_, err := b.Get("james/test")
			if err != nil {
				t.Error(err)
			}

			wg.Done()
		}()

		t.Log(c1, "...", openDescriptors())
	}

	wg.Wait()

	time.Sleep(time.Second)

	// See how many files are open now
	c2 := openDescriptors()

	t.Log(c1, "=>", c2)

	// Check for leak
	if c2 > c1 {
		t.Errorf("Leak: %d to %d descriptors", c1, c2)
	}
}

func openDescriptors() int {
	out, err := lsof(os.Getpid())
	if err != nil {
		return 0
	}
	return bytes.Count(out, []byte("\n"))
}

func lsof(pid int) ([]byte, error) {
	return exec.Command(
		"bash", "-c",
		fmt.Sprintf("lsof -p '%d'", pid),
	).Output()
}
