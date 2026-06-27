package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	DeviceCookieName = "device_id"
	DeviceMaxAge     = 365 * 24 * time.Hour // TODO: Make this configurable
)

func GetDeviceIDOrCreate(ctx *gin.Context) string {
	deviceID, err := ctx.Cookie(DeviceCookieName)
	if err == nil {
		return deviceID
	}

	// No Cookies found, create one for this device
	deviceID = uuid.NewString()
	ctx.SetCookie(DeviceCookieName, deviceID, int(DeviceMaxAge.Seconds()), "/", "", false, true)

	return deviceID
}

func GetExistingDeviceID(ctx *gin.Context) (string, error) {
	deviceID, err := ctx.Cookie(DeviceCookieName)
	if err == nil {
		return deviceID, nil
	}

	return "", errors.New("Missing device_id")
}
