package wallet

import (
	raven "github.com/getsentry/raven-go"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"

	"github.com/brave-intl/bat-go/wallet"
	// needed for magic migration
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Datastore holds the interface for the wallet datastore
type Datastore interface {
	SetAnonymousAddress(ID uuid.UUID, anonymousAddress uuid.UUID) error
	TxSetAnonymousAddress(tx *sqlx.Tx, ID uuid.UUID, anonymousAddress uuid.UUID) error
	TxGetByProviderAddress(tx *sqlx.Tx, providerAddress uuid.UUID) (*[]wallet.Info, error)
	LinkWallet(wallet uuid.UUID, anonymousAddress uuid.UUID) error
	GetWallet(id uuid.UUID) (*wallet.Info, error)
}

// ReadOnlyDatastore includes all database methods that can be made with a read only db connection
type ReadOnlyDatastore interface {
	// GetWallet by ID
	GetWallet(id uuid.UUID) (*wallet.Info, error)
}

// Postgres is a Datastore wrapper around a postgres database
type Postgres struct {
	*sqlx.DB
}

// GetWallet by ID
func (pg *Postgres) GetWallet(ID uuid.UUID) (*wallet.Info, error) {
	statement := "select * from wallets where id = $1"
	wallets := []wallet.Info{}
	err := pg.DB.Select(&wallets, statement, ID)
	if err != nil {
		return nil, err
	}

	if len(wallets) > 0 {
		return &wallets[0], nil
	}

	return nil, nil
}

// SetAnonymousAddress sets the anon addresses of the provided wallets
func (pg *Postgres) SetAnonymousAddress(ID uuid.UUID, anonymousAddress uuid.UUID) error {
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	err = pg.TxSetAnonymousAddress(tx, ID, anonymousAddress)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// TxSetAnonymousAddress pass a tx to set the anonymous address
func (pg *Postgres) TxSetAnonymousAddress(tx *sqlx.Tx, ID uuid.UUID, anonymousAddress uuid.UUID) error {
	statement := `
	UPDATE wallets
	SET anonymous_address = $2
	WHERE id = $1
	`
	_, err := tx.Exec(statement, ID.String(), anonymousAddress.String())
	return err
}

// TxGetByProviderAddress gets a wallet by a provider address
func (pg *Postgres) TxGetByProviderAddress(tx *sqlx.Tx, providerAddress uuid.UUID) (*[]wallet.Info, error) {
	statement := `
	SELECT wallets
	WHERE provider_address = $1
	`
	var wallets []wallet.Info
	err := tx.Select(&wallets, statement, providerAddress.String())
	if err != nil {
		return nil, err
	}
	return &wallets, nil
}

// LinkWallet links a wallet together
func (pg *Postgres) LinkWallet(id uuid.UUID, providerLinkingID uuid.UUID, anonymousAddress uuid.UUID) error {
	tx, err := pg.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	walletsMatchingProviderAddress, err := pg.TxGetByProviderAddress(tx, providerLinkingID)
	if err != nil {
		return errors.Wrap(err, "error looking up wallets by provider id")
	}
	walletLinkedLength := len(*walletsMatchingProviderAddress)
	if walletLinkedLength >= 3 {
		if walletLinkedLength > 3 {
			raven.CaptureMessage("to many cards linked", nil)
		}
		return errors.Wrap(err, "unable to add too many wallets to a single user")
	}
	err = pg.TxSetAnonymousAddress(tx, id, anonymousAddress)
	if err != nil {
		return errors.Wrap(err, "unable to set an anonymous address")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
