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

type DeviceHandler interface {
	Get(ctx *gin.Context)
	List(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	ListControls(ctx *gin.Context)
	CreateControl(ctx *gin.Context)
	UpdateControl(ctx *gin.Context)
	DeleteControl(ctx *gin.Context)
}

type deviceHandler struct {
	ds  application.DeviceService
	dcs application.DeviceControlService
	dm  mapper.DeviceMapper
	dcm mapper.DeviceControlMapper
}

func NewDeviceHandler(ds application.DeviceService, dcs application.DeviceControlService, dm mapper.DeviceMapper, dcm mapper.DeviceControlMapper) DeviceHandler {
	return &deviceHandler{ds, dcs, dm, dcm}
}

// DeviceGet godoc
//
//	@Summary	Get a single device
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		deviceId	path		string	true	"Device UUID"
//	@Success	200			{object}	dto.GetDeviceResponse
//	@Failure	400			{object}	errors.HTTPError
//	@Failure	401			{object}	errors.HTTPError
//	@Failure	404			{object}	errors.HTTPError
//	@Failure	500			{object}	errors.HTTPError
//	@Router		/devices/{deviceId} [get]
func (h *deviceHandler) Get(ctx *gin.Context) {
	var err error
	var deviceID, userID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	deviceID, err = h.getDeviceID(ctx)
	if err != nil {
		return
	}

	var device *domain.Device
	device, err = h.ds.Get(ctx, deviceID, userID)
	if err != nil {
		if errors.Is(err, ae.ErrDeviceNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusOK, h.dm.ModelToDTO(device))
}

// DeviceList godoc
//
//	@Summary	List devices
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		brokerId	query		string	false	"Broker UUID"	Format(UUID)
//	@Success	200			{array}		dto.GetDeviceResponse
//	@Failure	401			{object}	errors.HTTPError
//	@Failure	500			{object}	errors.HTTPError
//	@Router		/devices [get]
func (h *deviceHandler) List(ctx *gin.Context) {
	var err error
	var userID uuid.UUID
	var brokerID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	filters := &domain.ListDeviceFilters{
		UserID: userID,
	}

	query := &dto.DeviceQuery{}
	err = ctx.ShouldBindQuery(query)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	if query.BrokerID != nil {
		brokerID, err = uuid.Parse(*query.BrokerID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
			return
		}

		filters.BrokerID = uuid.NullUUID{UUID: brokerID, Valid: true}
	}

	var devices []*domain.Device
	devices, err = h.ds.List(ctx, filters)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	r := []dto.GetDeviceResponse{}

	for _, device := range devices {
		r = append(r, *h.dm.ModelToDTO(device))
	}

	ctx.JSON(http.StatusOK, r)
}

// DeviceCreate godoc
//
//	@Summary	Create a device
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		data	body		dto.CreateDeviceRequest	true	"Create data"
//	@Success	201		{object}	dto.CreateDeviceResponse
//	@Failure	400		{object}	errors.HTTPError
//	@Failure	401		{object}	errors.HTTPError
//	@Failure	404		{object}	errors.HTTPError
//	@Failure	500		{object}	errors.HTTPError
//	@Router		/devices [post]
func (h *deviceHandler) Create(ctx *gin.Context) {
	var err error
	var userID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	body := &dto.CreateDeviceRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	device := h.dm.CreateDTOToCreateModel(body)
	device.UserID = userID

	var deviceID uuid.UUID
	deviceID, err = h.ds.Create(ctx, device)
	if err != nil {
		if errors.Is(err, ae.ErrBrokerNotFound) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.CreateDeviceResponse{
		ID: deviceID,
	})
}

// DeviceUpdate godoc
//
//	@Summary	Update a device
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		deviceId	path	string					true	"Device UUID"
//	@Param		data		body	dto.UpdateDeviceRequest	true	"Update data"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/devices/{deviceId} [patch]
func (h *deviceHandler) Update(ctx *gin.Context) {
	var err error
	var userID, deviceID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	deviceID, err = h.getDeviceID(ctx)
	if err != nil {
		return
	}

	body := &dto.UpdateDeviceRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	device := h.dm.UpdateDTOToUpdateModel(body)
	device.UserID = userID
	device.ID = deviceID

	err = h.ds.Update(ctx, device)
	if err != nil {
		if errors.Is(err, ae.ErrDeviceNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		} else if errors.Is(err, ae.ErrBrokerNotFound) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeviceDelete godoc
//
//	@Summary	Delete a device
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		deviceId	path	string	true	"Device UUID"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/devices/{deviceId} [delete]
func (h *deviceHandler) Delete(ctx *gin.Context) {
	var err error
	var userID, deviceID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	deviceID, err = h.getDeviceID(ctx)
	if err != nil {
		return
	}

	err = h.ds.Delete(ctx, deviceID, userID)
	if err != nil {
		if errors.Is(err, ae.ErrDeviceNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeviceListControls godoc
//
//	@Summary	List a device controls
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		deviceId	path		string	true	"Device UUID"
//	@Success	200			{array}		dto.GetDeviceControlResponse
//	@Failure	400			{object}	errors.HTTPError
//	@Failure	401			{object}	errors.HTTPError
//	@Failure	404			{object}	errors.HTTPError
//	@Failure	500			{object}	errors.HTTPError
//	@Router		/devices/{deviceId}/controls [get]
func (h *deviceHandler) ListControls(ctx *gin.Context) {
	var err error
	var userID, deviceID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	deviceID, err = h.getDeviceID(ctx)
	if err != nil {
		return
	}

	controls, err := h.dcs.List(ctx, userID, deviceID)
	if err != nil {
		if errors.Is(err, ae.ErrDeviceNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	r := []dto.GetDeviceControlResponse{}

	for _, control := range controls {
		r = append(r, *h.dcm.ModelToDTO(control))
	}

	ctx.JSON(http.StatusOK, r)
}

// DeviceCreateControl godoc
//
//	@Summary	Create a device control
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		deviceId	path		string							true	"Device UUID"
//	@Param		data		body		dto.CreateDeviceControlRequest	true	"Create data"
//	@Success	201			{object}	dto.CreateDeviceControlResponse
//	@Failure	400			{object}	errors.HTTPError
//	@Failure	401			{object}	errors.HTTPError
//	@Failure	404			{object}	errors.HTTPError
//	@Failure	409			{object}	errors.HTTPError
//	@Failure	500			{object}	errors.HTTPError
//	@Router		/devices/{deviceId}/controls [post]
func (h *deviceHandler) CreateControl(ctx *gin.Context) {
	var err error
	var userID, deviceID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	deviceID, err = h.getDeviceID(ctx)
	if err != nil {
		return
	}

	body := &dto.CreateDeviceControlRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	control := h.dcm.CreateDTOToCreateModel(body)
	control.DeviceID = deviceID

	var controlID uuid.UUID
	controlID, err = h.dcs.Create(ctx, userID, control)
	if err != nil {
		if errors.Is(err, ae.ErrDeviceNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}
		if errors.Is(err, ae.ErrControlStateExists) {
			ctx.AbortWithStatusJSON(http.StatusConflict, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.CreateDeviceControlResponse{
		ID: controlID,
	})
}

// DeviceUpdateControl godoc
//
//	@Summary	Update a device control
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		deviceId	path	string							true	"Device UUID"
//	@Param		controlId	path	string							true	"Control UUID"
//	@Param		data		body	dto.UpdateDeviceControlRequest	true	"Update data"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	409	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/devices/{deviceId}/controls/{controlId} [patch]
func (h *deviceHandler) UpdateControl(ctx *gin.Context) {
	var err error
	var userID, deviceID, controlID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	deviceID, controlID, err = h.getDeviceControlIDs(ctx)
	if err != nil {
		return
	}

	body := &dto.UpdateDeviceControlRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return
	}

	control := h.dcm.UpdateDTOToUpdateModel(body)
	control.DeviceID = deviceID
	control.ID = controlID

	err = h.dcs.Update(ctx, userID, control)
	if err != nil {
		if errors.Is(err, ae.ErrDeviceNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}
		if errors.Is(err, ae.ErrDeviceControlNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}
		if errors.Is(err, ae.ErrControlStateExists) {
			ctx.AbortWithStatusJSON(http.StatusConflict, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeviceDeleteControl godoc
//
//	@Summary	Delete a device control
//	@Tags		Devices
//	@Security	BearerAuth
//	@Accept		json
//	@Produce	json
//	@Param		deviceId	path	string	true	"Device UUID"
//	@Param		controlId	path	string	true	"Control UUID"
//	@Success	204
//	@Failure	400	{object}	errors.HTTPError
//	@Failure	401	{object}	errors.HTTPError
//	@Failure	404	{object}	errors.HTTPError
//	@Failure	500	{object}	errors.HTTPError
//	@Router		/devices/{deviceId}/controls/{controlId} [delete]
func (h *deviceHandler) DeleteControl(ctx *gin.Context) {
	var err error
	var userID, deviceID, controlID uuid.UUID

	userID, err = h.getUserID(ctx)
	if err != nil {
		return
	}

	deviceID, controlID, err = h.getDeviceControlIDs(ctx)
	if err != nil {
		return
	}

	err = h.dcs.Delete(ctx, userID, deviceID, controlID)
	if err != nil {
		if errors.Is(err, ae.ErrDeviceNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}
		if errors.Is(err, ae.ErrDeviceControlNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, ae.NewHTTPError(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, ae.NewHTTPError(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *deviceHandler) getDeviceID(ctx *gin.Context) (uuid.UUID, error) {
	params := &dto.DeviceParams{}

	err := ctx.BindUri(params)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return uuid.Nil, err
	}

	var deviceID uuid.UUID
	deviceID, err = uuid.Parse(params.DeviceID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return uuid.Nil, err
	}

	return deviceID, nil
}

func (h *deviceHandler) getDeviceControlIDs(ctx *gin.Context) (uuid.UUID, uuid.UUID, error) {
	params := &dto.DeviceControlParams{}

	err := ctx.BindUri(params)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return uuid.Nil, uuid.Nil, err
	}

	var deviceID uuid.UUID
	deviceID, err = uuid.Parse(params.DeviceID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return uuid.Nil, uuid.Nil, err
	}

	var controlID uuid.UUID
	controlID, err = uuid.Parse(params.ControlID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return uuid.Nil, uuid.Nil, err
	}

	return deviceID, controlID, nil
}

func (h *deviceHandler) getUserID(ctx *gin.Context) (uuid.UUID, error) {
	userID, err := uuid.Parse(ctx.MustGet("UserID").(string))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ae.NewHTTPError(err))
		return uuid.Nil, err
	}

	return userID, nil
}
