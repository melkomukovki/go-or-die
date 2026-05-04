package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/go-or-die/payment/internal/api/payment/v1/mocks"
	"github.com/melkomukovki/go-or-die/payment/internal/converter"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

func TestAPI_PayOrder(t *testing.T) {
	t.Parallel()

	serviceErr := errors.New("ошибка сервиса")

	tests := []struct {
		name       string
		req        *paymentv1.PayOrderRequest
		prepare    func(s *mocks.PaymentService, req *paymentv1.PayOrderRequest)
		wantResp   *paymentv1.PayOrderResponse
		wantCode   codes.Code
		wantErrMsg string
	}{
		{
			name: "успешно - CARD",
			req: &paymentv1.PayOrderRequest{
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			prepare: func(s *mocks.PaymentService, req *paymentv1.PayOrderRequest) {
				expected := converter.PayRequestToModel(req)

				s.EXPECT().
					Pay(mock.Anything, expected).
					Return("tx-uuid", nil).
					Once()
			},
			wantResp: &paymentv1.PayOrderResponse{
				TransactionUuid: "tx-uuid",
			},
			wantCode: codes.OK,
		},
		{
			name: "успешно - SBP",
			req: &paymentv1.PayOrderRequest{
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
			},
			prepare: func(s *mocks.PaymentService, req *paymentv1.PayOrderRequest) {
				expected := converter.PayRequestToModel(req)

				s.EXPECT().
					Pay(mock.Anything, expected).
					Return("tx-sbp", nil).
					Once()
			},
			wantResp: &paymentv1.PayOrderResponse{
				TransactionUuid: "tx-sbp",
			},
			wantCode: codes.OK,
		},
		{
			name: "ошибка сервиса",
			req: &paymentv1.PayOrderRequest{
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			prepare: func(s *mocks.PaymentService, req *paymentv1.PayOrderRequest) {
				expected := converter.PayRequestToModel(req)

				s.EXPECT().
					Pay(mock.Anything, expected).
					Return("", serviceErr).
					Once()
			},
			wantCode:   codes.InvalidArgument,
			wantErrMsg: serviceErr.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := mocks.NewPaymentService(t)
			if tt.prepare != nil {
				tt.prepare(service, tt.req)
			}

			api := NewAPI(service)

			resp, err := api.PayOrder(context.Background(), tt.req)

			if tt.wantCode != codes.OK {
				require.Error(t, err)
				require.Nil(t, resp)

				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Equal(t, tt.wantErrMsg, st.Message())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantResp, resp)
		})
	}
}
