// Package main is entry point of the entire application
package main

//	@title			Go REST API
//	@version		1.0
//	@description	A production-ready REST API boilerplate with Gin, GORM, and MinIO.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	mrhpn
//	@contact.url	http://www.swagger.io/support
//	@contact.email	mr.htetphyonaing@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer " followed by a space and then your token.

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	runApplication()
}
