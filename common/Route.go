package handlers

import(
	model "github.com/MullionGroup/go-website-flintpro-example/models"
)

// Routes defines the type Routes which is just an array (slice) of Route structs.
type Routes []model.Route

var routesProtected = Routes{
	model.Route{
		"LoginGooglePageHandler", // Name
		"GET",     // HTTP method
		"/",      // Route pattern
		LoginGooglePageHandler,
		false,
	},

	model.Route{
		"IndexPageHandler", // Name
		"GET",         // HTTP method
		"/index",      // Route pattern
		IndexPageHandler,
		false,
	},

	model.Route{
		"LoginHandlerGoogleAuth", // Name
		"GET",         // HTTP method
		"/login_google",      // Route pattern
		LoginHandlerGoogleAuth,
		false,
	},

	model.Route{
		"RedirectHandlerGoogleAuth", // Name
		"GET",         // HTTP method
		"/load_google",      // Route pattern
		RedirectHandlerGoogleAuth,
		false,
	},

	model.Route{
		"AboutUsPageHandler", // Name
		"GET",         // HTTP method
		"/about_us",      // Route pattern
		AboutUsPageHandler,
		false,
	},

	//model.Route{
	//	"EFSheetDownloadHandler", // Name
	//	"GET",         // HTTP method
	//	"/ef_sheet/{sheet_id}/download",      // Route pattern
	//	EFSheetDownloadHandler,
	//	false,
	//},

	model.Route{
		"LogoutGooglePageHandler", // Name
		"POST",         // HTTP method
		"/logout_google",      // Route pattern
		LogoutGooglePageHandler,
		false,
	},

	model.Route{
		"PostSimulationDataPageHandler", // Name
		"POST",         // HTTP method
		"/load_data",      // Route pattern
		PostSimulationDataPageHandler,
		false,
	},

	//model.Route{
	//	"LoadSimulationDataPageHandler", // Name
	//	"GET",         // HTTP method
	//	"/load_data",      // Route pattern
	//	LoadSimulationDataPageHandler,
	//	false,
	//},

	model.Route{
		"RunSimulationHandler", // Name
		"POST",         // HTTP method
		"/run_simulation",      // Route pattern
		RunSimulationHandler,
		false,
	},

	model.Route{
		"ExportSimulationHandler", // Name
		"POST",         // HTTP method
		"/export_simulation",      // Route pattern
		ExportSimulationHandler,
		false,
	},
}

var routesOpenPrefix = Routes{
	model.Route{
		"AssetsHandler", // Name
		"GET",         // HTTP method
		"/assets",      // Route pattern
		AssetsHandler,
		false,
	},
}