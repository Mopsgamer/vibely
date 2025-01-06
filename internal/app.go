package internal

import (
	"fmt"
	"reflect"
	"restapp/internal/controller"
	"restapp/internal/controller/controller_http"
	"restapp/internal/controller/controller_ws"
	"restapp/internal/controller/database"
	"restapp/internal/controller/model_http"
	"restapp/internal/controller/model_ws"
	"restapp/internal/docsgen"
	"restapp/internal/environment"
	"restapp/websocket"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/graphql-go/graphql"
)

// Initialize gofiber application, including DB and view engine.
func NewApp() (*fiber.App, error) {
	environment.WaitForBuild()

	db, errDBLoad := database.InitDB()
	if errDBLoad != nil {
		log.Error(errDBLoad)
		return nil, errDBLoad
	}

	schema, graphqlFields, errGraphqlLoad := initGraphql(*db)
	if errGraphqlLoad != nil {
		log.Error(errGraphqlLoad)
		return nil, errGraphqlLoad
	}

	engine := NewAppHtmlEngine(db)
	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Use(logger.New())

	UseHttp := func(handler func(ctl controller_http.ControllerHttp) error) fiber.Handler {
		return func(ctx fiber.Ctx) error {
			ctl := controller_http.ControllerHttp{
				Ctx: ctx,
				DB:  *db,
			}

			if websocket.IsWebSocketUpgrade(ctx) {
				return ctx.Next()
			}

			return handler(ctl)
		}
	}

	UseHttpResp := func(resp controller_http.Response) fiber.Handler {
		return UseHttp(func(ctl controller_http.ControllerHttp) error {
			return resp.HandleHtmx(ctl)
		})
	}

	UseHttpPage := func(
		templatePath string,
		bind *fiber.Map,
		redirectOn controller_http.RedirectCompute,
		layouts ...string,
	) fiber.Handler {
		bindx := fiber.Map{
			"Title": "?",
		}
		bindx = controller.MapMerge(&bindx, bind)
		return UseHttp(func(ctl controller_http.ControllerHttp) error {
			x := new(model_http.MemberOfUriGroup)
			ctl.BindAll(x)
			rights, member, group, user := x.Rights(ctl)
			bindx = controller.MapMerge(&bindx, &fiber.Map{
				"User":   user,
				"Group":  group,
				"Member": member,
				"Rights": rights,
			})
			return ctl.RenderPage(
				templatePath,
				&bindx,
				redirectOn,
				layouts...,
			)
		})
	}

	UseWs := func(
		subscribe []string,
		handler func(ctl controller_ws.ControllerWs) error,
	) fiber.Handler {
		return func(c fiber.Ctx) error {
			ctlHttp := controller_http.ControllerHttp{
				Ctx: c,
				DB:  *db,
			}

			if !websocket.IsWebSocketUpgrade(ctlHttp.Ctx) {
				ctlHttp.Ctx.Next()
			}

			// ctlWs MUST be created before websocket handler,
			// because some methods become unavailable after protocol upgraded.
			ctlWs := controller_ws.New(ctlHttp)

			request := new(model_http.MemberOfUriGroup)
			ctlHttp.BindAll(request)
			ctlWs.User = request.User(ctlHttp)
			ctlWs.Group = request.Group(ctlHttp)

			websocket.New(func(conn *websocket.Conn) {
				ctlWs.Conn = conn

				controller_ws.UserSessionMap.Connect(ctlWs.User.Id, ctlWs)
				defer controller_ws.UserSessionMap.Close(ctlWs.User.Id, ctlWs)
				for !ctlWs.Closed {
					messageType, message, err := ctlWs.Conn.ReadMessage()
					if err != nil {
						break
					}

					// start := time.Now()
					ctlWs.MessageType = messageType
					ctlWs.Message = message
					err = handler(*ctlWs)

					// colorErr := fiber.DefaultColors.Green
					// if err != nil {
					// 	colorErr = fiber.DefaultColors.Red
					// }

					// fmt.Printf(
					// 	"%s | %s%3s%s | %13s | %15s | %d | %s%s%s \n",
					// 	time.Now().Format("15:04:05"),
					// 	colorErr,
					// 	"ws",
					// 	fiber.DefaultColors.Reset,
					// 	time.Since(start),
					// 	ws.IP,
					// 	ws.MessageType,
					// 	fiber.DefaultColors.Yellow,
					// 	ws.Message,
					// 	fiber.DefaultColors.Reset,
					// )
					if err != nil {
						break
					}
				}
				ctlWs.Conn.Close()
			})(ctlHttp.Ctx)

			return nil
		}
	}

	UseWsResp := func(
		subscribe []string,
		resp controller_ws.Response,
	) fiber.Handler {
		return UseWs(subscribe, func(ctl controller_ws.ControllerWs) error {
			return resp.HandleHtmx(ctl)
		})
	}

	// static
	app.Get("/static/*", static.New("./web/static", static.Config{Browse: true}))
	app.Get("/partials*", static.New("./web/templates/partials", static.Config{Browse: true}))

	// pages
	docs := docsgen.New()
	app.Get("/", UseHttpPage("homepage", &fiber.Map{"Title": "Discover", "IsHomePage": true}, func(r controller_http.ControllerHttp, bind *fiber.Map) string { return "" }, "partials/main"))
	app.Get("/terms", UseHttpPage("terms", &fiber.Map{"Title": "Terms", "CenterContent": true}, func(r controller_http.ControllerHttp, bind *fiber.Map) string { return "" }, "partials/main"))
	app.Get("/privacy", UseHttpPage("privacy", &fiber.Map{"Title": "Privacy", "CenterContent": true}, func(r controller_http.ControllerHttp, bind *fiber.Map) string { return "" }, "partials/main"))
	app.Get("/acknowledgements", UseHttpPage("acknowledgements", &fiber.Map{"Title": "Acknowledgements"}, func(r controller_http.ControllerHttp, bind *fiber.Map) string { return "" }, "partials/main"))
	app.Get("/docs", func(ctx fiber.Ctx) error {
		return ctx.Redirect().To("/docs/rest")
	})
	app.Get("/docs/rest", UseHttpPage("docs-rest", &fiber.Map{
		"Title":          "Rest Docs",
		"IsDocsPage":     true,
		"IsDocsPageRest": true,
		"Docs":           docs,
	}, func(r controller_http.ControllerHttp, bind *fiber.Map) string { return "" }, "partials/main"))
	app.Get("/docs/gql", UseHttpPage("docs-gql", &fiber.Map{
		"Title":          "GraphQL Docs",
		"IsDocsPage":     true,
		"IsDocsPageGql":  true,
		"GraphqlFields":  graphqlFields,
		"GraphqlTypes":   graphqlTypes,
		"GraphqlRequest": docsgen.FieldsOf(GraphqlInput{}),
	}, func(r controller_http.ControllerHttp, bind *fiber.Map) string { return "" }, "partials/main"))
	app.Get("/settings", UseHttpPage("settings", &fiber.Map{"Title": "Settings"},
		func(r controller_http.ControllerHttp, bind *fiber.Map) string {
			// FIXME:
			// if user := r.User(); user == nil {
			// 	return "/"
			// }
			return ""
		}, "partials/main"),
	)
	app.Get("/chat", UseHttpPage("chat", &fiber.Map{"Title": "Home", "IsChatPage": true},
		func(r controller_http.ControllerHttp, bind *fiber.Map) string {
			return ""
		}),
	)
	app.Get("/chat/groups/:group_id", UseHttpPage("chat", &fiber.Map{"Title": "Group", "IsChatPage": true},
		func(r controller_http.ControllerHttp, bind *fiber.Map) string {
			// FIXME:
			// member, _, group := r.Member()

			// if member == nil {
			// 	return "/chat"
			// }

			// (*bind)["Title"] = group.Nick
			return ""
		}),
	)
	app.Get("/chat/groups/join/:group_name", UseHttpPage("chat", &fiber.Map{"Title": "Join group", "IsChatPage": true},
		func(r controller_http.ControllerHttp, bind *fiber.Map) string {
			// FIXME:
			// member, group, _ := r.Member()
			// if group == nil {
			// 	return "/chat"
			// }

			// if member != nil {
			// 	return controller.PathRedirectGroup(group.Id)
			// }

			// (*bind)["Title"] = "Join " + group.Nick
			return ""
		}),
	)

	Listen := func(method, path string, response controller_http.Response) fiber.Router {
		return app.Add([]string{method}, path, UseHttpResp(response))
	}

	var ListenDoc = func(method, description string, path string, request []reflect.StructField, response controller_http.Response) fiber.Router {
		docs.HTTP[method] = append(docs.HTTP[method], docsgen.DocsHTTPMethod{
			Path:        path,
			Method:      method,
			Description: description,
			Request:     request,
			Response:    "Currently, it's always an html string response.",
		})

		return Listen(method, path, response)
	}

	// graphql
	app.Post("/gql", func(ctx fiber.Ctx) error {
		var input GraphqlInput
		if err := ctx.Bind().Body(&input); err != nil {
			return ctx.
				Status(fiber.StatusInternalServerError).
				SendString(err.Error())
		}

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  input.Query,
			OperationName:  input.OperationName,
			VariableValues: input.Variables,
		})

		ctx.Set("Content-Type", "application/graphql-response+json")
		return ctx.JSON(result)
	})

	// get
	ListenDoc(
		"get",
		"Get messages section.",
		"/groups/:group_id/messages/page/:messages_page",
		docsgen.FieldsOf(model_http.MessagesPage{}),
		new(model_http.MessagesPage),
	)
	ListenDoc(
		"get",
		"Get members section.",
		"/groups/:group_id/members/page/:members_page",
		docsgen.FieldsOf(model_http.MembersPage{}),
		new(model_http.MembersPage),
	)

	// post
	ListenDoc(
		"post",
		"Create new account.",
		"/account/create",
		docsgen.FieldsOf(model_http.UserSignUp{}),
		new(model_http.UserSignUp),
	)
	ListenDoc(
		"post",
		"Get new authorization token.",
		"/account/login",
		docsgen.FieldsOf(model_http.UserLogin{}),
		new(model_http.UserLogin),
	)
	ListenDoc(
		"post",
		"Create new group.",
		"/groups/create",
		docsgen.FieldsOf(model_http.GroupCreate{}),
		new(model_http.GroupCreate),
	)
	ListenDoc(
		"post",
		"Create (send) new message.",
		"/groups/:group_id/messages/create",
		docsgen.FieldsOf(model_http.MessageCreate{}),
		new(model_http.MessageCreate),
	)

	// put
	ListenDoc(
		"put",
		"Change name identificator.",
		"/account/change/name",
		docsgen.FieldsOf(model_http.UserChangeName{}),
		new(model_http.UserChangeName),
	)
	ListenDoc(
		"put",
		"Cahnge email.",
		"/account/change/email",
		docsgen.FieldsOf(model_http.UserChangeEmail{}),
		new(model_http.UserChangeEmail),
	)
	ListenDoc(
		"put",
		"Change phone.",
		"/account/change/phone",
		docsgen.FieldsOf(model_http.UserChangePhone{}),
		new(model_http.UserChangePhone),
	)
	ListenDoc(
		"put",
		"Change password.",
		"/account/change/password",
		docsgen.FieldsOf(model_http.UserChangePassword{}),
		new(model_http.UserChangePassword),
	)
	ListenDoc(
		"put",
		"Remove authorization cookie and refresh page.",
		"/account/logout",
		docsgen.FieldsOf(model_http.UserLogout{}),
		new(model_http.UserLogout),
	)
	ListenDoc(
		"put",
		"Join the group immediately.",
		"/groups/:group_id/join",
		docsgen.FieldsOf(model_http.GroupJoin{}),
		new(model_http.GroupJoin),
	)
	ListenDoc(
		"put",
		"Change group information and settings.",
		"/groups/:group_id/change",
		docsgen.FieldsOf(model_http.GroupChange{}),
		new(model_http.GroupChange),
	)

	// delete
	ListenDoc(
		"delete",
		"Leave the group immediately.",
		"/groups/:group_id/leave",
		docsgen.FieldsOf(model_http.GroupLeave{}),
		new(model_http.GroupLeave),
	)
	ListenDoc(
		"delete",
		"Delete group.",
		"/groups/:group_id",
		docsgen.FieldsOf(model_http.GroupDelete{}),
		new(model_http.GroupDelete),
	)
	ListenDoc(
		"delete",
		"Delete account.",
		"/account/delete",
		docsgen.FieldsOf(model_http.UserDelete{}),
		new(model_http.UserDelete),
	)

	// websoket
	app.Get("/groups/:group_id", UseWsResp(
		[]string{controller_ws.SubForMessages},
		&model_ws.WebsocketGroup{},
	))

	app.Use(UseHttpPage("partials/x", &fiber.Map{
		"Title":         fmt.Sprintf("%d", fiber.StatusNotFound),
		"StatusCode":    fiber.StatusNotFound,
		"StatusMessage": fiber.ErrNotFound.Message,
		"CenterContent": true,
	}, func(r controller_http.ControllerHttp, bind *fiber.Map) string { return "" }, "partials/main"))

	return app, nil
}
