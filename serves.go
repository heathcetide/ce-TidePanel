package ce_TidePanel

import (
	"github.com/gin-gonic/gin"
	"github.com/heathcetide/ce-TidePanel/config"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

func ServerStart() {
	db, _ := config.GetDB()
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	r := gin.Default()

	r.Static("/static", "./static")
	r.StaticFile("favicon.ico", "../static/img/favicon.ico")
	r.LoadHTMLGlob("templates/*")
	r = CollectionRoute(r)

	port := viper.GetString("server.port")
	r.Run(":" + port)
}
