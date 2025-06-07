package auth

import (
	"Git-Dominik/Schipperkes-Vereniging/db"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthManager struct {
	Admin *db.Admin
	DB    *db.SchipperkesDB
}

func (authManager *AuthManager) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionUUID := session.Get("session-id")
		admin := authManager.Admin
		if sessionUUID != admin.AdminUUID {
			ctx.Redirect(http.StatusFound, "/admin/login")
			return
		}

		ctx.Next()
	}
}

func (authManager *AuthManager) LoginHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	email := ctx.PostForm("adminEmail")
	password := ctx.PostForm("adminPassword")
	admin := authManager.Admin
	err := bcrypt.CompareHashAndPassword(admin.HashedPassword, []byte(password))
	if err != nil || email != admin.Email {
		fmt.Println("Login failed")
		ctx.HTML(http.StatusOK, "loginfailed.html", gin.H{})
		return
	}
	sessionUuid := uuid.New()
	admin.AdminUUID = sessionUuid.String()

	session.Set("session-id", admin.AdminUUID)
	session.Save()
	authManager.DB.GormDB.Save(admin)
	// Tell htmx to go to /admin
	ctx.Header("HX-Redirect", "/admin")
	// ctx.Redirect(http.StatusFound, "/admin")
}
