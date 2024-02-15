package server

import "github.com/gin-gonic/gin"

type RoomsAPI struct{}

func (api *RoomsAPI) CreateRoom(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "pong"})
}

func (api *RoomsAPI) ListRooms(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "pong"})
}
