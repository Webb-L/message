package main

import (
	lang "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"message/config"
	"message/database"
	"message/logs"
	"message/router"
	"message/utils"
)

//	@title			消息系统 API
//	@version		1.0
//	@description	简单又好用的消息服务。快来给你"项目"添加消息功能。
//	@termsOfService	https://github.com/webb-l

//	@contact.name	Webb
//	@contact.url	https://github.com/webb-l
//	@contact.email	822028533@qq.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:1204
//	@BasePath	/

//	@securityDefinitions.basic BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	// 初始化Gin引擎
	r := gin.Default()

	// 初始化配置
	config.InitConfig()

	// 初始日志
	logs.InitLog()
	defer logs.Sync()
	r.Use(utils.AccessLogger())

	// 国际化中间件
	r.Use(lang.Localize(lang.WithBundle(&lang.BundleCfg{
		DefaultLanguage:  language.Chinese,
		FormatBundleFile: "yaml",
		AcceptLanguage:   []language.Tag{language.Chinese, language.English},
		RootPath:         "resources/lang",
		UnmarshalFunc:    yaml.Unmarshal,
	})))

	// 初始化路由
	router.InitRouter(r)

	// 连接MySQL数据库
	database.InitMySQL()

	// 启动Gin引擎
	r.Run(":1204")
}
