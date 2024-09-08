package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"mic-study/account_srv/proto/pb"
	"mic-study/account_web/res"
	"mic-study/custom_error"
	"mic-study/log"
	"net/http"
)

func HandleError(err error) string {
	if err != nil {
		switch err.Error() {
		case custom_error.AccountNotFound:
			return custom_error.AccountNotFound
		case custom_error.SaltError:
			return custom_error.SaltError
		case custom_error.AccountExists:
			return custom_error.AccountExists
		default:
			return custom_error.InternalError
		}
	}
	return ""
}

func AccountListHandler(c *gin.Context) {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		s := fmt.Sprintf("AccountListHandler-Grpc拨号重构:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
	}
	client := pb.NewAccountServiceClient(conn)
	r, err := client.GetAccountList(c, &pb.PagingRequest{
		PageNo:   1,
		PageSize: 3,
	})
	if err != nil {
		s := fmt.Sprintf("AccountListHandler-GetAccountList:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
	}
	var resList []res.Account4Res
	for _, item := range r.AccountList {
		resList = append(resList, pb2Res(item))
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "ok",
		"total": r.Total,
		"data":  resList,
	})
}

func pb2Res(accountRes *pb.AccountRes) res.Account4Res {
	return res.Account4Res{
		Mobile:   accountRes.Mobile,
		NickName: accountRes.Nickname,
		Gender:   accountRes.Gender,
	}
}
