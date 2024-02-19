package p2p_server

import (
	"Geecache/group"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func GetRouter() *gin.Engine {
	r := gin.Default()

	BaseUrl := DefaultBasePath + "/:groupname" + "/:key"
	r.GET(BaseUrl, func(c *gin.Context) {
		// 获取参数
		groupname := c.Param("groupname")
		key := c.Param("key")
		g := group.GetGroup(groupname)

		// 请求头
		c.Header("Content-Type", "application/octet-stream")

		// 写出相关数据
		if b, e := g.Get(key); e == nil {
			fmt.Println(b.ByteSlice())
			c.Writer.Write(b.ByteSlice())
		} else {
			log.Println("group.Get error :", e)
			c.Writer.Write([]byte("error"))
		}
	})

	return r
}
