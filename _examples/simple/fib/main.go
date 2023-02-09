// Copyright © 2021 Hedzr Yeh.

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/hedzr/errors.v3"
	"strconv"

	"net/http"
)

func main() {
	app := gin.New()
	app.Use(gin.Logger(), gin.Recovery())
	app.GET("/fib/:n", fibonacciHandler)
	app.Run(":8105")
}

func fibonacciHandler(c *gin.Context) {
	n := c.Param("n")
	nn, err := strconv.Atoi(n)
	if err == nil {
		result := fibonacci(nn)
		_, _ = c.Writer.Write([]byte(fmt.Sprintf("%d", result)))
	} else {
		//_ = c.Error(errors.New("expecting a number with format: /fib/:number"))
		//c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		//	"msg": "expecting a number with format: /fib/:number",
		//})
		c.AbortWithStatusJSON(http.StatusBadRequest,
			c.Error(errors.New("expecting a number with format: /fib/:number")).JSON())
	}
}

func fibonacci(num int) int {
	if num < 2 {
		return 1
	}
	return fibonacci(num-1) + fibonacci(num-2)
}
