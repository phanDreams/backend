package appauth

import (
	"context"
	"fmt"
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
		//Validate the signing algorithm to prevent “alg: none” and similar attacks
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok ||
			t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
        return s.jwtSecret, nil
    })

	if e != nil || !token.Valid {
		return "", "", dom.ErrInvalidCredentials
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", "", dom.ErrInvalidCredentials
	}
	subj := claims.Subject

	//check redis
	if e := s.cache.Get(ctx, oldRefreshToken).Err(); e  != nil {
		return "", "", dom.ErrInvalidCredentials
	}

	//delete old, generate new tokens
	if err := s.cache.Del(ctx, oldRefreshToken).Err(); err != nil {
		s.logger.Warn("could not delete old refresh token from cache",
        zap.String("token", oldRefreshToken),
        zap.Error(err),
	  )
	}
	newAccessToken, err = s.GenerateJWT(subj, s.accessTTL)
	if err != nil {
        return "", "", err
    }

	newRefreshToken, err = s.GenerateJWT(subj, s.refreshTTL)
	if err != nil {
        return newAccessToken, "", err
    }

	if e := s.cache.Set(ctx, newRefreshToken, subj, s.refreshTTL).Err(); e != nil {
		return "", "", fmt.Errorf("unable to persist new refresh token: %w", e)
	}

	return newAccessToken, newRefreshToken, nil

}