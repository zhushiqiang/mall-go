package baseresponse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
	"mall-go/common/baseerror"
	"mall-go/common/i18n"
	"net/http"
)

func ParamError(ctx *gin.Context, err interface{}) {
	validErr,ok := err.(validator.ValidationErrors)
	if ok {
		fmt.Println(validErr)
		errMap := map[string]string{}
		for _, ve := range validErr {
			key := ve.FieldNamespace + "." + ve.Tag
			errMap[key] = i18n.ZhMessage[key]
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": errMap})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"message": i18n.ErrParam})
	return
}

func HttpResponse(ctx *gin.Context, res ,err interface{})  {
	baseerror, ok := err.(*baseerror.BaseError)
	if ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":baseerror.Error()})
		return
	}
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"message":err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data":res})
}
