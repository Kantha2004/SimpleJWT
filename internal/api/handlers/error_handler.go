package handlers

import (
	"fmt"
	"log"

	"github.com/Kantha2004/SimpleJWT/internal/api/services"
	apiresponse "github.com/Kantha2004/SimpleJWT/internal/apiResponse"
	"github.com/gin-gonic/gin"
)

func handleServiceError(c *gin.Context, err error, defaultMessage string) {
	switch e := err.(type) {
	case *services.ValidationError:
		apiresponse.SendValidationError(c, fmt.Errorf(e.Message))
	case *services.ConflictError:
		apiresponse.SendConflict(c, e.Message)
	case *services.NotFoundError:
		apiresponse.SendNotFound(c, e.Message)
	case *services.UnauthorizedError:
		apiresponse.SendUnauthorized(c, e.Message)
	case *services.InternalError:
		log.Printf("Internal service error: %v", e.Err)
		apiresponse.SendInternalError(c, defaultMessage)
	default:
		log.Printf("Unknown service error: %v", err)
		apiresponse.SendInternalError(c, defaultMessage)
	}
}
