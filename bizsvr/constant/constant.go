// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/14

package constant

import "time"

const (
	UserInfoCachePrefix  = "et_user:"
	SessionIdCachePrefix = "et_sid:"

	ThreeDaysToExpire = time.Hour * 72
	MaxImageSize      = 10 * 1024 * 1024 // 10M

	UploadFileDir = "./upload/images"
)

var ImageExtSupports = [5]string{".jpg", ".jpeg", ".png"}
