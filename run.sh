#!/bin/bash
go run cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go cmd/web/send-mail.go -dbname=bookings -dbuser=postgres -cache=false -production=false -dbpass=
