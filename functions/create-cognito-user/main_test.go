package main

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("success request", func(*testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "sebastian.ferreyra@devmente.com", Password: "PaSsWoRd_100", Name: "sanki0", Case: 0})
		if err != nil {
			t.Fatal("Errora")
		}
		if k != "" {
			t.Fatal("Error")
		}
	})
}

func TestHandler2(t *testing.T) {
	t.Run("success request", func(*testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "sebastian.ferreyra@devmente.com", Password: "PaSsWoRd_100", Username: "44e29d80-0bb0-45a4-9a63-79d4e77d3c4e", Case: 2, ConfirmationCode: "413110"})
		if err != nil {
			t.Fatal("Errora")
		}
		if k != "" {
			t.Fatal("Error")
		}
	})
}

func TestSingIn(t *testing.T) {
	t.Run("success request", func(*testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "sebastian.ferreyra@devmente.com", Password: "PaSsWoRd_104", Username: "44e29d80-0bb0-45a4-9a63-79d4e77d3c4e", Case: 4})
		if err != nil {
			t.Fatal("Errora")
		}
		if k != "" {
			t.Fatal("Error")
		}
	})
}

func TestChangePassword(t *testing.T) {
	t.Run("success request", func(*testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "sebastian.ferreyra@devmente.com", Password: "PaSsWoRd_100", Username: "44e29d80-0bb0-45a4-9a63-79d4e77d3c4e", NewPassword: "PaSsWoRd_102", Case: 5})
		if err != nil {
			t.Fatal("Errora")
		}
		if k != "" {
			t.Fatal("Error")
		}
	})
}

func TestForgotPassword(t *testing.T) {
	t.Run("success request", func(*testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "sebastian.ferreyra@devmente.com", Password: "PaSsWoRd_100", Username: "44e29d80-0bb0-45a4-9a63-79d4e77d3c4e", Case: 6})
		if err != nil {
			t.Fatal("Errora")
		}
		if k != "" {
			t.Fatal("Error")
		}
	})
}

func TestConfirmForgotPassword(t *testing.T) {
	t.Run("success request", func(*testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "sebastian.ferreyra@devmente.com", Password: "PaSsWoRd_104", Username: "44e29d80-0bb0-45a4-9a63-79d4e77d3c4e", Case: 7, ConfirmationCode: "178741"})
		if err != nil {
			t.Fatal("Errora")
		}
		if k != "" {
			t.Fatal("Error")
		}
	})
}
