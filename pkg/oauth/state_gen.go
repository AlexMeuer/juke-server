package oauth

import (
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type GinStateGenerator struct{}

func (g *GinStateGenerator) GenerateState(ctx *gin.Context) string {
	return ksuid.New().String()
}
