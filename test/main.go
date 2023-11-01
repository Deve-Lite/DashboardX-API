package test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/Deve-Lite/DashboardX-API/config"
	"github.com/Deve-Lite/DashboardX-API/internal/application"
	"github.com/Deve-Lite/DashboardX-API/internal/domain"
	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/http/rest"
	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/http/rest/handler"
	"github.com/Deve-Lite/DashboardX-API/internal/interfaces/http/rest/middleware"
	"github.com/Deve-Lite/DashboardX-API/pkg/postgres"
	"github.com/Deve-Lite/DashboardX-API/pkg/redis"
	"github.com/Deve-Lite/DashboardX-API/pkg/smtp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	gredis "github.com/redis/go-redis/v9"
)

type User struct {
	ID           uuid.UUID
	Email        string
	Password     string
	AccessToken  string
	RefreshToken string
}

type test struct {
	c  *config.Config
	d  *sqlx.DB
	ch *gredis.Client
	s  smtp.Client
}

type Test interface {
	SetupApp() (*gin.Engine, *application.Application)
	Teardown()
	CreateUser(app *application.Application, name string, password string, email string) *User
	DeleteUser(app *application.Application, userID uuid.UUID)
	MakeRequest(g *gin.Engine, method string, url string, payload io.Reader, token *string) *httptest.ResponseRecorder
}

func NewTest() Test {
	log.Print("Setup config & database")

	c := config.NewConfig("test.env")

	if postgres.Exists(c.Postgres) {
		postgres.Drop(c.Postgres)
	}

	redis.FlushDB(c.Redis)

	postgres.Create(c.Postgres)
	postgres.RunUp(c.Postgres)

	d := postgres.NewDB(c.Postgres)
	ch := redis.NewDB(c.Redis)

	s := smtp.NewClient(c.SMTP)

	return &test{c, d, ch, s}
}

func (t *test) MakeRequest(g *gin.Engine, method string, url string, payload io.Reader, token *string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, payload)

	if token != nil {
		req.Header.Set("Authorization", *token)
	}

	g.ServeHTTP(w, req)

	return w
}

func (t *test) SetupApp() (*gin.Engine, *application.Application) {
	log.Print("Setup test environment")

	gin.SetMode(gin.ReleaseMode)

	gin := gin.Default()

	app := application.NewApplication(t.c, t.d, t.ch, t.s)

	mRule := middleware.NewRule(app.AuthSrv, app.UserSrv)
	mInfo := middleware.NewInfo(t.c)

	userHnd := handler.NewUserHandler(t.c, app.UserSrv, app.UserMap)
	brokerHnd := handler.NewBrokerHandler(app.BrokerSrv, app.BrokerMap)
	deviceHnd := handler.NewDeviceHandler(app.DeviceSrv, app.ControlSrv, app.DeviceMap, app.ControlMap)

	rest.NewRouter(gin, mRule, mInfo, userHnd, brokerHnd, deviceHnd)

	return gin, app
}

func (t *test) Teardown() {
	log.Print("Teardown test environment")
	t.d.Close()

	postgres.RunDown(t.c.Postgres)
	postgres.Drop(t.c.Postgres)
}

func (t *test) CreateUser(app *application.Application, name string, password string, email string) *User {
	ctx := context.Background()
	defer ctx.Done()

	preUserID, err := app.UserSrv.PreCreate(ctx, &domain.CreateUser{
		Name:     name,
		Password: password,
		Email:    email,
	})
	if err != nil {
		log.Panic(err)
	}

	userID, err := app.UserSrv.Create(ctx, preUserID)
	if err != nil {
		log.Panic(err)
	}

	tokens, err := app.UserSrv.Login(ctx, &domain.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Panic(err)
	}

	return &User{
		ID:           userID,
		Email:        email,
		Password:     password,
		AccessToken:  fmt.Sprintf("Bearer %s", tokens.AccessToken),
		RefreshToken: fmt.Sprintf("Bearer %s", tokens.RefreshToken),
	}
}

func (t *test) DeleteUser(app *application.Application, userID uuid.UUID) {
	ctx := context.Background()
	defer ctx.Done()

	app.UserSrv.Delete(ctx, userID)
}
