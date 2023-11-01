package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Deve-Lite/DashboardX-API/test"
	"github.com/go-playground/assert"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Theme    string `json:"theme"`
	Language string `json:"language"`
}

func TestUserRegister(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, _ := tt.SetupApp()

	t.Run("should return 202 when a user has been registered", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"name": "new-user",
				"email": "new-user@test.com",
				"password": "secretpassword"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/register", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 202, w.Code)
	})

	t.Run("should return 400 when password is invalid", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"name": "new-user",
				"email": "new-user@test.com",
				"password": "123"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/register", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 400 when email is invalid", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"name": "new-user",
				"email": "not-an-email",
				"password": "test1234"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/register", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 400 when name is invalid", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"name": "1",
				"email": "new-user@test.com",
				"password": "test1234"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/register", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 409 when email already exists", func(t *testing.T) {
		p1 := strings.NewReader(`
			{
				"name": "new-user-2",
				"email": "new-user-2@test.com",
				"password": "test1234"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/register", p1)
		g.ServeHTTP(w, req)

		assert.Equal(t, 202, w.Code)

		p2 := strings.NewReader(`
			{
				"name": "new-user-3",
				"email": "new-user-2@test.com",
				"password": "test1234"
			}
		`)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/users/register", p2)
		g.ServeHTTP(w, req)

		assert.Equal(t, 409, w.Code)
	})
}

func TestUserLogin(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	tt.CreateUser(a, "login-user", "secretpassword", "new-login@test.com")

	t.Run("should return 200 when a user has been logged in", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"email": "new-login@test.com",
				"password": "secretpassword"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("should return tokens pair when user has been logged in", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"email": "new-login@test.com",
				"password": "secretpassword"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		data := Tokens{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			panic(err)
		}

		assert.NotEqual(t, data.AccessToken, "")
		assert.NotEqual(t, data.RefreshToken, "")
	})

	t.Run("should return tokens with valid payload", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"email": "new-login@test.com",
				"password": "secretpassword"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		data := Tokens{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			panic(err)
		}

		var user1ID, user2ID uuid.UUID
		var isAdmin1, isAdmin2 bool

		token1, _, err := jwt.NewParser().ParseUnverified(data.AccessToken, jwt.MapClaims{})
		if err != nil {
			panic(err)
		}

		if claims, ok := token1.Claims.(jwt.MapClaims); ok {
			user1ID = uuid.MustParse(claims["sub"].(string))
			isAdmin1 = claims["is_admin"].(bool)
		}

		user, err := a.UserSrv.Get(context.Background(), user1ID)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, user.ID, user1ID)
		assert.Equal(t, user.IsAdmin, isAdmin1)

		token2, _, err := jwt.NewParser().ParseUnverified(data.RefreshToken, jwt.MapClaims{})
		if err != nil {
			panic(err)
		}

		if claims, ok := token2.Claims.(jwt.MapClaims); ok {
			user2ID = uuid.MustParse(claims["sub"].(string))
			isAdmin2 = claims["is_admin"].(bool)
		}

		user, err = a.UserSrv.Get(context.Background(), user2ID)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, user.ID, user2ID)
		assert.Equal(t, user.IsAdmin, isAdmin2)

		assert.Equal(t, user1ID, user2ID)
		assert.Equal(t, isAdmin1, isAdmin2)
	})

	t.Run("should return 400 when password is invalid", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"email": "new-login@test.com",
				"password": "123"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 400 when password is incorrect", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"email": "new-login@test.com",
				"password": "incorrect"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 400 when email is invalid", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"email": "no-email",
				"password": "123"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("should return 404 when user does not exist", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"email": "no-user@test.com",
				"password": "test1234"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/login", p)
		g.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}

func TestUserTokensMe(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	t.Run("should return 200 when new tokens have been created", func(t *testing.T) {
		u := tt.CreateUser(a, "user1", "test1234", "user1@test.com")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.RefreshToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("should return tokens pair when tokens have been refreshed", func(t *testing.T) {
		u := tt.CreateUser(a, "user2", "test1234", "user2@test.com")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.RefreshToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		data := Tokens{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			panic(err)
		}

		assert.NotEqual(t, data.AccessToken, "")
		assert.NotEqual(t, data.RefreshToken, "")
	})

	t.Run("should return tokens with valid payload", func(t *testing.T) {
		u := tt.CreateUser(a, "user3", "test1234", "user3@test.com")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.RefreshToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		data := Tokens{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			panic(err)
		}

		var user1ID, user2ID uuid.UUID
		var isAdmin1, isAdmin2 bool

		token1, _, err := jwt.NewParser().ParseUnverified(data.AccessToken, jwt.MapClaims{})
		if err != nil {
			panic(err)
		}

		if claims, ok := token1.Claims.(jwt.MapClaims); ok {
			user1ID = uuid.MustParse(claims["sub"].(string))
			isAdmin1 = claims["is_admin"].(bool)
		}

		user, err := a.UserSrv.Get(context.Background(), user1ID)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, user.ID, user1ID)
		assert.Equal(t, user.IsAdmin, isAdmin1)

		token2, _, err := jwt.NewParser().ParseUnverified(data.RefreshToken, jwt.MapClaims{})
		if err != nil {
			panic(err)
		}

		if claims, ok := token2.Claims.(jwt.MapClaims); ok {
			user2ID = uuid.MustParse(claims["sub"].(string))
			isAdmin2 = claims["is_admin"].(bool)
		}

		user, err = a.UserSrv.Get(context.Background(), user2ID)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, user.ID, user2ID)
		assert.Equal(t, user.IsAdmin, isAdmin2)

		assert.Equal(t, user1ID, user2ID)
		assert.Equal(t, isAdmin1, isAdmin2)
	})

	t.Run("should accept newly created refresh token", func(t *testing.T) {
		u := tt.CreateUser(a, "user8", "test1234", "user8@test.com")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.RefreshToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		data := Tokens{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			panic(err)
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.RefreshToken))
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("should create valid access token", func(t *testing.T) {
		u := tt.CreateUser(a, "user10", "test1234", "user10@test.com")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.RefreshToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		data := Tokens{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			panic(err)
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.AccessToken))
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("should not accept an access token", func(t *testing.T) {
		u := tt.CreateUser(a, "user4", "test1234", "user4@test.com")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("should allow a refresh token to be used only once", func(t *testing.T) {
		u := tt.CreateUser(a, "user6", "test1234", "user6@test.com")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.RefreshToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.RefreshToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 401)
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", "invalid")
		g.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", "Bearer invalid")
		g.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("should return 401 when token is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		g.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("should return 404 when user does not exist", func(t *testing.T) {
		u := tt.CreateUser(a, "user5", "test1234", "user5@test.com")
		tt.DeleteUser(a, u.ID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/me/tokens", nil)
		req.Header.Set("Authorization", u.RefreshToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}

func TestUserGetMe(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	u1 := tt.CreateUser(a, "user1", "test1234", "user1@test.com")

	t.Run("should return 200 when fetched a user", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", u1.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("should return data in valid format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", u1.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		data := User{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, data.ID, u1.ID.String())
		assert.Equal(t, data.Email, "user1@test.com")
		assert.Equal(t, data.Name, "user1")
		assert.Equal(t, data.Language, "pl")
		assert.Equal(t, data.Theme, "inherit")
	})

	t.Run("should return 401 when token is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", "invalid")
		g.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("should return 401 when token is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
		g.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})

	t.Run("should return 404 when user does not exist", func(t *testing.T) {
		u := tt.CreateUser(a, "user5", "test1234", "user5@test.com")
		tt.DeleteUser(a, u.ID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})
}

func TestUserUpdateMe(t *testing.T) {
	tt := test.NewTest()
	defer tt.Teardown()
	g, a := tt.SetupApp()

	u := tt.CreateUser(a, "user1", "test1234", "user1@test.com")

	t.Run("should return 204 when updated a user", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"language": "en"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/users/me", p)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
	})

	t.Run("should save new values", func(t *testing.T) {
		p := strings.NewReader(`
			{
				"language": "it",
				"email": "new-mail@test.pl",
				"name": "new-name",
				"theme": "dark"
			}
		`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/v1/users/me", p)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("Authorization", u.AccessToken)
		g.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		data := User{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, data.Email, "new-mail@test.pl")
		assert.Equal(t, data.Name, "new-name")
		assert.Equal(t, data.Language, "it")
		assert.Equal(t, data.Theme, "dark")
	})
}
