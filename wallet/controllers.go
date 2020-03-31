package wallet

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/brave-intl/bat-go/middleware"
	"github.com/brave-intl/bat-go/utils/altcurrency"
	"github.com/brave-intl/bat-go/utils/handlers"
	"github.com/brave-intl/bat-go/utils/httpsignature"
	"github.com/brave-intl/bat-go/utils/requestutils"
	walletutils "github.com/brave-intl/bat-go/utils/wallet"
	"github.com/brave-intl/bat-go/utils/wallet/provider/uphold"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

// Router for suggestions endpoints
func Router(service *Service) chi.Router {
	r := chi.NewRouter()
	PostLinkWalletCompatHandler := middleware.HTTPSignedOnly(service)(middleware.InstrumentHandler("PostLinkWalletCompat", PostLinkWalletCompat(service)))
	if os.Getenv("ENV") != "local" {
		// only link from ledger
		PostLinkWalletCompatHandler = middleware.SimpleTokenAuthorizedOnly(PostLinkWalletCompatHandler)
	}
	r.Method("POST", "/{paymentId}/claim", PostLinkWalletCompatHandler)
	r.Method("GET", "/{paymentId}", middleware.InstrumentHandler("GetWallet", GetWallet(service)))
	r.Method("POST", "/", middleware.HTTPSignedOnly(service)(middleware.InstrumentHandler("PostCreateWallet", PostCreateWallet(service))))
	return r
}

// LookupPublicKey based on the HTTP signing keyID, which in our case is the walletID
func (service *Service) LookupPublicKey(ctx context.Context, keyID string) (*httpsignature.Verifier, error) {
	var publicKey httpsignature.Ed25519PubKey
	// hex encoded public key
	publicKey, err := hex.DecodeString(keyID)
	if err != nil {
		return nil, err
	}
	tmp := httpsignature.Verifier(publicKey)
	return &tmp, nil
}

// LinkWalletRequest holds the data necessary to update a wallet with an anonymous address
type LinkWalletRequest struct {
	SignedTx         string     `json:"signedTx"`
	AnonymousAddress *uuid.UUID `json:"anonymousAddress"`
}

// PostLinkWalletCompat links wallets using provided ids
func PostLinkWalletCompat(service *Service) handlers.AppHandler {
	return handlers.AppHandler(func(w http.ResponseWriter, r *http.Request) *handlers.AppError {
		paymentIDString := chi.URLParam(r, "paymentID")
		paymentID, err := uuid.FromString(paymentIDString)
		if err != nil {
			return handlers.ValidationError("url parameter", map[string]string{
				"paymentID": "must be a valid uuidv4",
			})
		}

		var body LinkWalletRequest
		err = requestutils.ReadJSON(r.Body, &body)
		if err != nil {
			return handlers.ValidationError("request body", map[string]string{
				"body": "unable to ready body",
			})
		}
		_, err = govalidator.ValidateStruct(body)
		if err != nil {
			return handlers.WrapValidationError(err)
		}

		wallet, err := service.ReadableDatastore().GetWallet(paymentID)
		if err != nil {
			return handlers.WrapError(err, "Error finding wallet", http.StatusNotFound)
		}
		err = service.LinkWallet(r.Context(), wallet, body.SignedTx, body.AnonymousAddress)
		if err != nil {
			return handlers.WrapError(err, "error linking wallet", http.StatusBadRequest)
		}

		return handlers.RenderContent(r.Context(), wallet, w, http.StatusOK)
	})
}

// PostCreateWalletResponse includes a ClaimID which can later be used to check the status of the claim
type PostCreateWalletResponse struct {
	Wallet     *walletutils.Info `json:"wallet"`
	PrivateKey *string           `json:"privateKey"`
}

// GetWalletResponse gets wallet info tied to a wallet id
type GetWalletResponse struct {
	Wallet *walletutils.Info `json:"wallet"`
}

// PostCreateWalletRequest has possible inputs for the new wallet
type PostCreateWalletRequest struct {
	Provider  string `json:"provider" valid:"in(brave|uphold)"`
	PublicKey string `json:"publicKey" valid:"required"`
	SignedTx  string `json:"signedTx" valid:"-"`
}

// PostCreateWallet creates a wallet
func PostCreateWallet(service *Service) handlers.AppHandler {
	return handlers.AppHandler(func(w http.ResponseWriter, r *http.Request) *handlers.AppError {
		limit := int64(1024 * 1024 * 10) // 10MiB
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, limit))
		if err != nil {
			return handlers.WrapError(err, "Error reading body", http.StatusBadRequest)
		}

		var req PostCreateWalletRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			return handlers.WrapError(err, "Error unmarshalling body", http.StatusBadRequest)
		}
		_, err = govalidator.ValidateStruct(req)
		if err != nil {
			return handlers.WrapValidationError(err)
		}

		publicKey, err := middleware.GetKeyID(r.Context())
		if err != nil {
			return handlers.WrapError(err, "unable to look up http signature info", http.StatusBadRequest)
		}
		if req.PublicKey != publicKey {
			return handlers.ValidationError("request signature", map[string]string{
				"publicKey": "publicKey must match signature",
			})
		}

		info, err := CreateWallet(req)
		if err != nil {
			return handlers.WrapError(err, "unable to save wallet", http.StatusServiceUnavailable)
		}
		err = service.Datastore.InsertWallet(&info)
		if err != nil {
			return handlers.WrapError(err, "unable to save wallet", http.StatusServiceUnavailable)
		}

		return handlers.RenderContent(r.Context(), info, w, http.StatusCreated)
	})
}

// GetWallet retrieves wallet information
func GetWallet(service *Service) handlers.AppHandler {
	return handlers.AppHandler(func(w http.ResponseWriter, r *http.Request) *handlers.AppError {
		paymentIDParam := chi.URLParam(r, "paymentId")
		paymentID, err := uuid.FromString(paymentIDParam)

		if err != nil {
			return handlers.ValidationError("request url parameter", map[string]string{
				"paymentId": "paymentId '" + paymentIDParam + "' is not supported",
			})
		}

		info, err := service.Datastore.GetWallet(paymentID)
		if err != nil {
			return handlers.WrapError(err, "Error getting wallet", http.StatusNotFound)
		}

		// just doing this until another way to track
		if info.AltCurrency == nil {
			tmp := altcurrency.BAT
			info.AltCurrency = &tmp
		}

		return handlers.RenderContent(r.Context(), info, w, http.StatusOK)
	})
}

// CreateWallet creates a new set of wallet info
func CreateWallet(req PostCreateWalletRequest) (walletutils.Info, error) {
	provider := req.Provider // client
	publicKey := req.PublicKey

	var info walletutils.Info
	info.ID = uuid.NewV4().String()
	info.Provider = provider
	{
		tmp := altcurrency.BAT
		info.AltCurrency = &tmp
	}

	info.PublicKey = publicKey

	if req.Provider == "uphold" {
		if req.SignedTx != "" {
			wallet := uphold.Wallet{
				Info:    info,
				PrivKey: ed25519.PrivateKey{},
				PubKey:  httpsignature.Ed25519PubKey([]byte(publicKey)),
			}
			err := wallet.SubmitRegistration(req.SignedTx)
			if err != nil {
				return info, err
			}
			info.ProviderID = wallet.GetWalletInfo().ProviderID
		}
	}
	return info, nil
}
