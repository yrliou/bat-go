// +build integration

package wallet

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brave-intl/bat-go/middleware"
	"github.com/brave-intl/bat-go/utils/altcurrency"
	"github.com/brave-intl/bat-go/utils/httpsignature"
	walletutils "github.com/brave-intl/bat-go/utils/wallet"
	uphold "github.com/brave-intl/bat-go/utils/wallet/provider/uphold"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type WalletControllersTestSuite struct {
	suite.Suite
}

func TestWalletControllersTestSuite(t *testing.T) {
	suite.Run(t, new(WalletControllersTestSuite))
}

func (suite *WalletControllersTestSuite) SetupSuite() {
	pg, err := NewPostgres("", false)
	suite.Require().NoError(err, "Failed to get postgres conn")

	m, err := pg.NewMigrate()
	suite.Require().NoError(err, "Failed to create migrate instance")

	ver, dirty, _ := m.Version()
	if dirty {
		suite.Require().NoError(m.Force(int(ver)))
	}
	if ver > 0 {
		suite.Require().NoError(m.Down(), "Failed to migrate down cleanly")
	}

	suite.Require().NoError(pg.Migrate(), "Failed to fully migrate")
}

func (suite *WalletControllersTestSuite) SetupTest() {
	suite.CleanDB()
}

func (suite *WalletControllersTestSuite) TearDownTest() {
	suite.CleanDB()
}

func (suite *WalletControllersTestSuite) CleanDB() {
	tables := []string{"claim_creds", "claims", "wallets", "issuers", "promotions"}

	pg, err := NewPostgres("", false)
	suite.Require().NoError(err, "Failed to get postgres conn")

	for _, table := range tables {
		_, err = pg.DB.Exec("delete from " + table)
		suite.Require().NoError(err, "Failed to get clean table")
	}
}
func noUUID() *uuid.UUID {
	return nil
}
func (suite *WalletControllersTestSuite) TestLinkWallet() {
	pg, err := NewPostgres("", false)
	suite.Require().NoError(err, "Failed to get postgres connection")

	service := &Service{
		Datastore: pg,
	}

	w1 := suite.NewWallet(service, "uphold")
	w2 := suite.NewWallet(service, "uphold")
	w3 := suite.NewWallet(service, "uphold")
	w4 := suite.NewWallet(service, "uphold")
	settlement := os.Getenv("BAT_SETTLEMENT_ADDRESS")

	upholdWallet1, ok := w1.(*uphold.Wallet)
	suite.Require().True(ok, "conversion to interface must succeed")
	upholdWallet2, ok := w2.(*uphold.Wallet)
	suite.Require().True(ok, "conversion to interface must succeed")
	upholdWallet3, ok := w3.(*uphold.Wallet)
	suite.Require().True(ok, "conversion to interface must succeed")
	upholdWallet4, ok := w4.(*uphold.Wallet)
	suite.Require().True(ok, "conversion to interface must succeed")

	anonCard1ID, err := upholdWallet1.CreateCardAddress("anonymous")
	suite.Require().NoError(err, "create anon card must not fail")
	anonCard1UUID := uuid.Must(uuid.FromString(anonCard1ID))

	anonCard2ID, err := upholdWallet2.CreateCardAddress("anonymous")
	suite.Require().NoError(err, "create anon card must not fail")
	anonCard2UUID := uuid.Must(uuid.FromString(anonCard2ID))

	w1ProviderID := upholdWallet1.GetWalletInfo().ProviderID
	w2ProviderID := upholdWallet2.GetWalletInfo().ProviderID
	w3ProviderID := upholdWallet3.GetWalletInfo().ProviderID

	amount := decimal.NewFromFloat(0)

	suite.claimCard(service, upholdWallet1, settlement, http.StatusOK, amount, noUUID())
	suite.claimCard(service, upholdWallet2, w1ProviderID, http.StatusOK, amount, &anonCard1UUID)
	suite.claimCard(service, upholdWallet2, w1ProviderID, http.StatusOK, amount, noUUID())
	suite.claimCard(service, upholdWallet3, w2ProviderID, http.StatusOK, amount, noUUID())
	suite.claimCard(service, upholdWallet4, w3ProviderID, http.StatusConflict, amount, noUUID())
	suite.claimCard(service, upholdWallet3, settlement, http.StatusOK, amount, &anonCard2UUID)
}

func (suite *WalletControllersTestSuite) claimCard(
	service *Service,
	w *uphold.Wallet,
	destination string,
	status int,
	amount decimal.Decimal,
	anonymousAddress *uuid.UUID,
) (*walletutils.Info, string) {
	signedTx, err := w.PrepareTransaction(*w.AltCurrency, amount, destination, "")
	suite.Require().NoError(err, "transaction must be signed client side")
	reqBody := LinkWalletRequest{
		SignedTx:         signedTx,
		AnonymousAddress: anonymousAddress,
	}
	body, err := json.Marshal(&reqBody)
	suite.Require().NoError(err, "unable to marshal claim body")
	handler := PostLinkWalletCompat(service)
	req, err := http.NewRequest("POST", "/v1/wallet/{paymentID}/claim", bytes.NewBuffer(body))
	suite.Require().NoError(err, "wallet claim request could not be created")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("paymentID", w.GetWalletInfo().ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(status, rr.Code, fmt.Sprintf("status is expected to match %d: %s", status, rr.Body.String()))
	linked, err := service.Datastore.GetWallet(uuid.Must(uuid.FromString(w.ID)))
	suite.Require().NoError(err, "retrieving the wallet did not cause an error")
	return linked, rr.Body.String()
}

func (suite *WalletControllersTestSuite) createBody(
	provider string,
	publicKey string,
	tx string,
) string {
	return `{"provider":"` + provider + `","publicKey":"` + publicKey + `","signedTx":"` + tx + `"}`
}

func (suite *WalletControllersTestSuite) NewWallet(service *Service, provider string) walletutils.Wallet {
	publicKey, privKey, err := httpsignature.GenerateEd25519Key(nil)
	publicKeyString := hex.EncodeToString(publicKey)

	bat := altcurrency.BAT
	info := walletutils.Info{
		ID:          uuid.NewV4().String(),
		PublicKey:   publicKeyString,
		Provider:    "uphold",
		AltCurrency: &bat,
	}
	wallet := &uphold.Wallet{
		Info:    info,
		PrivKey: privKey,
		PubKey:  publicKey,
	}

	reg, err := wallet.PrepareRegistration("Brave Browser Test Link")
	suite.Require().NoError(err, "unable to prepare transaction")

	createResp := suite.createWallet(
		service,
		suite.createBody(provider, publicKeyString, reg),
		http.StatusCreated,
		publicKey,
		privKey,
		true,
	)

	var returnedInfo walletutils.Info
	err = json.Unmarshal([]byte(createResp), &returnedInfo)
	suite.Require().NoError(err, "unable to create wallet")
	wallet.Info = returnedInfo
	return wallet
}

func (suite *WalletControllersTestSuite) TestPostCreateWallet() {
	pg, err := NewPostgres("", false)
	suite.Require().NoError(err, "Failed to get postgres connection")
	service := &Service{
		Datastore: pg,
	}

	publicKey, privKey, err := httpsignature.GenerateEd25519Key(nil)
	publicKeyString := hex.EncodeToString(publicKey)
	createBody := func(provider string, publicKey string) string {
		return `{"provider":"` + provider + `","publicKey":"` + publicKey + `"}`
	}
	badJSONBodyParse := suite.createWallet(
		service,
		``,
		http.StatusBadRequest,
		publicKey,
		privKey,
		true,
	)
	suite.Assert().JSONEq(`{
		"message": "Error unmarshalling body: unexpected end of JSON input",
		"code": 400
	}`, badJSONBodyParse, "should fail when parsing json")

	badFieldResponse := suite.createWallet(
		service,
		createBody("notaprovider", publicKeyString),
		http.StatusBadRequest,
		publicKey,
		privKey,
		true,
	)
	suite.Assert().JSONEq(`{
		"code": 400,
		"message": "Error validating request body",
		"data": {
			"validationErrors": {
				"provider": "notaprovider does not validate as in(brave|uphold)"
			}
		}
	}`, badFieldResponse, "field is not valid")

	// assume 403 is already covered
	// fail because of lacking signature presence
	notSignedResponse := suite.createWallet(
		service,
		createBody("brave", publicKeyString),
		http.StatusBadRequest,
		publicKey,
		privKey,
		false,
	)
	suite.Assert().Equal(`Bad Request
`, notSignedResponse, "not signed creation requests should fail")
	// body public key does not match signature public key
	notMatchingResponse := suite.createWallet(
		service,
		createBody("brave", uuid.NewV4().String()),
		http.StatusBadRequest,
		publicKey,
		privKey,
		true,
	)
	suite.Assert().JSONEq(`{
		"message": "Error validating request signature",
		"code": 400,
		"data": {
			"validationErrors": {
				"publicKey": "publicKey must match signature"
			}
		}
	}`, notMatchingResponse, "body should not match keyId")
	createResp := suite.createWallet(
		service,
		createBody("brave", publicKeyString),
		http.StatusCreated,
		publicKey,
		privKey,
		true,
	)

	var created walletutils.Info
	err = json.Unmarshal([]byte(createResp), &created)
	suite.Require().NoError(err, "unable to unmarshal response")

	getResp := suite.getWallet(service, uuid.Must(uuid.FromString(created.ID)), http.StatusOK)

	var gotten walletutils.Info
	err = json.Unmarshal([]byte(getResp), &gotten)
	suite.Require().NoError(err, "unable to unmarshal response")
	// gotten.PrivateKey = created.PrivateKey
	suite.Require().Equal(created, gotten, "the get and create return the same structure")
}

func (suite *WalletControllersTestSuite) getWallet(
	service *Service,
	paymentId uuid.UUID,
	code int,
) string {
	handler := GetWallet(service)

	req, err := http.NewRequest("GET", "/v1/wallet/"+paymentId.String(), nil)
	suite.Require().NoError(err, "a request should be created")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("paymentId", paymentId.String())
	joined := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(joined)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	suite.Require().Equal(code, rr.Code, "known status code should be sent")

	return rr.Body.String()
}

func (suite *WalletControllersTestSuite) createWallet(
	service *Service,
	body string,
	code int,
	publicKey httpsignature.Ed25519PubKey,
	privateKey ed25519.PrivateKey,
	shouldSign bool,
) string {
	handler := middleware.HTTPSignedOnly(service)(PostCreateWallet(service))

	bodyBuffer := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest("POST", "/v1/wallet", bodyBuffer)
	suite.Require().NoError(err, "a request should be created")

	if shouldSign {
		suite.SignRequest(
			req,
			publicKey,
			privateKey,
		)
	}

	rctx := chi.NewRouteContext()
	joined := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(joined)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Require().Equal(code, rr.Code, "known status code should be sent: "+rr.Body.String())

	return rr.Body.String()
}

func (suite *WalletControllersTestSuite) SignRequest(
	req *http.Request,
	publicKey httpsignature.Ed25519PubKey,
	privateKey ed25519.PrivateKey,
) {
	var s httpsignature.Signature
	s.Algorithm = httpsignature.ED25519
	s.KeyID = hex.EncodeToString(publicKey)
	s.Headers = []string{"digest", "(request-target)"}

	err := s.Sign(privateKey, crypto.Hash(0), req)
	suite.Require().NoError(err)
}
