package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	n "github.com/Deve-Lite/DashboardX-API/pkg/nullable"
	"github.com/Deve-Lite/DashboardX-API/test"
	"github.com/go-playground/assert"
	"github.com/google/uuid"
)

type BrokerCredentials struct {
	ID       string  `json:"id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func patchURL(brokerID string) string {
	return fmt.Sprintf(`/api/v1/brokers/%s`, brokerID)
}

func credentialsURL(brokerID string) string {
	return fmt.Sprintf(`/api/v1/brokers/%s/credentials`, brokerID)
}

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

		w := tt.MakeRequest(g, "POST", "/api/v1/brokers", p, &u.AccessToken)
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

		w := tt.MakeRequest(g, "POST", "/api/v1/brokers", p, &u.AccessToken)
		assert.Equal(t, 201, w.Code)

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

		w = tt.MakeRequest(g, "POST", "/api/v1/brokers", p, &u.AccessToken)
		assert.Equal(t, 409, w.Code)
	})
}

func TestBrokerUpdate(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	u := tt.CreateUser(a, "user1", "test123", "user1@user.com")

	ctx := context.Background()

	bid, _ := a.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              u.ID,
		Name:                "broker-test",
		Server:              "some.server.com",
		Port:                1000,
		KeepAlive:           10,
		IconName:            "icon.jpg",
		IconBackgroundColor: "#fafefa",
		IsSSL:               false,
	})

	bid2, _ := a.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              u.ID,
		Name:                "broker-test",
		Server:              "some-2.server.com",
		Port:                1000,
		KeepAlive:           10,
		IconName:            "icon.jpg",
		IconBackgroundColor: "#fafefa",
		IsSSL:               false,
	})

	ctx.Done()

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

		w := tt.MakeRequest(g, "PATCH", patchURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 204, w.Code)
	})

	t.Run("should return 409 when a broker with provided server already exists", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"server": "server-exists.hivemq.com"
			}
		`)

		w := tt.MakeRequest(g, "PATCH", patchURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 204, w.Code)

		p = strings.NewReader(`
			{
				"server": "server-exists.hivemq.com"
			}
		`)

		w = tt.MakeRequest(g, "PATCH", patchURL(bid2.String()), p, &u.AccessToken)
		assert.Equal(t, 409, w.Code)
	})
}

func TestBrokerSetCredentials(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	dp := strings.NewReader(`
		{
			"username": "def-user",
			"password": "def-pass"
		}
	`)

	u := tt.CreateUser(a, "user1", "test123", "user1@user.com")

	ctx := context.Background()

	bid, _ := a.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              u.ID,
		Name:                "broker-test",
		Server:              "some.server.com",
		Port:                1000,
		KeepAlive:           10,
		IconName:            "icon.jpg",
		IconBackgroundColor: "#fafefa",
		IsSSL:               true,
	})

	ctx.Done()

	t.Run("should return 204 when successfully set up credentials", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"username": "user-123",
				"password": "secret123!"
			}
		`)

		w := tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 204, w.Code)
	})

	t.Run("should return 204 when successfully unset credentials", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"username": null,
				"password": null
			}
		`)

		w := tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 204, w.Code)
	})

	t.Run("should save new credentials", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"username": "test-user-123",
				"password": "someWeakPass"
			}
		`)

		w := tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 204, w.Code)

		w = tt.MakeRequest(g, "GET", credentialsURL(bid.String()), nil, &u.AccessToken)
		data := BrokerCredentials{}
		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			panic(err)
		}

		assert.Equal(t, data.ID, bid.String())
		assert.Equal(t, data.Username, "test-user-123")
		assert.Equal(t, data.Password, "someWeakPass")
	})

	t.Run("should unset credentials when both are set to null", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"username": null,
				"password": null
			}
		`)

		w := tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 204, w.Code)

		w = tt.MakeRequest(g, "GET", credentialsURL(bid.String()), nil, &u.AccessToken)
		data := BrokerCredentials{}
		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			panic(err)
		}

		assert.Equal(t, data.ID, bid.String())
		assert.Equal(t, data.Username, nil)
		assert.Equal(t, data.Password, nil)
	})

	t.Run("should return 400 when brokerId is invalid", func(t *testing.T) {
		w := tt.MakeRequest(g, "PUT", credentialsURL("invalid"), dp, &u.AccessToken)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 400 when credentials are missing", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"username": "test"
			}
		`)

		w := tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 400, w.Code)

		p = strings.NewReader(`
			{
				"password": "test"
			}
		`)

		w = tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 400, w.Code)

		p = strings.NewReader(`{}`)

		w = tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 400 when credentials are invalid", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"username": "",
				"password": ""
			}
		`)

		w := tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 400, w.Code)

		p = strings.NewReader(`
			{
				"username": 123,
				"password": false
			}
		`)

		w = tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), p, &u.AccessToken)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 401 when token is missing", func(t *testing.T) {
		w := tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), dp, nil)
		assert.Equal(t, 401, w.Code)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		tk := "invalid-token"
		w := tt.MakeRequest(g, "PUT", credentialsURL(bid.String()), dp, &tk)
		assert.Equal(t, 401, w.Code)
	})

	t.Run("should return 404 when broker does not exist", func(t *testing.T) {
		w := tt.MakeRequest(g, "PUT", credentialsURL(uuid.New().String()), dp, &u.AccessToken)
		assert.Equal(t, 404, w.Code)
	})
}

func TestBrokerGetCredentials(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	u := tt.CreateUser(a, "user1", "test123", "user1@user.com")

	ctx := context.Background()

	bid, _ := a.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              u.ID,
		Name:                "broker-test",
		Server:              "some.server.com",
		Port:                1000,
		KeepAlive:           10,
		IconName:            "icon.jpg",
		IconBackgroundColor: "#fafefa",
		IsSSL:               true,
	})

	a.BrokerSrv.SetCredentials(ctx, &domain.UpdateBroker{
		UserID:   u.ID,
		ID:       bid,
		Password: n.NewString("secretPass123", false, true),
		Username: n.NewString("user1", false, true),
	})

	a.BrokerSrv.Create(ctx, &domain.CreateBroker{
		UserID:              u.ID,
		Name:                "broker-test-2",
		Server:              "some-2.server.com",
		Port:                1000,
		KeepAlive:           10,
		IconName:            "icon.jpg",
		IconBackgroundColor: "#fafefa",
		IsSSL:               true,
	})

	ctx.Done()

	t.Run("should return 200 when successfully fetched credentials", func(t *testing.T) {
		w := tt.MakeRequest(g, "GET", credentialsURL(bid.String()), nil, &u.AccessToken)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("should fetch credentials", func(t *testing.T) {
		w := tt.MakeRequest(g, "GET", credentialsURL(bid.String()), nil, &u.AccessToken)
		assert.Equal(t, 200, w.Code)

		data := BrokerCredentials{}
		if err := json.Unmarshal(w.Body.Bytes(), &data); err != nil {
			panic(err)
		}

		assert.Equal(t, data.ID, bid.String())
		assert.Equal(t, data.Username, "user1")
		assert.Equal(t, data.Password, "secretPass123")
	})

	t.Run("should return 400 when brokerId is invalid", func(t *testing.T) {
		w := tt.MakeRequest(g, "GET", credentialsURL("invalid"), nil, &u.AccessToken)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 401 when token is missing", func(t *testing.T) {
		w := tt.MakeRequest(g, "GET", credentialsURL(bid.String()), nil, nil)
		assert.Equal(t, 401, w.Code)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		tk := "invalid-token"
		w := tt.MakeRequest(g, "GET", credentialsURL(bid.String()), nil, &tk)
		assert.Equal(t, 401, w.Code)
	})

	t.Run("should return 404 when broker does not exist", func(t *testing.T) {
		w := tt.MakeRequest(g, "GET", credentialsURL(uuid.New().String()), nil, &u.AccessToken)
		assert.Equal(t, 404, w.Code)
	})
}
