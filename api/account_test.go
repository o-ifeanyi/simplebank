package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	user, _ := createRandomUser(t)
	acc := createRandomAccount(user.Username)
	arg := db.CreateAccountParams{
		Owner:    acc.Owner,
		Balance:  0,
		Currency: acc.Currency,
	}

	testSuite := []struct {
		name          string
		arg           db.CreateAccountParams
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(w *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			arg:  arg,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(acc, nil)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatchAccount(t, w.Body, acc)
			},
		},
		{
			name: "NoAuthorization",
			arg:  arg,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "StatusBadRequest",
			arg:  db.CreateAccountParams{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			arg:  arg,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
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
			server := newTestServer(t, store)

			tc.buildStubs(store)

			createAccountReq := createAccountReq{
				Currency: tc.arg.Currency,
			}
			reqVal, err := json.Marshal(createAccountReq)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(reqVal))

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(w, req)

			tc.checkResponse(w)
		})

	}

}

func TestGetAccountAPI(t *testing.T) {
	user, _ := createRandomUser(t)
	acc := createRandomAccount(user.Username)

	testSuite := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(w *httptest.ResponseRecorder)
	}{
		{
			name:      "StatusOK",
			accountID: acc.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc.ID)).
					Times(1).
					Return(acc, nil)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatchAccount(t, w.Body, acc)
			},
		},
		{
			name:      "StatusBadRequest",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:      "StatusNotFound",
			accountID: acc.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:      "StatusInternalServerError",
			accountID: acc.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(acc.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
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
			server := newTestServer(t, store)

			tc.buildStubs(store)

			url := fmt.Sprintf("/accounts/%v", tc.accountID)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(w, req)

			tc.checkResponse(w)
		})

	}
}

func TestListAccountsAPI(t *testing.T) {
	user, _ := createRandomUser(t)
	arg := db.ListAccountsParams{
		Owner:  user.Username,
		Offset: 0,
		Limit:  10,
	}

	testSuite := []struct {
		name          string
		arg           db.ListAccountsParams
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(w *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			arg:  arg,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Account{}, nil)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "StatusBadRequest",
			arg:  db.ListAccountsParams{},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			arg:  arg,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
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
			server := newTestServer(t, store)

			tc.buildStubs(store)

			listAccountReq := listAccountsReq{
				PageID:   1,
				PageSize: tc.arg.Limit,
			}
			url := fmt.Sprintf("/accounts?page_id=%v&page_size=%v", listAccountReq.PageID, listAccountReq.PageSize)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(w, req)

			tc.checkResponse(w)
		})

	}

}

func createRandomAccount(username string) db.Account {
	return db.Account{
		ID:       int64(util.RandomInt(1, 1000)),
		Owner:    username,
		Currency: util.RandomCurrency(),
		Balance:  int64(util.RandomAmount()),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, acc db.Account) {
	var gotAccount db.Account
	err := json.Unmarshal(body.Bytes(), &gotAccount)
	require.NoError(t, err)

	require.Equal(t, gotAccount, acc)
}
