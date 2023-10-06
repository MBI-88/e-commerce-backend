package main

import (
	"fiber-backend/router"
	"fiber-backend/models"
	"fiber-backend/settings"
	"flag"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

var (
	debug *bool
)

func init() {
	debug = flag.Bool("debug", true, "set env vars")
	flag.Usage = func() {
		info := "[*] Settings:\nrunserver: start server\nmigrate: migrate models to online database"
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "%s", info)
	}
}

// interfaceServer show option for setting up the server
func interfaceServer() {
	flag.Parse()
	args := flag.Arg(0)
	op := settings.Settings{}.GetEnvVar(*debug)
	switch args {
	case "migrate":
		models.Migrate(op.DSN)

	case "runserver":
		models.DialDb(op.DSN,"./logs/error.log")
		app := fiber.New(fiber.Config{
			ServerHeader: "Nginx",
			AppName:      "D'Todo",
		})
		
		app.Use(cors.New(cors.Config{
			AllowOrigins:     op.AllowOrigins,
			AllowHeaders:     op.AllowHeaders,
			AllowCredentials: op.AllowCredentials,
			AllowMethods:     op.AllowMethods,
			MaxAge:           op.MaxAge,
		}), helmet.New(helmet.Config{
			XSSProtection: op.XSSprotection,
			CrossOriginOpenerPolicy: op.CrossOriginOpenerPolicy,
			CrossOriginResourcePolicy: op.CrossOriginResourcePolicy,
		}))
		
		/* 
		csrf.New(csrf.Config{
			Expiration:     op.CsrfExpi,
			CookieSecure:   op.CookieSecure,
			CookieSameSite: op.CookieSameSite,
		})
		*/

		app.Static("/static", "./static", fiber.Static{
			MaxAge: 5,
		})

		router.Router(app)
		app.Listen(fmt.Sprintf("%s:%d", op.Host, op.Port))

	default:
		fmt.Println("[-] Set options!")

	}

}

func main() {
	interfaceServer()
}
