package wallet

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Storer interface {
	Wallets() ([]Wallet, error)
	QueryWalletsWithType(wt string) ([]Wallet, error)
	QueryWalletsByUserId(uid string) ([]Wallet, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

// WalletHandler
//
// @Summary		Get all wallets
// @Description	Get all wallets
// @Tags		wallet
// @Accept		json
// @Param       wallet_type query string  ture "Wallet type  fillter"
// @Produce		json
// @Success		200	{object}	Wallet
// @Router		/api/v1/wallets [get]
// @Failure		500	{object}	Err
func (h *Handler) WalletHandler(c echo.Context) error {
	walletType := c.QueryParam("wallet_type")
	var wallets []Wallet
	var err error
	if walletType != "" {

		wallets, err = h.store.QueryWalletsWithType(walletType)
	} else {

		wallets, err = h.store.Wallets()

	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, wallets)
}

// GetWalletByIdHandler
//
// @Summary		Get  wallets by user id
// @Description	Get  wallets  by user id
// @Tags		wallet
// @Accept		json
// @Param      	id path  string true  "user id"
// @Produce		json
// @Success		200	{object}	Wallet
// @Router      /api/v1/user/{id}/wallets [get]
// @Failure		500	{object}	Err
func (h *Handler) GetWalletByIdHandler(c echo.Context) error {
	user_id := c.Param("id")
	wallets, err := h.store.QueryWalletsByUserId(user_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)

}
