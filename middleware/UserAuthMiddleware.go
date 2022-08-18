package middleware

import (
	"audioPhile/database/helper"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

func UserAddressAndQuantityValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := helper.GetContextData(request)
		if ctx == nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err := helper.GetUserAddress(ctx.ID)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		productID := chi.URLParam(request, "productID")
		productId, err := uuid.FromString(productID)
		if err != nil {
			logrus.Error("ProductQuantityValidateMiddleware:Error in conversion: %v", err)
			return
		}

		CartItemQuantity, err := helper.GetProductQuantity(productId)
		if err != nil {
			logrus.Error("ProductQuantityValidateMiddleware:Error in fetching cart quantity: %v", err)
			return
		}

		inventoryId, err := helper.GetInventoryId(productId)
		if err != nil {
			logrus.Error("ProductQuantityValidateMiddleware:Error in getting inventory Id: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		productQuantity, err := helper.GetCartItemInventoryQuantity(inventoryId)
		if err != nil {
			logrus.Error(" ProductQuantityValidateMiddleware:Error in getting product quantity: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if CartItemQuantity > productQuantity {
			logrus.Error(" ProductQuantityValidateMiddleware: Please enter a valid quantity %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
