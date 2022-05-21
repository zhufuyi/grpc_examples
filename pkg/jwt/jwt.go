package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	signingKey    = []byte("zaq12wsxmko0") // 默认
	expire        = 2 * time.Hour          // 默认过期时长
	issuer        = ""
	signingMethod = jwt.SigningMethodHS256 // 默认HS256
)

var (
	formatErr       = errors.New("invalid token format")
	expiredErr      = errors.New("token has expired")
	unverifiableErr = errors.New("the token could not be verified due to a signing problem")
	signatureErr    = errors.New("signature failure")
)

//SetSigningKey 自定义设置签名key
func SetSigningKey(key string) {
	if key != "" {
		signingKey = []byte(key)
	}
}

// SetExpire 设置过期时间，至少大于1s
func SetExpire(d time.Duration) {
	if d < time.Second {
		expire = time.Second
	}
	expire = d
}

// SetIssuer 设置发布人
func SetIssuer(s string) {
	if s != "" {
		issuer = s
	}
}

// SetIssuer 设置发布人，"HS256"、"HS384"、"HS512"
func SetSigningMethod(name string) {
	switch name {
	case "HS256":
		signingMethod = jwt.SigningMethodHS256
	case "HS384":
		signingMethod = jwt.SigningMethodHS384
	case "HS512":
		signingMethod = jwt.SigningMethodHS512
	}
}

// GenerateTokenStandard 生成token
func GenerateTokenStandard() (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expire).Unix(),
		Issuer:    issuer,
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString(signingKey)
}

// VerifyTokenStandard 验证token
func VerifyTokenStandard(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	// token有效
	if token.Valid {
		return nil
	}

	ve, ok := err.(*jwt.ValidationError)
	if ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return formatErr
		} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return expiredErr
		} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
			return unverifiableErr
		} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
			return signatureErr
		} else {
			return ve // 其他错误
		}
	}

	return signatureErr
}

// ------------------------------------------------------------------------------------------

// CustomClaims 自定义Claims
type CustomClaims struct {
	Uid  string `json:"uid"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// GenerateTokenStandard 生成token
func GenerateTokenWithCustom(uid string, role ...string) (string, error) {
	roleVal := ""
	if len(role) > 1 {
		roleVal = role[0]
	}
	claims := CustomClaims{
		uid,
		roleVal,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire).Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString(signingKey)
}

// VerifyTokenCustom 验证token
func VerifyTokenCustom(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, formatErr
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, expiredErr
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				return nil, unverifiableErr
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, signatureErr
			} else {
				return nil, ve
			}
		}
		return nil, signatureErr
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, signatureErr
}
