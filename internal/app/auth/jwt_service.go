package appauth

import (
	"context"
	dom "pethelp-backend/internal/domain/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func (s *AuthService) GenerateJWT(subject string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject: subject,
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,  claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) RefreshToken(ctx context.Context, oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {
	//parse token
	token, e := jwt.ParseWithClaims(oldRefreshToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
        return s.jwtSecret, nil
    })

	if e != nil || !token.Valid {
		return "", "", dom.ErrInvalidCredentials
	}

	claims, _ := token.Claims.(*jwt.RegisteredClaims)
	subj := claims.Subject

	//check redis
	if e := s.cache.Get(ctx, oldRefreshToken).Err(); e  != nil {
		return "", "", dom.ErrInvalidCredentials
	}

	//delete old, generate new tokens
	s.cache.Del(ctx, oldRefreshToken)
	newAccessToken, err = s.GenerateJWT(subj, s.accessTTL)
	if err != nil {
        return "", "", err
    }

	newRefreshToken, err = s.GenerateJWT(subj, s.refreshTTL)
	if err != nil {
        return newAccessToken, "", err
    }

	if e := s.cache.Set(ctx, newRefreshToken, subj, s.refreshTTL).Err(); e != nil {
		s.logger.Warn("could not store new refresh token", zap.Error(e))
	}

	return newAccessToken, newRefreshToken, nil

}