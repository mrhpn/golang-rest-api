package main

func runApplication() {
	cfg := setupConfig()                              // config
	logger := setupLogger(cfg)                        // logger
	db := setupDatabase(cfg)                          // database
	media := setupMedia(cfg)                          // storage
	appCtx := setupAppContext(cfg, db, logger, media) // app context
	router := setupRouter(appCtx)                     // router
	server := setupHTTPServer(cfg, router)            // server
	gracefulShutdown(cfg, server, db)                 // graceful shutdown
}
