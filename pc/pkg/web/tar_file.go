package web

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"plugin_center/ent/ruleinfo"
	"plugin_center/pkg/store"
	"strings"
	"time"
)

func init() {
	RegisterHandler(func(engine *gin.Engine) {
		engine.GET("/opa/download/bundle.tar.gz", func(c *gin.Context) {
			desc := writeRuleToFile(c)
			defer deleteRuleDir(desc)

			gw := gzip.NewWriter(c.Writer)
			defer gw.Close()

			tw := tar.NewWriter(gw)
			defer tw.Close()

			info, err := os.Stat(desc)
			if err != nil {
				panic(err)
			}

			var baseDir string
			if info.IsDir() {
				baseDir = filepath.Base(desc)
			}

			fn := func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				header, err := tar.FileInfoHeader(info, info.Name())
				if err != nil {
					return err
				}

				if baseDir != "" {
					header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, desc))
				}

				if err := tw.WriteHeader(header); err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(tw, file)
				return err
			}
			if err = filepath.Walk(desc, fn); err != nil {
				panic(err)
			}
		})
	})
}

func deleteRuleDir(desc string) {
	_ = os.RemoveAll(desc)
}

const rootName = "rules"

// /tmp/xxx/rules/module/abcd.rego return /tmp/xxx
func writeRuleToFile(c *gin.Context) string {
	dir := filepath.Join(tempdir, fmt.Sprintf("%d", time.Now().UnixNano()), rootName)
	if err := os.MkdirAll(dir, 0777); err != nil {
		panic(err)
	}
	//wirte .manifest
	manifestContent := `{
  "revision" : "1",
  "roots": ["rules"]
}`
	file, err := os.OpenFile(filepath.Join(dir, ".manifest"), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err = file.WriteString(manifestContent); err != nil {
		panic(err)
	}
	x := store.Client.RuleInfo.Query().Where(ruleinfo.EnableEQ(true)).AllX(c)
	for _, info := range x {
		filename := fmt.Sprintf("%s.rego", info.RuleName)
		base := filepath.Join(dir, info.Module)
		err := os.MkdirAll(base, 0777)
		if err != nil {
			panic(err)
		}
		file, err := os.OpenFile(filepath.Join(dir, info.Module, filename), os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(info.RuleBody)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
	return dir
}
