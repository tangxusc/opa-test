package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/logging"
	"github.com/open-policy-agent/opa/sdk"
	"strings"
	"testing"
)

func TestBoundle(t *testing.T) {
	ctx := context.Background()

	engine := gin.Default()
	engine.Static("/", "/Users/tangxu/GolandProjects/opa-test")
	//server := http.FileServer(http.Dir("Users/tangxu/GolandProjects/opa-test/testdata/test"))
	//http.Handle("/test", server)
	//http.Handle("/air", http.FileServer(http.Dir("Users/tangxu/GolandProjects/opa-test/testdata/air/")))
	go func() {
		err := engine.Run(":8899")
		//err := http.ListenAndServe(":8899", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	config := fmt.Sprintf(`{
		"services": {
			"test": {
				"url": %q
			},
			"air": {
				"url": %q
			}
		},
		"bundles": {
			"test": {
				"resource": "testdata/bundle.tar.gz"
			},
			"air": {
				"resource": "testdata2/bundle.tar.gz"
			}
		},
		"decision_logs": {
			"console":true
		}

	}`, "http://localhost:8899/", "http://localhost:8899/")

	fmt.Println(config)

	standardLogger := logging.New()
	standardLogger.SetLevel(logging.Debug)

	opa, err := sdk.New(ctx, sdk.Options{
		Config: strings.NewReader(config),
		Logger: standardLogger,
	})
	if err != nil {
		t.Fatal(err)
	}

	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(`{
			"sacNo": "4613810890"
		}`), &m)
	if err != nil {
		println(err)
	}

	defer opa.Stop(ctx)

	result, err := opa.Decision(ctx, sdk.DecisionOptions{Path: "/system", Input: m})
	if err != nil {
		println(err)
		panic(err)
	}
	marshal, _ := json.Marshal(result)
	fmt.Println(string(marshal))
}
