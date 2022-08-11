// @Description 过滤器
// @Author weitao.yin@shopee.com
// @Since 2022/7/12

package base

type Filter interface {
	DoFilter(ctx *GContext) bool

	MatchUrl(ctx *GContext) bool
}
