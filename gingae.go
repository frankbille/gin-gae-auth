// Package gingae provides Gin middleware, handling authentication against the
// Google App Engine users service.
package gingae

import (
	"appengine"
	"appengine/user"
	"github.com/gin-gonic/gin"
)

const (
	// ContextKey is the key name for the GAE context within the Gin context
	ContextKey = "GaeContext"

	// UserKey is the key name for the GAE user within the Gin context
	UserKey = "GaeUser"

	// UserOAuthErrorKey is the key name for an OAuth error message within the Gin context
	UserOAuthErrorKey = "GaeUserOAuthError"
)

type gaeContextProvider func(c *gin.Context) appengine.Context

// GaeContext sets a variable on the Gin context, containing the GAE Context.
func GaeContext() gin.HandlerFunc {
	return gaeContextFromProvider(func(c *gin.Context) appengine.Context {
		return appengine.NewContext(c.Request)
	})
}

func gaeContextFromProvider(gaeContextProvider gaeContextProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx := gaeContextProvider(c)
		c.Set(ContextKey, gaeCtx)
	}
}

// GaeUser sets a variable on the Gin context, containing the GAE User.
func GaeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx, exists := c.Get(ContextKey)
		if exists == false {
			panic("Must use the GaeContext middleware before the GaeUser")
		}
		gaeUser := user.Current(gaeCtx.(appengine.Context))
		c.Set(UserKey, gaeUser)
	}
}

// GaeUserOAuth sets a variable on the Gin context, containing the GAE User, logged in using OAuth.
func GaeUserOAuth(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		gaeCtx, exists := c.Get(ContextKey)
		if exists == false {
			panic("Must use the GaeContext middleware before the GaeUserOAuth")
		}
		gaeUser, err := user.CurrentOAuth(gaeCtx.(appengine.Context), scope)
		if err != nil {
			c.Set(UserOAuthErrorKey, err)
		} else {
			c.Set(UserKey, gaeUser)
		}
	}
}
