package handler_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Deve-Lite/DashboardX-API-PoC/test"
	"github.com/go-playground/assert"
)

func TestBrokerCreate(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	u := tt.CreateUser(a, "user1", "test123", "user1@user.com")

	t.Run("should return 201 when broker has been created", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"name": "Test Devices",
				"server": "broker.hivemq.com",
				"port": 8884,
				"keepAlive": 60,
				"icon": {
				  "name": "Home",
				  "backgroundColor": "#ff00ff"
				},
				"isSsl": true,
				"clientId": "123"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/brokers", p)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	})

	t.Run("should return 409 when a broker with provided server already exists", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"name": "Test Devices",
				"server": "server-exists.com",
				"port": 8884,
				"keepAlive": 60,
				"icon": {
				  "name": "Home",
				  "backgroundColor": "#ff00ff"
				},
				"isSsl": true
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/brokers", p)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		p = strings.NewReader(`
			{
				"name": "Test Devices 2",
				"server": "server-exists.com",
				"port": 8884,
				"keepAlive": 60,
				"icon": {
				  "name": "Home",
				  "backgroundColor": "#ff00ff"
				},
				"isSsl": true
			}
		`)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/brokers", p)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		log.Print()

		assert.Equal(t, 409, w.Code)
	})
}

func TestBrokerUpdate(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	u := tt.CreateUser(a, "user1", "test123", "user1@user.com")

	p := strings.NewReader(`
		{
			"name": "Test Devices",
			"server": "broker-new-2.hivemq.com",
			"port": 8884,
			"keepAlive": 60,
			"icon": {
			"name": "Home",
			"backgroundColor": "#ff00ff"
			},
			"isSsl": true,
			"clientId": "123"
		}
	`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/brokers", p)
	req.Header.Set("Authorization", u.AccessToken)
	g.ServeHTTP(w, req)

	t.Run("should return 204 when broker has been updated", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"name": "Test Devices",
				"server": "broker.hivemq.com",
				"port": 8884,
				"keepAlive": 60,
				"icon": {
				  "name": "Home",
				  "backgroundColor": "#ff00ff"
				},
				"isSsl": true,
				"clientId": "123"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/brokers", p)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
	})

	t.Run("should return 409 when a broker with provided server already exists", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"name": "Test Devices",
				"server": "server-exists.com",
				"port": 8884,
				"keepAlive": 60,
				"icon": {
				  "name": "Home",
				  "backgroundColor": "#ff00ff"
				},
				"isSsl": true
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/brokers", p)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		p = strings.NewReader(`
			{
				"name": "Test Devices 2",
				"server": "server-exists.com",
				"port": 8884,
				"keepAlive": 60,
				"icon": {
				  "name": "Home",
				  "backgroundColor": "#ff00ff"
				},
				"isSsl": true
			}
		`)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/brokers", p)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		log.Print()

		assert.Equal(t, 409, w.Code)
	})
}
