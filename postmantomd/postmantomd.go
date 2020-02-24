package postmantomd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"myTestGo/exception"
	"myTestGo/util"
	"strconv"
	"strings"
	"time"
)

type PostmanRoot struct {
	Info struct {
		Name string `json:"name"`
		Description string `json:"description"`
		PostmanID string `json:"_postman_id"`
	} `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanItem struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Request     PostmanRequest `json:"request"`
	Response    []struct {
		Name            string         `json:"name"`
		OriginalRequest PostmanRequest `json:"originalRequest"`
		Status          string         `json:"status"`
		Code            int            `json:"code"`
		Body            string         `json:"body"`
	} `json:"response"`
	Item []PostmanItem `json:"item"`
}

type PostmanRequest struct {
	Method string `json:"method"`
	Header []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"header"`
	Url struct {
		Raw   string   `json:"raw"`
		Host  []string `json:"host"`
		Path  []string `json:"path"`
		Query []struct {
			Key         string `json:"key"`
			Value       string `json:"value"`
			Description string `json:"description"`
		} `json:"query"`
		Variable []struct {
			Key         string `json:"key"`
			Value       string `json:"value"`
			Description string `json:"description"`
		} `json:"variable"`
	} `json:"url"`
	Body struct {
		Mode     string `json:"mode"`
		Formdata []struct {
			Key         string `json:"key"`
			Value       string `json:"value"`
			Description string `json:"description"`
			Type        string `json:"type"`
		} `json:"formdata"`
		Raw     string `json:"raw"`
		Options struct {
			Raw struct {
				Language string `json:"language"`
			} `json:"raw"`
		} `json:"options"`
	} `json:"body"`
	Description string `json:"description"`
}

func PostManJSONToMarddown(inPath, outPath string) error {
	bytes, err := ioutil.ReadFile(inPath)
	if err != nil {
		return exception.NewSysErrorException(err.Error())
	}

	pm := PostmanRoot{}
	err = json.Unmarshal([]byte(string(bytes)), &pm)
	if err != nil {
		return exception.NewSysErrorException(err.Error())
	}

	out := "\n\n " + pm.Info.Description + "\n\n---\n\n"



	menu, content := parseItem(pm.Item)

	out += menu + "\n\n---\n\n" + content
	fmt.Println(out)
	err = ioutil.WriteFile(outPath, []byte(out), 0777)
	if err != nil {
		return exception.NewSysErrorException(err.Error())
	}

	return nil
}

func parseItem(items []PostmanItem) (string, string) {
	menu := ""
	content := ""

	for _, item := range items {
		if item.Request.Method != "" {
			// 请求
			id := strings.ReplaceAll(strings.Join(item.Request.Url.Path, "."),":","")
			name := strings.Split(item.Name, " ")[0]
			menu += "* [" + name + "](#" + id + ")\n\n"
			content += "\n\n### <a id=\"" + id + "\">API: " + name + "</a>\n\n"

			// markdown
			if item.Request.Description != "" {
				content += "**接口描述**\n\n" + item.Request.Description + "\n\n"
			}
			content += "**请求描述**\n\n```http\n" + item.Request.Method + `  ` + strings.Join(item.Request.Url.Path, "/") + "\n```\n\n"

			paramReq := item.Request
			typeTitle := "参考值"
			if len(item.Response) > 0 {
				typeTitle = "类型"
				paramReq = item.Response[0].OriginalRequest
			}
			if len(paramReq.Url.Variable) > 0 || len(paramReq.Url.Query) > 0 ||
				(paramReq.Body.Mode == "formdata" && len(paramReq.Body.Formdata) > 0) ||
				paramReq.Body.Mode == "raw" {
				content += "**参数描述**\n\n"

				params := "名称|" + typeTitle + "|描述\n---|---|---\n"
				if len(paramReq.Url.Variable) > 0 {
					content += "path:\n\n" + params
					for _, p := range paramReq.Url.Variable {
						content += p.Key + "|" + p.Value + "|" + p.Description + "\n"
					}
					content += "\n"
				}
				if len(paramReq.Url.Query) > 0 {
					content += "query:\n\n" + params
					for _, p := range paramReq.Url.Query {
						content += p.Key + "|" + p.Value + "|" + p.Description + "\n"
					}
					content += "\n"
				}
				if paramReq.Body.Mode == "formdata" && len(paramReq.Body.Formdata) > 0 {
					content += "body:\n\n" + params
					for _, p := range paramReq.Body.Formdata {
						if p.Type == "file" {
							p.Value = "选择本地文件对象"
						}
						content += p.Key + "|" + p.Value + "|" + p.Description + "\n"
					}
					content += "\n"
				}
				if paramReq.Body.Mode == "raw" {
					content += "特殊参数：\n\n"
					switch paramReq.Body.Options.Raw.Language {
					case "text":
						content += "`Content-Type`:`text/plain`\n\n```text\n" + paramReq.Body.Raw + "\n```"
					case "javascript":
						content += "`Content-Type`:`application/javascript`\n\n```javascript\n" + paramReq.Body.Raw + "\n```"
					case "json":
						content += "`Content-Type`:`application/json`\n\n```json\n" + paramReq.Body.Raw + "\n```"
					case "html":
						content += "`Content-Type`:`text/html`\n\n```html\n" + paramReq.Body.Raw + "\n```"
					case "xml":
						content += "`Content-Type`:`application/xml`\n\n```xml\n" + paramReq.Body.Raw + "\n```"
					default:
						content += "```\n" + paramReq.Body.Raw + "\n```"
					}
				}
			}

			if len(item.Response) > 0 {
				content += "\n\n**响应描述**\n\n"

				for _, res := range item.Response {
					content += "```json\n" + res.Body + "\n```\n\n"
				}
			}

			content += "---\n\n"
		} else {
			// 目录
			now := time.Now()
			id := "item_" + util.ToPinYin(item.Name) + "_" + strconv.Itoa(now.Nanosecond())
			menu += "<details open>\n<summary>" + item.Name + "</summary><blockquote>\n\n"
			content += "## <a id=\"" + id + "\">" + item.Name + "</a>\n\n" + item.Description + "\n\n"

			if len(item.Item) <= 0 {
				continue
			}
			m, c := parseItem(item.Item)
			menu += m + "</details>\n\n"
			content += c
		}
	}
	return menu, content
}
