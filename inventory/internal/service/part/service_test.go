package part

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/melkomukovki/go-or-die/inventory/internal/model"
	mockPart "github.com/melkomukovki/go-or-die/inventory/internal/service/part/mocks"
)

func TestService_Get(t *testing.T) {
	t.Parallel()
	validID := uuid.New()
	repoErr := errors.New("ошибка репозитория")

	tests := []struct {
		name    string
		id      uuid.UUID
		prepare func(repo *mockPart.PartRepository)
		wantErr error
	}{
		{
			name: "корректный uuid",
			id:   validID,
			prepare: func(repo *mockPart.PartRepository) {
				repo.EXPECT().
					Get(mock.Anything, validID).
					Return(model.Part{}, nil).
					Once()
			},
		},
		{
			name: "ошибка репозитория",
			id:   validID,
			prepare: func(repo *mockPart.PartRepository) {
				repo.EXPECT().
					Get(mock.Anything, validID).
					Return(model.Part{}, repoErr).
					Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mockPart.NewPartRepository(t)
			tt.prepare(repo)
			s := NewService(repo)
			_, err := s.Get(context.Background(), tt.id)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	validID1 := uuid.New()
	validID2 := uuid.New()
	repoErr := errors.New("repo error")

	tests := []struct {
		name    string
		filter  model.PartFilter
		prepare func(repo *mockPart.PartRepository)
		wantErr error
	}{
		{
			name: "пустые uuid",
			filter: model.PartFilter{
				UUIDs: nil,
			},
			prepare: func(repo *mockPart.PartRepository) {
				repo.EXPECT().
					List(mock.Anything, model.PartFilter{UUIDs: nil}).
					Return([]model.Part{}, nil).
					Once()
			},
		},
		{
			name: "корректный uuid",
			filter: model.PartFilter{
				UUIDs: []uuid.UUID{validID1, validID2},
			},
			prepare: func(repo *mockPart.PartRepository) {
				repo.EXPECT().
					List(mock.Anything, model.PartFilter{
						UUIDs: []uuid.UUID{validID1, validID2},
					}).
					Return([]model.Part{}, nil).
					Once()
			},
		},
		{
			name: "ошибка репозитория",
			filter: model.PartFilter{
				UUIDs: []uuid.UUID{validID1},
			},
			prepare: func(repo *mockPart.PartRepository) {
				repo.EXPECT().
					List(mock.Anything, model.PartFilter{
						UUIDs: []uuid.UUID{validID1},
					}).
					Return(nil, repoErr).
					Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mockPart.NewPartRepository(t)
			if tt.prepare != nil {
				tt.prepare(repo)
			}
			s := NewService(repo)
			_, err := s.List(context.Background(), tt.filter)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
