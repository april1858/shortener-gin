package service

import (
	"errors"
	"testing"

	repoMock "github.com/april1858/shortener-gin/internal/app/repository/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestService_CreatorShortened(t *testing.T) {
	var ctx *gin.Context = &gin.Context{}
	ctx.Set("UID", "1234567")
	mockCtrl := gomock.NewController(t)
	repo := repoMock.NewMockRepository(mockCtrl)

	TestService_CreatorSh := &Service{r: repo}

	tests := []struct {
		name        string
		originalURL string
		short       string
		err         error
	}{
		{name: "first", originalURL: "http://abcd1234.ru", short: "any", err: nil},
		{name: "second", originalURL: "http://abcd1234.ru", short: "any", err: errors.New("23505")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().Store(ctx, tt.originalURL, "1234567").Return(tt.short, tt.err)
			TestService_CreatorSh.CreatorShortened(ctx, tt.originalURL)
		})
	}

}
