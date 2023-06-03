package main

import (
	"content_portal/routes"
)

func main() {
	r := routes.StartApp()

	r.Run(":8000")
}
