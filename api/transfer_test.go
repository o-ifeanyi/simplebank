package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	acc := createRandomAccount()
	acc.Currency = util.USD
	acc2 := createRandomAccount()
	acc2.Currency = util.NGN
	transfer := createRandomTransfer()
	arg := db.CreateTransferParams{
		FromAccountID: transfer.FromAccountID,
		ToAccountID:   transfer.ToAccountID,
		Amount:        transfer.Amount,
	}

	testSuite := []struct {
		name          string
		arg           db.CreateTransferParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(w *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			arg:  arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.FromAccountID)).
					Times(1).
					Return(acc, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.ToAccountID)).
					Times(1).
					Return(acc, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TransferTxResult{}, nil)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "StatusBadRequest",
			arg:  db.CreateTransferParams{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.FromAccountID)).
					Times(0)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "StatusBadRequest Currency Mismatch",
			arg:  arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.FromAccountID)).
					Times(1).
					Return(acc, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.ToAccountID)).
					Times(1).
					Return(acc2, nil)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "StatusNotFound",
			arg:  arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.FromAccountID)).
					Times(1).
					Return(acc, sql.ErrNoRows)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			arg:  arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.FromAccountID)).
					Times(1).
					Return(acc, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.ToAccountID)).
					Times(1).
					Return(acc, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "StatusInternalServerError Invalid Account 1",
			arg:  arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.FromAccountID)).
					Times(1).
					Return(acc, sql.ErrConnDone)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "StatusInternalServerError Invalid Account 2",
			arg:  arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.FromAccountID)).
					Times(1).
					Return(acc, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg.ToAccountID)).
					Times(1).
					Return(acc, sql.ErrConnDone)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tc := range testSuite {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := NewServer(store)

			tc.buildStubs(store)

			transferReq := transferReq{
				FromAccountID: tc.arg.FromAccountID,
				ToAccountID:   tc.arg.ToAccountID,
				Amount:        tc.arg.Amount,
				Currency:      util.USD,
			}
			reqVal, err := json.Marshal(transferReq)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/transfers", bytes.NewBuffer(reqVal))
			server.router.ServeHTTP(w, req)

			tc.checkResponse(w)
		})

	}

}

func createRandomTransfer() db.Transfer {
	return db.Transfer{
		ID:            int64(util.RandomInt(1, 1000)),
		FromAccountID: int64(util.RandomInt(1, 1000)),
		ToAccountID:   int64(util.RandomInt(1, 1000)),
		Amount:        int64(util.RandomAmount()),
	}
}
