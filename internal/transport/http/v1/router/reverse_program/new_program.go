package reverse_program

import (
	v1handler "reverse-study-server/internal/transport/http/v1/handler/reverse_program"

	"github.com/gin-gonic/gin"
)

func RegisterReverseProgramRouter(r gin.IRouter) {
	reverseProgramGroup := r.Group("/reverse-program")
	{
		reverseProgramGroup.GET("/list", v1handler.ListPrograms)
		reverseProgramGroup.POST("/new", v1handler.NewProgram)
		reverseProgramGroup.PUT("/:id", v1handler.UpdateProgram)
		reverseProgramGroup.PUT("/:id/publish", v1handler.PublishProgram)
		reverseProgramGroup.DELETE("/:id", v1handler.DeleteProgram)
	}
}
