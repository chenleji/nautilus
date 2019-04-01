package helper

//func TestGet(t *testing.T) {
//	Data, err := json.Marshal(&map[string]string{
//		"active": "ok",
//	})
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	// start web server
//	go func (){
//		http.HandleFunc("/health",
//			func(w http.ResponseWriter, req *http.Request) {
//				w.Write(Data)
//			},
//		)
//
//		http.ListenAndServe(":9095", nil)
//	}()
//
//	// registry consul
//	consul := Consul{}.New()
//	consul.RegistryService("service", "9095", "/health")
//
//	time.Sleep(15 * time.Second)
//
//	// verify get operation
//	resp := make(map[string]string, 0)
//	client := HttpClient{Service: "service"}
//	if err := client.Get(&resp, "/health", []string{}); err != nil {
//		t.Error(err)
//	}
//
//	t.Log("response:", resp)
//}
