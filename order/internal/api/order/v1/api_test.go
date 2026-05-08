package v1

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/melkomukovki/go-or-die/order/internal/api/order/v1/mocks"
	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func TestAPI_CreateOrder(t *testing.T) {
	t.Parallel()

	serviceErr := errors.New("ошибка сервиса")
	orderUUID := uuid.New()
	hullUUID := uuid.New()
	engineUUID := uuid.New()

	tests := []struct {
		name    string
		req     *orderv1.CreateOrderRequest
		prepare func(s *mocks.OrderService)
		wantRes orderv1.CreateOrderRes
		wantErr bool
	}{
		{
			name: "успешный",
			req: &orderv1.CreateOrderRequest{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(model.Order{
						UUID:       orderUUID,
						TotalPrice: 1000,
					}, nil).
					Once()
			},
			wantRes: &orderv1.CreateOrderResponse{
				OrderUUID:  orderUUID,
				TotalPrice: 1000,
			},
		},
		{
			name: "некорретный uuid",
			req: &orderv1.CreateOrderRequest{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(model.Order{}, errs.ErrInvalidUUID).
					Once()
			},
			wantRes: &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: errs.ErrInvalidUUID.Error(),
			},
		},
		{
			name: "деталь не найдена",
			req: &orderv1.CreateOrderRequest{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(model.Order{}, errs.ErrPartNotFound).
					Once()
			},
			wantRes: &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: errs.ErrPartNotFound.Error(),
			},
		},
		{
			name: "нет в наличии",
			req: &orderv1.CreateOrderRequest{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(model.Order{}, errs.ErrOutOfStock).
					Once()
			},
			wantRes: &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: errs.ErrOutOfStock.Error(),
			},
		},
		{
			name: "внутренняя ошибка",
			req: &orderv1.CreateOrderRequest{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(model.Order{}, serviceErr).
					Once()
			},
			wantRes: &orderv1.CreateOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewOrderService(t)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			api := NewAPI(service)
			res, err := api.CreateOrder(context.Background(), tt.req)

			require.NoError(t, err)
			require.Equal(t, tt.wantRes, res)
		})
	}
}

func TestAPI_GetOrder(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	hullUUID := uuid.New()
	engineUUID := uuid.New()
	now := time.Now()

	tests := []struct {
		name    string
		params  orderv1.GetOrderParams
		prepare func(s *mocks.OrderService)
		wantRes orderv1.GetOrderRes
	}{
		{
			name: "успешный",
			params: orderv1.GetOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(model.Order{
						UUID:       orderUUID,
						HullUUID:   hullUUID,
						EngineUUID: engineUUID,
						TotalPrice: 1000,
						Status:     model.OrderStatusPendingPayment,
						CreatedAt:  now,
					}, nil).
					Once()
			},
			wantRes: &orderv1.OrderDto{
				OrderUUID:  orderUUID,
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
				TotalPrice: 1000,
				Status:     orderv1.OrderStatusPENDINGPAYMENT,
				CreatedAt:  now,
			},
		},
		{
			name: "заказ не найден",
			params: orderv1.GetOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound).
					Once()
			},
			wantRes: &orderv1.GetOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "заказ не найден",
			},
		},
		{
			name: "внутренняя ошибка",
			params: orderv1.GetOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Get(mock.Anything, orderUUID).
					Return(model.Order{}, errors.New("error")).
					Once()
			},
			wantRes: &orderv1.GetOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewOrderService(t)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			api := NewAPI(service)
			res, err := api.GetOrder(context.Background(), tt.params)

			require.NoError(t, err)
			require.Equal(t, tt.wantRes, res)
		})
	}
}

func TestAPI_PayOrder(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()
	txUUID := uuid.New()

	tests := []struct {
		name    string
		req     *orderv1.PayOrderRequest
		params  orderv1.PayOrderParams
		prepare func(s *mocks.OrderService)
		wantRes orderv1.PayOrderRes
	}{
		{
			name: "успех",
			req: &orderv1.PayOrderRequest{
				PaymentMethod: orderv1.PaymentMethodCARD,
			},
			params: orderv1.PayOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Pay(mock.Anything, orderUUID, model.PaymentMethodCard).
					Return(txUUID, nil).
					Once()
			},
			wantRes: &orderv1.PayOrderResponse{
				TransactionUUID: txUUID,
			},
		},
		{
			name: "заказ уже оплачен",
			req: &orderv1.PayOrderRequest{
				PaymentMethod: orderv1.PaymentMethodCARD,
			},
			params: orderv1.PayOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Pay(mock.Anything, orderUUID, model.PaymentMethodCard).
					Return(uuid.Nil, errs.ErrOrderAlreadyPaid).
					Once()
			},
			wantRes: &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: errs.ErrOrderAlreadyPaid.Error(),
			},
		},
		{
			name: "заказ не найден",
			req: &orderv1.PayOrderRequest{
				PaymentMethod: orderv1.PaymentMethodCARD,
			},
			params: orderv1.PayOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Pay(mock.Anything, orderUUID, model.PaymentMethodCard).
					Return(uuid.Nil, errs.ErrOrderNotFound).
					Once()
			},
			wantRes: &orderv1.PayOrderNotFound{
				Code:    http.StatusNotFound,
				Message: errs.ErrOrderNotFound.Error(),
			},
		},
		{
			name: "некорректный запрос",
			req: &orderv1.PayOrderRequest{
				PaymentMethod: orderv1.PaymentMethodCARD,
			},
			params: orderv1.PayOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Pay(mock.Anything, orderUUID, model.PaymentMethodCard).
					Return(uuid.Nil, errs.ErrInvalidUUID).
					Once()
			},
			wantRes: &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: errs.ErrInvalidUUID.Error(),
			},
		},
		{
			name: "внутренняя ошибка",
			req: &orderv1.PayOrderRequest{
				PaymentMethod: orderv1.PaymentMethodCARD,
			},
			params: orderv1.PayOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Pay(mock.Anything, orderUUID, model.PaymentMethodCard).
					Return(uuid.Nil, errors.New("error")).
					Once()
			},
			wantRes: &orderv1.PayOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewOrderService(t)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			api := NewAPI(service)
			res, err := api.PayOrder(context.Background(), tt.req, tt.params)

			require.NoError(t, err)
			require.Equal(t, tt.wantRes, res)
		})
	}
}

func TestAPI_CancelOrder(t *testing.T) {
	t.Parallel()

	orderUUID := uuid.New()

	tests := []struct {
		name    string
		params  orderv1.CancelOrderParams
		prepare func(s *mocks.OrderService)
		wantRes orderv1.CancelOrderRes
	}{
		{
			name: "успешный",
			params: orderv1.CancelOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Cancel(mock.Anything, orderUUID).
					Return(nil).
					Once()
			},
			wantRes: &orderv1.CancelOrderResponse{},
		},
		{
			name: "заказ не найден",
			params: orderv1.CancelOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Cancel(mock.Anything, orderUUID).
					Return(errs.ErrOrderNotFound).
					Once()
			},
			wantRes: &orderv1.CancelOrderNotFound{
				Code:    http.StatusNotFound,
				Message: errs.ErrOrderNotFound.Error(),
			},
		},
		{
			name: "конфликт",
			params: orderv1.CancelOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Cancel(mock.Anything, orderUUID).
					Return(errs.ErrOrderPendingPaymentMismatch).
					Once()
			},
			wantRes: &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: errs.ErrOrderPendingPaymentMismatch.Error(),
			},
		},
		{
			name: "внутренняя ошибка",
			params: orderv1.CancelOrderParams{
				OrderUUID: orderUUID,
			},
			prepare: func(s *mocks.OrderService) {
				s.EXPECT().
					Cancel(mock.Anything, orderUUID).
					Return(errors.New("error")).
					Once()
			},
			wantRes: &orderv1.CancelOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewOrderService(t)
			if tt.prepare != nil {
				tt.prepare(service)
			}

			api := NewAPI(service)
			res, err := api.CancelOrder(context.Background(), tt.params)

			require.NoError(t, err)
			require.Equal(t, tt.wantRes, res)
		})
	}
}
