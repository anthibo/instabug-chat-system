package main

import (
	"message_service/internal/handlers"
)

type AppDIContainer struct {
	ApiCmdHandlers *handlers.ApiCmdHandlers
}
