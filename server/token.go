package server

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/base64"
	"time"
)

type TokenIdGeneratorFunc func() string

func GenerateTokenId() string {

	token := uuid.New()
	token = base64.StdEncoding.EncodeToString([]byte(token))
	return token
}

type TokenGenerator interface {
	GenerateAccessToken(serverConfig *Config, grant Grant) *Token
	GenerateRefreshToken(serverConfig *Config, grant Grant) *Token
}

type DefaultTokenGenerator struct {
	tokenIdGenerator TokenIdGeneratorFunc
}

func (generator *DefaultTokenGenerator) GenerateAccessToken(config *Config, grant Grant) *Token {

	var expiration int64

	expiration = grant.AccessTokenExpiration()

	if expiration == 0 {

		expiration = config.DefaultAccessTokenExpires
	}

	return &Token{
		generator.tokenIdGenerator(),
		time.Now().UTC().Add(time.Duration(expiration) * time.Second).Unix(),
	}
}

func (generator *DefaultTokenGenerator) GenerateRefreshToken(config *Config, grant Grant) *Token {

	var expiration int64

	expiration = config.DefaultRefreshTokenExpires

	return &Token{
		generator.tokenIdGenerator(),
		time.Now().UTC().Add(time.Duration(expiration) * time.Second).Unix(),
	}
}

func (generator *DefaultTokenGenerator) TokenIdGenerator() TokenIdGeneratorFunc {

	return generator.tokenIdGenerator
}

func NewDefaultTokenGenerator() *DefaultTokenGenerator {

	return &DefaultTokenGenerator{GenerateTokenId}
}
