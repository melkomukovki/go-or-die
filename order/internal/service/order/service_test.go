package order

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	"github.com/melkomukovki/go-or-die/order/internal/service/order/mocks"
)

type mockTxManager struct {
	mock.Mock
}

func (m *mockTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(mock.Anything, mock.Anything)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return fn(ctx)
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	type fields struct {
		inventoryClient *mocks.InventoryClient
		orderRepo       *mocks.OrderRepository
		txManager       *mockTxManager
	}

	hullUUID := uuid.New()
	engineUUID := uuid.New()
	shieldUUID := uuid.New()

	tests := []struct {
		name    string
		req     model.CreateOrderInput
		setup   func(f fields)
		wantErr error
	}{
		{
			name: "успешное создание заказа",
			req: model.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
				ShieldUUID: &shieldUUID,
			},
			setup: func(f fields) {
				f.inventoryClient.On("ListParts", mock.Anything, []uuid.UUID{hullUUID, engineUUID, shieldUUID}).
					Return([]model.Part{
						{UUID: hullUUID, Price: 100, StockQuantity: 10},
						{UUID: engineUUID, Price: 200, StockQuantity: 5},
						{UUID: shieldUUID, Price: 50, StockQuantity: 1},
					}, nil)
				f.orderRepo.On("Create", mock.Anything, mock.MatchedBy(func(o model.Order) bool {
					return o.TotalPrice() == 350 && o.Status == model.OrderStatusPendingPayment
				})).Return(nil)
				f.txManager.On("Do", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "ошибка: детали нет на складе",
			req: model.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setup: func(f fields) {
				f.inventoryClient.On("ListParts", mock.Anything, []uuid.UUID{hullUUID, engineUUID}).
					Return([]model.Part{
						{UUID: hullUUID, Price: 100, StockQuantity: 0},
						{UUID: engineUUID, Price: 200, StockQuantity: 5},
					}, nil)
			},
			wantErr: errs.ErrOutOfStock,
		},
		{
			name: "ошибка клиента инвентаризации",
			req: model.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setup: func(f fields) {
				f.inventoryClient.On("ListParts", mock.Anything, []uuid.UUID{hullUUID, engineUUID}).
					Return(nil, errors.New("inventory error"))
			},
			wantErr: errors.New("inventory error"),
		},
		{
			name: "ошибка репозитория при создании",
			req: model.CreateOrderInput{
				HullUUID:   hullUUID,
				EngineUUID: engineUUID,
			},
			setup: func(f fields) {
				f.inventoryClient.On("ListParts", mock.Anything, []uuid.UUID{hullUUID, engineUUID}).
					Return([]model.Part{
						{UUID: hullUUID, Price: 100, StockQuantity: 10},
						{UUID: engineUUID, Price: 200, StockQuantity: 5},
					}, nil)
				f.txManager.On("Do", mock.Anything, mock.Anything).Return(errors.New("repo error"))
			},
			wantErr: errors.New("repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := fields{
				inventoryClient: mocks.NewInventoryClient(t),
				orderRepo:       mocks.NewOrderRepository(t),
				txManager:       &mockTxManager{},
			}

			if tt.setup != nil {
				tt.setup(f)
			}

			s := NewService(f.inventoryClient, nil, f.orderRepo, f.txManager)
			res, err := s.Create(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res.UUID)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	t.Parallel()

	type fields struct {
		orderRepo *mocks.OrderRepository
		txManager *mockTxManager
	}

	orderUUID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func(f fields)
		wantErr error
	}{
		{
			name: "успешное получение заказа",
			id:   orderUUID,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{UUID: orderUUID}, nil)
			},
			wantErr: nil,
		},
		{
			name: "заказ не найден",
			id:   orderUUID,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			wantErr: errs.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := fields{
				orderRepo: mocks.NewOrderRepository(t),
				txManager: &mockTxManager{},
			}

			if tt.setup != nil {
				tt.setup(f)
			}

			s := NewService(nil, nil, f.orderRepo, f.txManager)
			res, err := s.Get(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr) || err.Error() == tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, res.UUID)
			}
		})
	}
}

func TestService_Pay(t *testing.T) {
	t.Parallel()

	type fields struct {
		paymentClient *mocks.PaymentClient
		orderRepo     *mocks.OrderRepository
		txManager     *mockTxManager
	}

	orderUUID := uuid.New()
	transactionUUID := uuid.New()
	method := model.PaymentMethodCard

	tests := []struct {
		name    string
		id      uuid.UUID
		method  model.PaymentMethod
		setup   func(f fields)
		wantErr error
	}{
		{
			name:   "успешная оплата",
			id:     orderUUID,
			method: method,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusPendingPayment}, nil)
				f.paymentClient.On("PayOrder", mock.Anything, orderUUID, method).
					Return(transactionUUID, nil)
				f.orderRepo.On("Update", mock.Anything, mock.MatchedBy(func(o model.Order) bool {
					return o.UUID == orderUUID && o.Status == model.OrderStatusPaid && *o.TransactionUUID == transactionUUID
				})).Return(nil)
				f.txManager.On("Do", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "заказ не в статусе ожидания оплаты",
			id:     orderUUID,
			method: method,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusPaid}, nil)
			},
			wantErr: errs.ErrOrderAlreadyPaid,
		},
		{
			name:   "ошибка платежного клиента",
			id:     orderUUID,
			method: method,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusPendingPayment}, nil)
				f.paymentClient.On("PayOrder", mock.Anything, orderUUID, method).
					Return(uuid.UUID{}, errors.New("payment error"))
			},
			wantErr: errors.New("payment error"),
		},
		{
			name:   "ошибка репозитория при обновлении",
			id:     orderUUID,
			method: method,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusPendingPayment}, nil)
				f.paymentClient.On("PayOrder", mock.Anything, orderUUID, method).
					Return(transactionUUID, nil)
				f.txManager.On("Do", mock.Anything, mock.Anything).Return(errors.New("update error"))
			},
			wantErr: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := fields{
				paymentClient: mocks.NewPaymentClient(t),
				orderRepo:     mocks.NewOrderRepository(t),
				txManager:     &mockTxManager{},
			}

			if tt.setup != nil {
				tt.setup(f)
			}

			s := NewService(nil, f.paymentClient, f.orderRepo, f.txManager)
			res, err := s.Pay(context.Background(), tt.id, tt.method)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, transactionUUID, res)
			}
		})
	}
}

func TestService_Cancel(t *testing.T) {
	t.Parallel()

	type fields struct {
		orderRepo *mocks.OrderRepository
		txManager *mockTxManager
	}

	orderUUID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func(f fields)
		wantErr error
	}{
		{
			name: "успешная отмена заказа",
			id:   orderUUID,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusPendingPayment}, nil)
				f.orderRepo.On("Update", mock.Anything, mock.MatchedBy(func(o model.Order) bool {
					return o.UUID == orderUUID && o.Status == model.OrderStatusCancelled
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "заказ не найден",
			id:   orderUUID,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			wantErr: errs.ErrOrderNotFound,
		},
		{
			name: "заказ не в статусе ожидания оплаты",
			id:   orderUUID,
			setup: func(f fields) {
				f.orderRepo.On("Get", mock.Anything, orderUUID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusPaid}, nil)
			},
			wantErr: errs.ErrOrderAlreadyPaid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := fields{
				orderRepo: mocks.NewOrderRepository(t),
				txManager: &mockTxManager{},
			}

			if tt.setup != nil {
				tt.setup(f)
			}

			s := NewService(nil, nil, f.orderRepo, f.txManager)
			err := s.Cancel(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr) || err.Error() == tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
