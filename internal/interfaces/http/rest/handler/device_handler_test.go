package handler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Deve-Lite/DashboardX-API/test"
	"github.com/go-playground/assert"
	"github.com/google/uuid"
)

func patchControlURL(deviceID, controlID uuid.UUID) string {
	return fmt.Sprintf(`/api/v1/devices/%s/controls/%s`, deviceID.String(), controlID.String())
}

func createControlURL(deviceID uuid.UUID) string {
	return fmt.Sprintf(`/api/v1/devices/%s/controls`, deviceID.String())
}

func TestCreateDeviceControl(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	usr := tt.CreateUser(a, "user1", "test123", "user1@user.com")
	bID := tt.CreateBroker(a, usr.ID)
	dID := tt.CreateDevice(a, usr.ID, bID)

	t.Run("should return 201 when device control has been created", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"type": "state",
				"attributes": {
					"onPayload": "test-1",
					"offPayload": "test-3"
				},
				"canDisplayName": true,
				"canNotifyOnPublish": true,
				"icon": {
					"name": "Home",
					"backgroundColor": "#ff00ff"
				},
				"isAvailable": true,
				"isConfirmationRequired": true,
				"name": "Control",
				"qualityOfService": 0,
				"topic": "test"
			}
		`)

		w := tt.MakeRequest(g, "POST", createControlURL(dID), p, &usr.AccessToken)
		assert.Equal(t, 201, w.Code)
	})

	t.Run("should return 409 when state control already exists", func(t *testing.T) {
		dID2 := tt.CreateDevice(a, usr.ID, bID)

		p := strings.NewReader(`
			{
				"type": "state",
				"attributes": {
					"onPayload": "test-1",
					"offPayload": "test-3"
				},
				"canDisplayName": true,
				"canNotifyOnPublish": true,
				"icon": {
					"name": "Home",
				  	"backgroundColor": "#ff00ff"
				},
				"isAvailable": true,
				"isConfirmationRequired": true,
				"name": "Control",
				"qualityOfService": 0,
				"topic": "test"
			}
		`)

		w := tt.MakeRequest(g, "POST", createControlURL(dID2), p, &usr.AccessToken)
		assert.Equal(t, 201, w.Code)

		p2 := strings.NewReader(`
			{
				"type": "state",
				"attributes": {
					"onPayload": "test-1",
					"offPayload": "test-3"
				},
				"canDisplayName": true,
				"canNotifyOnPublish": true,
				"icon": {
					"name": "Home",
					"backgroundColor": "#ff00ff"
				},
				"isAvailable": true,
				"isConfirmationRequired": true,
				"name": "Control",
				"qualityOfService": 0,
				"topic": "test"
			}
		`)

		w2 := tt.MakeRequest(g, "POST", createControlURL(dID2), p2, &usr.AccessToken)
		assert.Equal(t, 409, w2.Code)
	})
}

func TestUpdateDeviceControl(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	usr := tt.CreateUser(a, "user1", "test123", "user1@user.com")
	bID := tt.CreateBroker(a, usr.ID)
	dID := tt.CreateDevice(a, usr.ID, bID)
	dID2 := tt.CreateDevice(a, usr.ID, bID)
	dcID := tt.CreateDeviceControl(a, usr.ID, dID)

	t.Run("should return 204 when device control has been updated", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"topic": "new-topic"
			}
		`)

		w := tt.MakeRequest(g, "PATCH", patchControlURL(dID, dcID), p, &usr.AccessToken)
		assert.Equal(t, 204, w.Code)
	})

	t.Run("should be able to update existing state control", func(t *testing.T) {
		dcID4 := tt.CreateDeviceControl(a, usr.ID, dID2)
		p := strings.NewReader(`
			{
				"type": "state",
				"attributes": { "onPayload": "test1", "offPayload": "test2" }
			}
		`)

		w := tt.MakeRequest(g, "PATCH", patchControlURL(dID2, dcID4), p, &usr.AccessToken)
		assert.Equal(t, 204, w.Code)

		p2 := strings.NewReader(`
			{
				"icon": {
					"name": "Garage",
					"backgroundColor": "#013311"
				}
			}
		`)

		w2 := tt.MakeRequest(g, "PATCH", patchControlURL(dID2, dcID4), p2, &usr.AccessToken)
		assert.Equal(t, 204, w2.Code)
	})

	t.Run("should return 409 when state control already exists", func(t *testing.T) {
		dcID2 := tt.CreateDeviceControl(a, usr.ID, dID)
		p := strings.NewReader(`
			{
				"type": "state",
				"attributes": { "onPayload": "test1", "offPayload": "test2" }
			}
		`)

		w := tt.MakeRequest(g, "PATCH", patchControlURL(dID, dcID2), p, &usr.AccessToken)
		assert.Equal(t, 204, w.Code)

		dcID3 := tt.CreateDeviceControl(a, usr.ID, dID)
		p2 := strings.NewReader(`
			{
				"type": "state",
				"attributes": { "onPayload": "test1", "offPayload": "test2" }
			}
		`)

		w2 := tt.MakeRequest(g, "PATCH", patchControlURL(dID, dcID3), p2, &usr.AccessToken)
		assert.Equal(t, 409, w2.Code)
	})
}
