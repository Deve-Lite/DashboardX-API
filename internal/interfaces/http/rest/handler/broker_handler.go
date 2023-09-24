package handler

import (
	"errors"
	"net/http"

	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application/dto"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/application/mapper"
	"github.com/Deve-Lite/DashboardX-API-PoC/internal/domain"
	ae "github.com/Deve-Lite/DashboardX-API-PoC/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BrokerHandler interface {
	Get(ctx *gin.Context)
	List(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type brokerHandler struct {
	bs application.BrokerService
	m  mapper.BrokerMapper
}

func NewBrokerHandler(bs application.BrokerService, m mapper.BrokerMapper) BrokerHandler {
	return &brokerHandler{bs, m}
}

// BrokerGet godoc
//
//	@Summary	Get a broker
//	@Tags		Brokers
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		brokerId	path		string	true	"Broker UUID"
//	@Success	200			{object}	dto.GetBrokerResponse
//	@Failure	400			{object}	errors.HTTPError
//	@Failure	401			{object}	errors.HTTPError
//	@Failure	404			{object}	errors.HTTPError
//	@Failure	500			{object}	errors.HTTPError
//	@Router		/brokers/{brokerId} [get]
func (h *brokerHandler) Get(ctx *gin.Context) {
	var err error
	var brokerID, userID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	brokerID, err = h.getBrokerID(ctx)
	if err != nil {
		return
	}

	var broker *domain.Broker
	broker, err = h.bs.Get(ctx, brokerID, userID)
	if err != nil {
		var code int
		if errors.Is(err, ae.ErrBrokerNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusOK, h.m.ModelToDTO(broker))
}

// BrokerList godoc
//
//	@Summary	List brokers
//	@Tags		Brokers
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		dto.GetBrokerResponse
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/brokers [get]
func (h *brokerHandler) List(ctx *gin.Context) {
	var err error
	var userID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	var brokers []*domain.Broker
	brokers, err = h.bs.List(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	r := []dto.GetBrokerResponse{}

	for _, broker := range brokers {
		r = append(r, *h.m.ModelToDTO(broker))
	}

	ctx.JSON(http.StatusOK, r)
}

// BrokerCreate godoc
//
//	@Summary	Create a broker
//	@Tags		Brokers
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		data	body		dto.CreateBrokerRequest	true	"Create data"
//	@Success	201		{object}	dto.CreateBrokerResponse
//	@Failure	400		{object}	errors.HTTPError
//	@Failure	401		{object}	errors.HTTPError
//	@Failure	404		{object}	errors.HTTPError
//	@Failure	500		{object}	errors.HTTPError
//	@Router		/brokers [post]
func (h *brokerHandler) Create(ctx *gin.Context) {
	var err error
	var userID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	body := &dto.CreateBrokerRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	broker := h.m.CreateDTOToCreateModel(body)
	broker.UserID = userID

	var brokerID uuid.UUID
	brokerID, err = h.bs.Create(ctx, broker)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.CreateBrokerResponse{
		ID: brokerID,
	})
}

// BrokerUpdate godoc
//
//	@Summary	Update a broker
//	@Tags		Brokers
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		brokerId	path	string					true	"Broker UUID"
//	@Param		data		body	dto.UpdateBrokerRequest	true	"Update data"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/brokers/{brokerId} [patch]
func (h *brokerHandler) Update(ctx *gin.Context) {
	var err error
	var userID, brokerID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	brokerID, err = h.getBrokerID(ctx)
	if err != nil {
		return
	}

	body := &dto.UpdateBrokerRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	broker := h.m.UpdateDTOToUpdateModel(body)
	broker.UserID = userID
	broker.ID = brokerID

	err = h.bs.Update(ctx, broker)
	if err != nil {
		var code int
		if errors.Is(err, ae.ErrBrokerNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// BrokerDelete godoc
//
//	@Summary	Delete a broker
//	@Tags		Brokers
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		brokerId	path	string	true	"Broker UUID"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/brokers/{brokerId} [delete]
func (h *brokerHandler) Delete(ctx *gin.Context) {
	var err error
	var userID, brokerID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	brokerID, err = h.getBrokerID(ctx)
	if err != nil {
		return
	}

	err = h.bs.Delete(ctx, brokerID, userID)
	if err != nil {
		var code int
		if errors.Is(err, ae.ErrBrokerNotFound) {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}

		ctx.AbortWithStatusJSON(code, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *brokerHandler) getBrokerID(ctx *gin.Context) (uuid.UUID, error) {
	params := &dto.BrokerParams{}

	err := ctx.BindUri(params)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return uuid.Nil, err
	}

	var brokerID uuid.UUID
	brokerID, err = uuid.Parse(params.BrokerID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return uuid.Nil, err
	}

	return brokerID, nil
}

func (h *brokerHandler) getUserID(ctx *gin.Context) (uuid.UUID, error) {
	userID, err := uuid.Parse(ctx.MustGet("UserID").(string))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.HTTPError{Message: "could not retrieve user info"})
		return uuid.Nil, err
	}

	return userID, nil
}
