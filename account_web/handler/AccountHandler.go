package handler

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mic-study/account_srv/proto/pb"
	"mic-study/account_web/req"
	"mic-study/account_web/res"
	"mic-study/custom_error"
	"mic-study/internal"
	"mic-study/jwt_op"
	"mic-study/log"
	"net/http"
	"strconv"
	"time"
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

func init() {
	err := initConsul()
	if err != nil {
		panic(err)
	}
	err = initGrpcClient()
	if err != nil {
		panic(err)
	}
}
func initConsul() error {
	defaultConfig := api.DefaultConfig()
	consulAddr := fmt.Sprintf("%s:%d",
		internal.AppConf.ConsulConfig.Host,
		internal.AppConf.ConsulConfig.Port)
	defaultConfig.Address = consulAddr
	consulClient, err := api.NewClient(defaultConfig)
	if err != nil {
		zap.S().Error("AccountListHandler, 创建consul client失败: " + err.Error())
		return err
	}

	serviceList, err := consulClient.Agent().ServicesWithFilter("Service=account-srv")
	if err != nil {
		zap.S().Error("AccountListHandler, 创建consul 获取服务列表失败: " + err.Error())
		return err
	}
	for _, v := range serviceList {
		accountSrvHost = v.Address
		accountSrvPort = v.Port
	}
	return nil
}

func initGrpcClient() error {
	grpcAddr := fmt.Sprintf("%s:%d", accountSrvHost, accountSrvPort)
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		s := fmt.Sprintf("AccountListHandler-Grpc拨号重构:%s", err.Error())
		log.Logger.Info(s)
		return err
	}
	client = pb.NewAccountServiceClient(conn)
	return nil
}

var accountSrvHost string
var accountSrvPort int
var client pb.AccountServiceClient

func AccountListHandler(c *gin.Context) {
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "3")
	pageNo, err := strconv.ParseInt(pageNoStr, 10, 32)
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)

	r, err := client.GetAccountList(c, &pb.PagingRequest{
		PageNo:   uint32(pageNo),
		PageSize: uint32(pageSize),
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

func LoginByPasswordHandler(c *gin.Context) {
	var loginParam req.LoginByPassword
	err := c.ShouldBind(&loginParam)
	if err != nil {
		log.Logger.Error("LoginByPassword出错:" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "解析参数错误",
		})
		return
	}
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		log.Logger.Error("LoginByPassword 拨号出错:" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	client := pb.NewAccountServiceClient(conn)
	r, err := client.GetAccountByMobile(context.Background(), &pb.MobileRequest{Mobile: loginParam.Mobile})
	if err != nil {
		log.Logger.Error("GetAccountByMobile 出错:" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	cheRes, err := client.CheckPassword(c, &pb.CheckPasswordRequest{
		Password:    loginParam.Password,
		HasPassword: r.Password,
		AccountId:   uint32(r.Id),
	})
	if err != nil {
		log.Logger.Error("CheckPassword 出错:" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	checkResult := "登录失败"
	if cheRes.Result {
		checkResult = "登录成功"
		j := jwt_op.NewJWT()
		now := time.Now()
		claims := jwt_op.CustomClaims{
			StandardClaims: jwt.StandardClaims{
				NotBefore: now.Unix(),
				ExpiresAt: now.Add(time.Hour).Unix(),
			},
			ID:          r.Id,
			NickName:    r.Nickname,
			AuthorityId: int32(r.Role),
		}
		token, err := j.GenerateJWT(claims)
		if err != nil {
			log.Logger.Error("GenerateJWT 出错:" + err.Error())
			e := HandleError(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": e,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":    "ok",
			"result": checkResult,
			"token":  token,
		})
	}
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
