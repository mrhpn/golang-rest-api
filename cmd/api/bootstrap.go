package main

func runApplication() {
	cfg := setupConfig()                       // config
	logger := setupLogger(cfg)                 // logger
	db := setupDatabase(cfg)                   // database
	appCtx := setupAppContext(cfg, db, logger) // app context
	router := setupRouter(appCtx)              // router
	server := setupHTTPServer(cfg, router)     // server
	gracefulShutdown(cfg, server, db)          // graceful shutdown
}
