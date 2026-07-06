package utility

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	RefreshRedisKey  = "staff:refresh"
	UserRefreshRedisKey = "user:refresh"
)

type StaffClaims struct {
	StaffId   int64  `json:"staff_id"`
	Username  string `json:"username"`
	RealName  string `json:"real_name"`
	TokenType string `json:"token_type"`
	TokenId   string `json:"token_id,omitempty"`
	jwt.RegisteredClaims
}

type UserClaims struct {
	UserId    int64  `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	TokenId   string `json:"token_id,omitempty"`
	jwt.RegisteredClaims
}

func JwtSecret(ctx context.Context) []byte {
	secret, _ := g.Cfg().Get(ctx, "jwt.secret")
	return []byte(secret.String())
}

func JwtAccessExpire(ctx context.Context) time.Duration {
	v, _ := g.Cfg().Get(ctx, "jwt.access_expire")
	return time.Duration(v.Int64()) * time.Second
}

func JwtRefreshExpire(ctx context.Context) time.Duration {
	v, _ := g.Cfg().Get(ctx, "jwt.refresh_expire")
	return time.Duration(v.Int64()) * time.Second
}

func GenerateAccessToken(ctx context.Context, staffId int64, username, realName string) (string, error) {
	now := time.Now()
	claims := &StaffClaims{
		StaffId:   staffId,
		Username:  username,
		RealName:  realName,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(JwtAccessExpire(ctx))),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "gf-eshop",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret(ctx))
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	ExpireIn     int64  `json:"expire_in"`
	RefreshToken string `json:"refresh_token"`
	RefreshIn    int64  `json:"refresh_in"`
}

func GenerateTokenPair(ctx context.Context, staffId int64, username, realName string) (*TokenPair, error) {
	now := time.Now()
	accessExpire := JwtAccessExpire(ctx)
	refreshExpire := JwtRefreshExpire(ctx)

	accessClaims := &StaffClaims{
		StaffId:   staffId,
		Username:  username,
		RealName:  realName,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "gf-eshop",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(JwtSecret(ctx))
	if err != nil {
		return nil, err
	}

	tokenId := guid.S()
	refreshClaims := &StaffClaims{
		StaffId:   staffId,
		Username:  username,
		RealName:  realName,
		TokenType: TokenTypeRefresh,
		TokenId:   tokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "gf-eshop",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(JwtSecret(ctx))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		ExpireIn:     int64(accessExpire.Seconds()),
		RefreshToken: refreshToken,
		RefreshIn:    int64(refreshExpire.Seconds()),
	}, nil
}

func ParseStaffToken(ctx context.Context, tokenStr string) (*StaffClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &StaffClaims{}, func(t *jwt.Token) (interface{}, error) {
		return JwtSecret(ctx), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*StaffClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func GenerateUserAccessToken(ctx context.Context, userId int64, username string) (string, error) {
	now := time.Now()
	claims := &UserClaims{
		UserId:    userId,
		Username:  username,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(JwtAccessExpire(ctx))),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "gf-eshop",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret(ctx))
}

func GenerateUserTokenPair(ctx context.Context, userId int64, username string) (*TokenPair, error) {
	now := time.Now()
	accessExpire := JwtAccessExpire(ctx)
	refreshExpire := JwtRefreshExpire(ctx)

	accessClaims := &UserClaims{
		UserId:    userId,
		Username:  username,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "gf-eshop",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(JwtSecret(ctx))
	if err != nil {
		return nil, err
	}

	tokenId := guid.S()
	refreshClaims := &UserClaims{
		UserId:    userId,
		Username:  username,
		TokenType: TokenTypeRefresh,
		TokenId:   tokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "gf-eshop",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(JwtSecret(ctx))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		ExpireIn:     int64(accessExpire.Seconds()),
		RefreshToken: refreshToken,
		RefreshIn:    int64(refreshExpire.Seconds()),
	}, nil
}

func ParseUserToken(ctx context.Context, tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return JwtSecret(ctx), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func GetUserClaims(ctx context.Context) *UserClaims {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil
	}
	if v := r.GetCtxVar("user_claims"); !v.IsNil() {
		if claims, ok := v.Val().(*UserClaims); ok {
			return claims
		}
	}
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return nil
	}
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	} else {
		return nil
	}
	claims, err := ParseUserToken(ctx, tokenStr)
	if err != nil {
		return nil
	}
	return claims
}

func GetStaffClaims(ctx context.Context) *StaffClaims {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil
	}
	if v := r.GetCtxVar("staff_claims"); !v.IsNil() {
		if claims, ok := v.Val().(*StaffClaims); ok {
			return claims
		}
	}
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return nil
	}
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	} else {
		return nil
	}
	claims, err := ParseStaffToken(ctx, tokenStr)
	if err != nil {
		return nil
	}
	return claims
}
