/*
 * Juke It Out
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.1
 * Contact: alex@alexmeuer.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"github.com/gin-gonic/gin"
)

type RoomsAPI interface {
	// CreateRoom - Create room
	CreateRoom(ctx *gin.Context)
	// ListRooms - List all rooms that are visible to you
	ListRooms(ctx *gin.Context)
}
