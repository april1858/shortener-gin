package service

import (
	"errors"
	"testing"

	repoMock "github.com/april1858/shortener-gin/internal/app/repository/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestService_CreatorShortened(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set("UID", "1234567")
	mockCtrl := gomock.NewController(t)
	repo := repoMock.NewMockRepository(mockCtrl)

	TestServiceCreatorSh := &Service{r: repo}

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
			TestServiceCreatorSh.CreatorShortened(ctx, tt.originalURL)
		})
	}

}

func TestService_FindOriginalURL(t *testing.T) {
	ctx := &gin.Context{}
	//ctx.Set("UID", "1234567")
	mockCtrl := gomock.NewController(t)
	repo := repoMock.NewMockRepository(mockCtrl)

	TestServiceFindOriginal := &Service{r: repo}

	tests := []struct {
		name        string
		originalURL string
		short       string
		err         error
	}{
		{name: "first", originalURL: "http://abcd1234.ru", short: "12345678", err: nil},
		{name: "second", originalURL: "http://bcd1234.ru", short: "12345678", err: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().Find(ctx, tt.short).Return(tt.originalURL, tt.err)
			TestServiceFindOriginal.FindOriginalURL(ctx, tt.short)
		})
	}

}
