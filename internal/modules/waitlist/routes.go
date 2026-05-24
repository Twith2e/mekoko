package waitlist

import "github.com/gin-gonic/gin"

func AddRoute(rg *gin.RouterGroup, handler *Handler) {
	waitlist := rg.Group("/waitlist")

	{
		waitlist.POST("/join", handler.JoinWaitlist)
		waitlist.GET("/count", handler.GetWaitlistCount)
		waitlist.GET("/entries", handler.FetchWaitlistedEmails)
	}
}
