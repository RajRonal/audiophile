package handlers

import (
	"audioPhile/database"
	"audioPhile/database/helper"
	"audioPhile/models"
	cloud "cloud.google.com/go/storage"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func CreateProductCategory(writer http.ResponseWriter, request *http.Request) {
	var productCategory models.ProductCategory
	var err error
	err = json.NewDecoder(request.Body).Decode(&productCategory)
	if err != nil {
		logrus.Error("SignUp: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	categoryId, err := helper.CreateProductCategory(productCategory)
	if err != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(writer).Encode(categoryId)
	if err != nil {
		logrus.Error("CreateProductCategory: Error in Encoding details %v", err)
		writer.WriteHeader(http.StatusBadRequest)
	}

}

func CreateProductInventory(writer http.ResponseWriter, request *http.Request) {
	var productQuantity models.ProductInventory
	categoryId := chi.URLParam(request, "categoryID")
	err := json.NewDecoder(request.Body).Decode(&productQuantity)
	if err != nil {
		logrus.Error("SignUp: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = helper.CreateProductInventory(productQuantity.Quantity, categoryId)
	if err != nil {
		logrus.Error("CreateProductInventory: Error in creating inventory%v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}
func CreateProduct(writer http.ResponseWriter, request *http.Request) {
	productDetails := &models.BulkProduct{}
	categoryId := chi.URLParam(request, "categoryID")
	categoryID, err := uuid.FromString(categoryId)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		return
	}

	inventoryId := chi.URLParam(request, "inventoryID")
	inventoryID, err := uuid.FromString(inventoryId)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		return
	}

	err = json.NewDecoder(request.Body).Decode(&productDetails)
	if err != nil {
		logrus.Error("CreateProduct: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err = helper.CreateProduct(productDetails.Products, categoryID, inventoryID, tx)
		if err != nil {
			logrus.Error("CreateProduct: Error in adding product %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return err
		}
		return err
	})
	if txErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	productId := chi.URLParam(request, "productID")
	productID, err := uuid.FromString(productId)
	if err != nil {
		logrus.Error("DeleteProduct:Error in conversion: %v", err)
		return
	}
	err = helper.DeleteProduct(productID)
	if err != nil {
		logrus.Error("DeleteProduct: Error in deleting product %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}
func UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	var product models.UpdateProduct
	productId := chi.URLParam(request, "productID")
	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		logrus.Error("UpdateProduct: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	productID, err := uuid.FromString(productId)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		return
	}
	err = helper.UpdateProduct(productID, product.ProductName, product.ProductDescription)
	if err != nil {
		logrus.Error("UpdateProduct: Error in updating product %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func CreateCoupon(writer http.ResponseWriter, request *http.Request) {
	var couponDetails models.Coupon
	err := json.NewDecoder(request.Body).Decode(&couponDetails)
	if err != nil {
		logrus.Error("CreateCoupon: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = helper.CreateCoupon(couponDetails)
	if err != nil {
		logrus.Error("CreateCoupon: Error in creating coupon%v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

func UpdateProductCategory(writer http.ResponseWriter, request *http.Request) {
	var categoryDetails models.ProductCategory
	categoryId := chi.URLParam(request, "categoryID")
	err := json.NewDecoder(request.Body).Decode(&categoryDetails)
	if err != nil {
		logrus.Error("UpdateProductCategory: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	categoryID, err := uuid.FromString(categoryId)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		return
	}
	err = helper.UpdateCategory(categoryID, categoryDetails)
	if err != nil {
		logrus.Error("UpdateProductCategory: Error in updating categoryDetails %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func DeleteProductCategory(writer http.ResponseWriter, request *http.Request) {
	categoryId := chi.URLParam(request, "categoryID")
	categoryID, err := uuid.FromString(categoryId)
	if err != nil {
		logrus.Error("DeleteProductCategory:Error in conversion: %v", err)
		return
	}
	err = helper.DeleteProductCategory(categoryID)
	if err != nil {
		logrus.Error("DeleteProductCategory: Error in deleting product category %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func UpdateProductInventory(writer http.ResponseWriter, request *http.Request) {
	var inventoryDetails models.ProductInventory
	categoryId := chi.URLParam(request, "categoryID")
	inventoryId := chi.URLParam(request, "inventoryID")
	err := json.NewDecoder(request.Body).Decode(&inventoryDetails)
	if err != nil {
		logrus.Error("UpdateProductCategory: Error in decoding json %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	categoryID, err := uuid.FromString(categoryId)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		return
	}
	inventoryID, err := uuid.FromString(inventoryId)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		return
	}
	err = helper.UpdateProductInventory(categoryID, inventoryID, inventoryDetails.Quantity)
	if err != nil {
		logrus.Error("UpdateProductCategory: Error in updating categoryDetails %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func DeleteProductInventory(writer http.ResponseWriter, request *http.Request) {
	inventoryId := chi.URLParam(request, "inventoryID")
	inventoryID, err := uuid.FromString(inventoryId)
	if err != nil {
		logrus.Error(" DeleteProductInventory:Error in conversion: %v", err)
		return
	}
	err = helper.DeleteProductInventory(inventoryID)
	if err != nil {
		logrus.Error(" DeleteProductInventory: Error in deleting product inventory %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}

func UploadProductImage(writer http.ResponseWriter, request *http.Request) {
	var err error
	productId := chi.URLParam(request, "productID")
	productID, err := uuid.FromString(productId)
	if err != nil {
		logrus.Error(" UploadProductImage:Error in conversion: %v", err)
		return
	}
	client := models.App{}

	client.Ctx = context.Background()
	credentialsFile := option.WithCredentialsJSON([]byte(os.Getenv("FIRE_KEY")))
	//fmt.Println(credentialsFile)
	app, err := firebase.NewApp(client.Ctx, nil, credentialsFile)
	if err != nil {
		logrus.Error(err)
		return
	}

	client.Client, err = app.Firestore(client.Ctx)
	if err != nil {
		logrus.Error(err)
		return
	}

	client.Storage, err = cloud.NewClient(client.Ctx, credentialsFile)
	if err != nil {
		logrus.Error(err)
		return
	}

	file, _, err := request.FormFile("image")
	err = request.ParseMultipartForm(10 << 20)
	if err != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()
	imagePath := strconv.Itoa(int(time.Now().Unix()))
	bucket := "audiophile-1606d.appspot.com"
	bucketStorage := client.Storage.Bucket(bucket).Object(imagePath).NewWriter(client.Ctx)

	_, err = io.Copy(bucketStorage, file)
	if err != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := bucketStorage.Close(); err != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	signedUrl := &cloud.SignedURLOptions{
		Scheme:  cloud.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}
	url, err := client.Storage.Bucket(bucket).SignedURL(imagePath, signedUrl)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Println(url)
	errs := json.NewEncoder(writer).Encode(url)
	if errs != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = helper.InsertImageDetails(imagePath, productID)
	if err != nil {
		logrus.Error(" UploadProductImage: Error in adding to image table: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

}
