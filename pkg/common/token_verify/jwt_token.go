package token_verify

import (
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	"Share-Wallet/pkg/db"
	"time"

	"github.com/pkg/errors"

	go_redis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UID      string
	Platform string //login platform
	jwt.RegisteredClaims
}

func BuildClaims(uid, platform string, ttl int64) Claims {
	now := time.Now()
	return Claims{
		UID:      uid,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*24) * time.Hour)), //Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                        //Issuing time
			NotBefore: jwt.NewNumericDate(now),                                        //Begin Effective time
		}}
}

func DeleteToken(userID string, platformID int, gAuthTypeToken bool) error {
	m, err := db.DB.RedisDB.GetTokenMapByUidPid(userID, constant.PlatformIDToName(platformID))
	if err != nil && err != go_redis.Nil {
		return err
	}
	var deleteTokenKey []string
	for k, v := range m {
		_, err = GetClaimFromToken(k, gAuthTypeToken)
		if err != nil || v != constant.NormalToken {
			deleteTokenKey = append(deleteTokenKey, k)

		}
	}
	if len(deleteTokenKey) != 0 {
		err = db.DB.RedisDB.DeleteTokenByUidPid(userID, platformID, deleteTokenKey)
		return err
	}
	return nil
}

func CreateToken(userID string, platformID int, gAuthTypeToken bool) (string, int64, error) {
	claims := BuildClaims(userID, constant.PlatformIDToName(platformID), config.Config.TokenPolicy.AccessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var tokenString string
	var err error
	if gAuthTypeToken {
		tokenString, err = token.SignedString([]byte(config.Config.TokenPolicy.AccessSecretGAuth))
	} else {
		tokenString, err = token.SignedString([]byte(config.Config.TokenPolicy.AccessSecret))
	}
	if err != nil {
		return "", 0, err
	}
	//remove Invalid token
	m, err := db.DB.RedisDB.GetTokenMapByUidPid(userID, constant.PlatformIDToName(platformID))
	if err != nil && err != go_redis.Nil {
		return "", 0, err
	}
	var deleteTokenKey []string
	//One account login restriction
	for k, _ := range m {
		// _, err = GetClaimFromToken(k, gAuthTypeToken)
		// if err != nil || v != constant.NormalToken {
		deleteTokenKey = append(deleteTokenKey, k)
		// }
	}
	if len(deleteTokenKey) != 0 {
		err = db.DB.RedisDB.DeleteTokenByUidPid(userID, platformID, deleteTokenKey)
		if err != nil {
			return "", 0, err
		}
	}
	err = db.DB.RedisDB.AddTokenFlag(userID, platformID, tokenString, constant.NormalToken)
	if err != nil {
		return "", 0, err
	}
	return tokenString, claims.ExpiresAt.Time.Unix(), err
}

func secret(gAuthTypeToken bool) jwt.Keyfunc {
	if gAuthTypeToken {
		return func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.TokenPolicy.AccessSecretGAuth), nil
		}
	}
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenPolicy.AccessSecret), nil
	}
}

func GetClaimFromToken(tokensString string, gAuthTypeToken bool) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secret(gAuthTypeToken))
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, &constant.ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, &constant.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, &constant.ErrTokenNotValidYet
			} else {
				return nil, &constant.ErrTokenUnknown
			}
		} else {
			return nil, &constant.ErrTokenNotValidYet
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			//log.NewDebug("", claims.UID, claims.Platform)
			return claims, nil
		}
		return nil, &constant.ErrTokenNotValidYet
	}
}

func GetUserIDFromToken(token string, operationID string) (bool, string, string) {
	gAuthTypeToken := false
	claims, err := ParseToken(token, operationID, gAuthTypeToken)
	if err != nil {
		log.Error(operationID, "ParseToken failed, ", err.Error(), token)
		return false, "", err.Error()
	}
	log.Debug(operationID, "token claims.ExpiresAt.Second() ", claims.ExpiresAt.Unix())
	return true, claims.UID, ""
}

func GetAdminUserIDFromToken(token string, operationID string, gAuthTypeToken bool) (bool, string, string) {
	claims, err := ParseToken(token, operationID, gAuthTypeToken)
	if err != nil {
		log.Error(operationID, "ParseToken failed, ", err.Error(), token)
		return false, "", err.Error()
	}
	log.Debug(operationID, "token claims.ExpiresAt.Second() ", claims.ExpiresAt.Unix())
	return true, claims.UID, ""
}

func GetUserIDFromTokenExpireTime(token string, operationID string) (bool, string, string, int64) {
	gAuthTypeToken := false
	claims, err := ParseToken(token, operationID, gAuthTypeToken)
	if err != nil {
		log.Error(operationID, "ParseToken failed, ", err.Error(), token)
		return false, "", err.Error(), 0
	}
	return true, claims.UID, "", claims.ExpiresAt.Unix()
}

func ParseTokenGetUserID(token string, operationID string) (error, string) {
	gAuthTypeToken := false
	claims, err := ParseToken(token, operationID, gAuthTypeToken)
	if err != nil {
		return err, ""
	}
	return nil, claims.UID
}

func ParseToken(tokensString, operationID string, gAuthTypeToken bool) (claims *Claims, err error) {
	claims, err = GetClaimFromToken(tokensString, gAuthTypeToken)
	if err != nil {
		log.NewError(operationID, "token validate err", err.Error(), tokensString)
		return nil, err
	}

	m, err := db.DB.RedisDB.GetTokenMapByUidPid(claims.UID, claims.Platform)
	if err != nil {
		log.NewError(operationID, "get token from redis err", err.Error(), tokensString)
		return nil, errors.Wrap(&constant.ErrTokenInvalid, "get token from redis err")
	}
	if m == nil {
		log.NewError(operationID, "get token from redis err", "m is nil", tokensString)
		return nil, errors.Wrap(&constant.ErrTokenInvalid, "get token from redis err")
	}
	if v, ok := m[tokensString]; ok {
		switch v {
		case constant.NormalToken:
			log.NewDebug(operationID, "this is normal return", claims)
			return claims, nil
		case constant.InValidToken:
			return nil, errors.Wrap(&constant.ErrTokenInvalid, "")
		case constant.KickedToken:
			log.Error(operationID, "this token has been kicked by other same terminal ", constant.ErrTokenKicked)
			return nil, errors.Wrap(&constant.ErrTokenKicked, "this token has been kicked by other same terminal ")
		case constant.ExpiredToken:
			return nil, errors.Wrap(&constant.ErrTokenExpired, "")
		default:
			return nil, errors.Wrap(&constant.ErrTokenUnknown, "")
		}
	}
	log.NewError(operationID, "redis token map not find", constant.ErrTokenUnknown)
	return nil, errors.Wrap(&constant.ErrTokenUnknown, "redis token map not find")
}

func ParseRedisInterfaceToken(redisToken interface{}, gAuthTypeToken bool) (*Claims, error) {
	return GetClaimFromToken(string(redisToken.([]uint8)), gAuthTypeToken)
}

func DeleteAdminTokenOnLogout(userID string, gAuthTypeToken bool) error {
	m, err := db.DB.RedisDB.GetTokenMapByUidPid(userID, constant.PlatformIDToName(constant.AdminPlatformID))
	if err != nil && err != go_redis.Nil {
		return errors.Wrap(err, "")
	}
	var deleteTokenKey []string
	for k, v := range m {
		deleteTokenKey = append(deleteTokenKey, k)
		log.Error("Token Added for Delete, ", k, v)
	}
	if len(deleteTokenKey) != 0 {
		err = db.DB.RedisDB.DeleteTokenByUidPid(userID, constant.AdminPlatformID, deleteTokenKey)
		return errors.Wrap(err, "")
	}
	return nil
}

// Validation token, false means failure, true means successful verification
func VerifyToken(token, uid string) (bool, error) {
	gAuthTypeToken := false
	claims, err := ParseToken(token, "", gAuthTypeToken)
	if err != nil {
		return false, errors.Wrap(err, "ParseToken failed")
	}
	if claims.UID != uid {
		return false, &constant.ErrTokenUnknown
	}
	log.NewDebug("", claims.UID, claims.Platform)
	return true, nil
}

func VerifyManagementToken(token, uid string) (bool, error) {
	gAuthTypeToken := false
	claims, err := ParseToken(token, "", gAuthTypeToken)

	if err != nil {
		return false, errors.Wrap(err, "ParseToken failed")
	}
	if claims.UID != uid {
		return false, &constant.ErrTokenUnknown
	}
	if claims.Platform != constant.PlatformIDToName(constant.AdminPlatformID) {
		return false, &constant.ErrTokenInvalid
	}
	log.NewDebug("", claims.UID, claims.Platform)
	return true, nil
}
