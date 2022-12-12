package main

func AuthenticationControllers() {
	Auth := Authentication{}

	AuthenticationRouter.POST("/new/session", Auth.NewSession)
	AuthenticationRouter.POST("/new/record/user", Auth.CreateAccount)
	AuthenticationRouter.POST("/destroy/session", Auth.DestroySession)
	AuthenticationRouter.GET("/request/data", Auth.RequestData)
}