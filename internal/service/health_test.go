package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/service"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestHealthServiceHealth(t *testing.T) {
	var cases = map[string]struct {
		mocking     func(healthRepository *mock.MockHealthRepository)
		expectedErr error
	}{
		"should health": {
			mocking: func(healthRepository *mock.MockHealthRepository) {
				healthRepository.EXPECT().Health(gomock.Any()).Return(nil)
			},
		},
		"should throw error when health": {
			mocking: func(healthRepository *mock.MockHealthRepository) {
				healthRepository.EXPECT().Health(gomock.Any()).Return(fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			healthRepositoryMock := mock.NewMockHealthRepository(ctrl)
			healthService := service.NewHealthService(healthRepositoryMock)

			cs.mocking(healthRepositoryMock)

			// when
			err := healthService.Health(ctx)

			// then
			assert.Equal(t, err, cs.expectedErr)
		})
	}
}

func BenchmarkHealthServiceHealth(b *testing.B) {
	// given
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	healthRepositoryMock := mock.NewMockHealthRepository(ctrl)
	healthService := service.NewHealthService(healthRepositoryMock)

	healthRepositoryMock.EXPECT().Health(gomock.Any()).AnyTimes().Return(nil)

	// when
	for i := 0; i < b.N; i++ {
		healthService.Health(ctx)
	}
}
