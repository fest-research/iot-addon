package handler

var registeredAPIHandlers []*APIHandler = make([]*APIHandler, 0)

func registerAPIHandler(handler *APIHandler) {
	registeredAPIHandlers = append(registeredAPIHandlers, handler)
}

func GetAPIHandlers() []*APIHandler {
	return registeredAPIHandlers
}
