package wallet

import (
	"database/sql"
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
	CreateWallet(wallet Wallet) error
	UpdateWallet(wallet Wallet) error
	DeleteWalletById(uid string) error
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

// CreateWalletHandler
//
// @Summary		Create  wallets
// @Description	Create  wallets
// @Tags		wallet
// @Accept		json
// @Produce		json
// @Param      	wallet body  Wallet true  "Wallet object"
// @Success		200	{object}	Wallet
// @Router      /api/v1/wallets [POST]
// @Failure		500	{object}	Err
func (h *Handler) CreateWalletHandler(c echo.Context) error {
	var wallet Wallet
	if err := c.Bind(&wallet); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	if err := h.store.CreateWallet(wallet); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, wallet)
}

// UpdateWalletHandler
//
// @Summary		UPDATE Wallets by id
// @Description	UPDATE Wallets by id
// @Tags		wallet
// @Accept		json
// @Param      	wallet body  Wallet true  "Wallet object"
// @Produce		json
// @Success		200	{object}	Wallet
// @Router		/api/v1/wallets [PUT]
// @Failure		500	{object}	Err
func (h *Handler) UpdateWalletHandler(c echo.Context) error {
	var wallet Wallet
	if err := c.Bind(&wallet); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	if err := h.store.UpdateWallet(wallet); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallet)
}

// DeleteWalletByIdHandler
//
// @Summary		Delete  wallets by user id
// @Description	Delete  wallets  by user id
// @Tags		wallet
// @Accept		json
// @Param      	id path  string true  "user id"
// @Produce		json
// @Success		200	{object}	Wallet
// @Router      /api/v1/user/{id}/wallets [DELETE]
// @Failure		500	{object}	Err
func (h *Handler) DeleteWalletByIdHandler(c echo.Context) error {
	user_id := c.Param("id")
	if err := h.store.DeleteWalletById(user_id); err != nil {
		if err != nil {
			// Handle specific error cases (e.g., user not found, database error)
			if err == sql.ErrNoRows {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"error": "User not found",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to delete user",
			})
		}
	}
	return c.JSON(http.StatusNoContent, map[string]interface{}{
		"Message": "Wallet should be deleted!",
	})

}
