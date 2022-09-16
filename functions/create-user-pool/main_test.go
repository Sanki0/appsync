package main

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("success request", func(*testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "sebastian.ferreyra@devmente.com", Password: "PaSsWoRd_100", Name: "sanki0", Case: 1})
		if err != nil {
			t.Fatal("Errora")
		}
		if k != "" {
			t.Fatal("Error")
		}
	})
}
