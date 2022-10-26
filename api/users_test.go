package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type ArgMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (a ArgMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CompareHashAndPassword(arg.HashedPassword, a.password)
	if err != nil {
		return false
	}

	a.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(a.arg, arg)
}

func (a ArgMatcher) String() string {
	return fmt.Sprintf("%v matches %v", a.arg.HashedPassword, a.password)
}

func EqArgMatcher(arg db.CreateUserParams, password string) gomock.Matcher {
	return ArgMatcher{arg: arg, password: password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := createRandomUser(t)

	testSuite := []struct {
		name          string
		arg           gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(w *httptest.ResponseRecorder)
	}{
		{
			name: "StatusOK",
			arg: gin.H{
				"username": user.Username,
				"password": password,
				"fullname": user.Fullname,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					Email:          user.Email,
					HashedPassword: user.HashedPassword,
					Fullname:       user.Fullname,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqArgMatcher(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatchUser(t, w.Body, user)
			},
		},
		{
			name: "StatusBadRequest",
			arg: gin.H{
				"username": user.Username,
				"password": password,
				"fullname": user.Fullname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "StatusInternalServerError",
			arg: gin.H{
				"username": user.Username,
				"password": password,
				"fullname": user.Fullname,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "StatusForbidden",
			arg: gin.H{
				"username": user.Username,
				"password": password,
				"fullname": user.Fullname,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
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

			reqVal, err := json.Marshal(tc.arg)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqVal))
			server.router.ServeHTTP(w, req)

			tc.checkResponse(w)
		})

	}

}

func createRandomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Fullname:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	var gotUser db.User
	err := json.Unmarshal(body.Bytes(), &gotUser)
	require.NoError(t, err)

	fmt.Println("body", gotUser, "compare", user)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Fullname, gotUser.Fullname)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
