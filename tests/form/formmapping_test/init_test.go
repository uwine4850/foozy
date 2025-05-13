package formmappingtest

// func TestMain(m *testing.M) {
// 	testinitcnf.InitCnf()
// 	mddl := middlewares.NewMiddleware()
// 	mddl.AsyncMddl(builtin_mddl.GenerateAndSetCsrf(1800, nil))
// 	render, err := tmlengine.NewRender()
// 	if err != nil {
// 		panic(err)
// 	}
// 	newRouter := router.NewRouter(manager.NewManager(render))
// 	newRouter.SetMiddleware(mddl)
// 	newRouter.Post("/mp-default-struct", mpDefaultStruct)
// 	newRouter.Post("/mp-empty-string-0-err", mpEmptyString0Err)
// 	newRouter.Post("/mp-empty-string-1-err", mpEmptyString1Err)
// 	newRouter.Post("/mp-empty-file-err", mpEmptyFileErr)
// 	newRouter.Post("/mp-empty-value", mpEmptyValue)
// 	newRouter.Post("/fill", fill)
// 	newRouter.Post("/fill-reflect-value", fillReflectValue)
// 	newRouter.Post("/mp-typed-struct", mpTypedMapper)
// 	serv := server.NewServer(tconf.PortFormMapping, newRouter, nil)
// 	go func() {
// 		err = serv.Start()
// 		if err != nil && !errors.Is(err, http.ErrServerClosed) {
// 			panic(err)
// 		}
// 	}()
// 	if err := server.WaitStartServer(tconf.PortFormMapping, 5); err != nil {
// 		panic(err)
// 	}
// 	exitCode := m.Run()
// 	os.Exit(exitCode)
// 	err = serv.Stop()
// 	if err != nil {
// 		panic(err)
// 	}
// }
