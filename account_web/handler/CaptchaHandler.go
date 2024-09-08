package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mic-study/internal"
	"net/http"
	"os"
	"time"
)

func CaptchaHandler(c *gin.Context) {
	mobile, ok := c.GetQuery("mobile")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	fileName := "data.png"
	f, err := os.Create(fileName)
	if err != nil {
		zap.S().Error("获取验证码失败")
		return
	}
	defer f.Close()
	var w io.WriterTo
	d := captcha.RandomDigits(captcha.DefaultLen)
	w = captcha.NewImage("", d, captcha.StdWidth, captcha.StdHeight)
	_, err = w.WriteTo(f)
	if err != nil {
		zap.S().Error("CaptchaHandler() 失败")
		return
	}
	fmt.Println(d)
	captcha := ""
	for _, item := range d {
		captcha += string(item)
	}
	fmt.Println(captcha)
	internal.RedisClient.Set(c, mobile, captcha, 180*time.Second)
	b64, err := GetBase64(fileName)
	if err != nil {
		zap.S().Error("CaptchaHandler() 失败")
		return
	}
	fmt.Println(b64)
	c.JSON(http.StatusOK, gin.H{
		"captcha": b64,
	})
}

func GetBase64(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	b := make([]byte, 102300)
	base64.StdEncoding.Encode(b, file)
	s := string(b)
	fmt.Println(s)
	return s, nil
}
