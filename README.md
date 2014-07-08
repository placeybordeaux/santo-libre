santo-libre
===========

This is a small add on for martini to expose all the routes for an API

Simply add ```m.Get("/api", libre.ExposeRoutes(m))``` as the last route added and enjoy!

Depends on https://github.com/danielgtaylor/aglio for now.
