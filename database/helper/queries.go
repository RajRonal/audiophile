package helper

import (
	"audioPhile/claims"
	"audioPhile/database"
	"audioPhile/models"
	"database/sql"
	"github.com/elgris/sqrl"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func CreateUser(user models.CreateUser) (string, error) {
	SQL := `INSERT INTO users (first_name, last_name, email, contact_number, user_name, password)
			VALUES($1, $2, $3, $4,$5,$6)
			returning user_id`
	var userID string
	err := database.DB.Get(&userID, SQL, user.FirstName, user.LastName, user.Email, user.ContactNumber, user.UserName, user.Password)
	if err != nil {
		logrus.Error("CreateUser:Error in creating user %v", err)
		return "", err
	}
	return userID, nil
}

func CreateRole(userRole string, userId string, tx *sqlx.Tx) error {
	SQL := `INSERT INTO roles(user_id, user_role) VALUES($1,$2)`
	_, err := tx.Exec(SQL, userId, userRole)
	if err != nil {
		logrus.Error("CreateRole:Error in assigning role %v", err)
		return err
	}
	return nil
}

func UserLogin(username string) (*models.AddLogin, error) {
	SQL := `SElECT user_id,password from users where user_name=$1`
	var pass models.AddLogin
	err := database.DB.Get(&pass, SQL, username)
	if err != nil && err != sql.ErrNoRows {
		logrus.Error("LoginUser:Error in logging in %v", err)
		return nil, err
	}
	return &pass, nil
}

func GetRole(userId string) (string, error) {
	Id, err := uuid.FromString(userId)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		return "", err
	}
	SQL := `SElECT user_role from roles where user_id=$1`
	var pass string
	err = database.DB.Get(&pass, SQL, Id)
	if err != nil {
		logrus.Error("GetRole: error in fetching data: %v", err)
		return "", err
	}
	return pass, nil
}

func CreateSession(userId string, expiredAt time.Time) (string, error) {
	SQL := `INSERT INTO sessions (user_id, expired_at ) 	
			VALUES ($1, $2)  
			returning session_id;`
	var sessID string
	err := database.DB.Get(&sessID, SQL, userId, expiredAt)
	if err != nil {
		logrus.Error("CreateSession:Error in creating Session %v", err)
		return "", err
	}
	return sessID, nil
}

func SessionExist(sessionID string) (bool, error) {
	var isExpired bool
	query := `SELECT count(*) > 0
			  FROM sessions 
			WHERE session_id=$1  and expired_at >now() and archived_at is null`
	checkSessionErr := database.DB.Get(&isExpired, query, sessionID)
	if checkSessionErr != nil {
		logrus.Error("SessionExist:Error in checking existence of session %v", checkSessionErr)
		return isExpired, checkSessionErr
	}
	return isExpired, nil
}

func GetContextData(request *http.Request) *claims.MapClaims {
	uc, ok := request.Context().Value(models.ClaimKey).(*claims.MapClaims)
	if !ok {
		logrus.Error("GetContextData: Error In Parsing Context")
		return nil
	}
	return uc
}

func CreateProductCategory(productType models.ProductCategory) (string, error) {
	var categoryId string
	SQL := `INSERT INTO product_categories (category_name, category_description) 	
			VALUES ($1, $2)  
			returning category_id`
	err := database.DB.Get(&categoryId, SQL, productType.CategoryName, productType.CategoryDescription)
	if err != nil {
		logrus.Error("CreateProductCategory: error in fetching data: %v", err)
		return "", err
	}
	return categoryId, nil
}

func CreateProductInventory(quantity int, categoryId string) error {
	categoryID, err := uuid.FromString(categoryId)
	if err != nil {
		logrus.Error("Error in conversion: %v", err)
		return err
	}

	SQL := `INSERT INTO product_inventory (quantity, category_id) 	
			VALUES ($1, $2)`
	_, err = database.DB.Exec(SQL, quantity, categoryID)
	if err != nil {
		logrus.Error("CreateProductInventor:Error in Creating Inventory %v", err)
		return err
	}

	return nil
}

func CreateProduct(product []models.Product, categoryId, inventoryId uuid.UUID, tx *sqlx.Tx) error {
	psql := sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)
	insertQuery := psql.Insert("product").Columns("product_name", "product_description", "regular_price", "discounted_price", "category_id", "inventory_id")
	for _, products := range product {

		insertQuery.Values(products.ProductName, products.ProductDescription, products.RegularPrice, products.DiscountedPrice, categoryId, inventoryId)
	}

	SQL, args, err := insertQuery.ToSql()
	if err != nil {
		logrus.Error("CreateProducts : Error in making the query")
		return err
	}

	_, err = tx.Exec(SQL, args...)
	if err != nil {
		logrus.Error("CreateProduct : Error in Adding product")
		return err
	}

	return nil
}

func DeleteProduct(productId uuid.UUID) error {
	SQL := `UPDATE product
	SET archived_at = now()
	WHERE product_id = $1 and archived_at is null`
	_, err := database.DB.Exec(SQL, productId)
	if err != nil {
		logrus.Error("DeleteProduct:error in deleting Product %v", err)
		return err
	}

	return nil
}

func UpdateProduct(productId uuid.UUID, productName, productDescription string) error {
	SQL := `UPDATE product
	SET product_name=$2,product_description=$3
	WHERE product_id=$1 and archived_at is null`
	_, err := database.DB.Exec(SQL, productId, productName, productDescription)
	if err != nil {
		logrus.Error("UpdateProduct: Error in updating product %v", err)
		return err
	}

	return nil
}

func CreateCoupon(couponDetails models.Coupon) error {
	SQL := `INSERT INTO discount (coupon_name, coupon_description,discount_percentage,discount_status) 	
			VALUES ($1, $2,$3,$4)`
	_, err := database.DB.Exec(SQL, couponDetails.CouponName, couponDetails.CouponDescription, couponDetails.DiscountPercentage, couponDetails.DiscountStatus)
	if err != nil {
		logrus.Error("CreateCoupon:Error in Creating Coupon %v", err)
		return err
	}

	return nil
}

func GetAllProducts(pageNo, taskSize int, searchProduct string) (models.PaginatedInventoryProductDetails, error) {
	var data models.PaginatedInventoryProductDetails
	SQL := `SELECT count(*) total_count,
										product.product_id,
										product_name,
										product_description,
										regular_price,
										discounted_price,
										   ARRAY_AGG(image_id) AS image_id
								 FROM product
										  JOIN image_details ON product.product_id = image_details.product_id
								 where product.product_name ILIKE '%' || $3 || '%'
								   AND archived_at IS NULL
                                    GROUP BY product.product_id
                                    LIMIT $1
                                    OFFSET $2
								
								
			
			`

	products := make([]models.InventoryProductDetails, 0)
	err := database.DB.Select(&products, SQL, taskSize, pageNo*taskSize, searchProduct)
	if err != nil {
		logrus.Error("SearchProduct: error in fetching Product: %v", err)
		return data, err
	}

	if len(products) == 0 {
		return data, err
	}

	data.TotalCount = products[0].TotalCount
	data.Details = products
	return data, err
}
func UpdateCategory(categoryId uuid.UUID, categoryDetails models.ProductCategory) error {
	SQL := `UPDATE product_categories
	SET category_name=$2,category_description=$3
	WHERE  category_id=$1 and archived_at is null`
	_, err := database.DB.Exec(SQL, categoryId, categoryDetails.CategoryName, categoryDetails.CategoryDescription)
	if err != nil {
		logrus.Error("UpdateCategory: Error in updating category %v", err)
		return err
	}

	return nil
}

func DeleteProductCategory(categoryId uuid.UUID) error {
	SQL := `UPDATE product_categories
	SET archived_at = now()
	WHERE category_id = $1 and archived_at is null`
	_, err := database.DB.Exec(SQL, categoryId)
	if err != nil {
		logrus.Error("DeleteProductCategory:error in deleting Product  category%v", err)
		return err
	}

	return nil
}

func UpdateProductInventory(categoryId, inventoryId uuid.UUID, quantity int) error {
	SQL := `UPDATE product_inventory
	SET quantity = $2
	WHERE category_id = $1 and  inventory_id =$3 and archived_at is null`
	_, err := database.DB.Exec(SQL, categoryId, quantity, inventoryId)
	if err != nil {
		logrus.Error("UpdateProductInventory:error in updating Product  inventory%v", err)
		return err
	}

	return nil
}

func DeleteProductInventory(inventoryId uuid.UUID) error {
	SQL := `UPDATE product_inventory
	SET archived_at = now()
	WHERE inventory_id = $1 and archived_at is null`
	_, err := database.DB.Exec(SQL, inventoryId)
	if err != nil {
		logrus.Error("DeleteProductInventory:error in deleting Product  inventory%v", err)
		return err
	}

	return nil
}

func InsertImageDetails(imageId string, productId uuid.UUID) error {
	SQL := `INSERT INTO image_details (image_id,product_id) 	
			VALUES ($1, $2)`
	_, err := database.DB.Exec(SQL, imageId, productId)
	if err != nil {
		logrus.Error("InsertImageDetails:Error in inserting into image table %v", err)
		return err
	}

	return nil
}

//func GetAllImageId(pageNo, taskSize int) (models.PaginatedImageDetails, error) {
//	var data models.PaginatedImageDetails
//	SQL := `WITH getImages AS (SELECT  count(*) over ()total_count,image_id, product_id
//			FROM image_details)
//
//			 SELECT  total_count,image_id,product_id from getImages
//			        LIMIT $1
//					OFFSET $2`
//
//	images := make([]models.ImageDetails, 0)
//	err := database.DB.Select(&images, SQL, taskSize, pageNo*taskSize)
//	if err != nil {
//		logrus.Error("SGetAllImageId: error in fetching Images: %v", err)
//		return data, err
//	}
//
//	if len(images) == 0 {
//		return data, err
//	}
//
//	data.TotalCount = images[0].TotalCount
//	data.Details = images
//	return data, err
//}

func AddToCart(sessionId uuid.UUID, productDetails models.CartProduct) error {
	SQL := `INSERT INTO cart_item (session_id,product_id,coupon_id,quantity) 	
			VALUES ($1, $2,$3,$4)`
	_, err := database.DB.Exec(SQL, sessionId, productDetails.ProductId, productDetails.CouponId, productDetails.Quantity)
	if err != nil {
		logrus.Error("AddToCart:Error in inserting into Cart %v", err)
		return err
	}

	return nil
}

func GetInventoryId(productId uuid.UUID) (uuid.UUID, error) {
	var inventoryId uuid.UUID
	SQL := `SELECT inventory_id FROM product 
            WHERE product_id=$1 AND archived_at IS NULL`
	err := database.DB.Get(&inventoryId, SQL, productId)
	if err != nil {
		logrus.Error("GetInventoryId:Error in getting inventory Id %v", err)
		return inventoryId, err
	}

	return inventoryId, nil
}

func GetCartItemInventoryQuantity(inventoryId uuid.UUID) (int, error) {
	var quantity int
	SQL := `SELECT quantity FROM product_inventory 
            WHERE inventory_id=$1 AND archived_at IS NULL`
	err := database.DB.Get(&quantity, SQL, inventoryId)
	if err != nil {
		logrus.Error("GetCartItemQuantity:Error in getting quantity of product %v", err)
		return quantity, err
	}

	return quantity, nil
}

func GetUserAddress(userId string) (bool, error) {
	var isExist bool
	SQL := `SELECT count(*) > 0
			  FROM user_address 
			WHERE user_id=$1 AND archived_at IS NULL `
	err := database.DB.Get(&isExist, SQL, userId)
	if err != nil {
		logrus.Error("GetUserAddressError in getting user address %v", err)
		return isExist, err
	}

	return isExist, nil
}

func AddUserAddress(userId uuid.UUID, addressDetails models.UserAddress) error {
	SQL := `INSERT INTO user_address (user_id,address_line_1,landmark,city,postal_code) 	
			VALUES ($1, $2,$3,$4,$5)`
	_, err := database.DB.Exec(SQL, userId, addressDetails.AddressLine1, addressDetails.Landmark, addressDetails.City, addressDetails.PostalCode)
	if err != nil {
		logrus.Error("AddUserAddress:Error in inserting user address %v", err)
		return err
	}

	return nil
}

func UpdateAddress(addressId uuid.UUID, addressDetails models.UserAddress) error {
	SQL := `UPDATE user_address
	SET address_line_1 = $1,landmark=$2,city=$3,postal_code=$4
	WHERE address_id = $5 AND archived_at IS NULL`
	_, err := database.DB.Exec(SQL, addressDetails.AddressLine1, addressDetails.Landmark, addressDetails.City, addressDetails.PostalCode, addressId)
	if err != nil {
		logrus.Error("UpdateAddress:error in updating user address%v", err)
		return err
	}

	return nil
}

func DeleteAddress(addressId uuid.UUID) error {
	SQL := `UPDATE user_address
	SET archived_at = now()
	WHERE address_id = $1 and archived_at is null`
	_, err := database.DB.Exec(SQL, addressId)
	if err != nil {
		logrus.Error("DeleteAddress:error in deleting user Address%v", err)
		return err
	}

	return nil
}

func AddPaymentDetails(userId uuid.UUID, paymentDetails models.PaymentDetails) error {
	SQL := `INSERT INTO payment (user_id,payment_type) 	
			VALUES ($1, $2)`
	_, err := database.DB.Exec(SQL, userId, paymentDetails.PaymentType)
	if err != nil {
		logrus.Error("AddPaymentDetails:Error in inserting user Payment Details %v", err)
		return err
	}

	return nil

}

func GetProductQuantity(productId uuid.UUID) (int, error) {
	var quantity int
	SQL := `SELECT quantity FROM cart_item 
            WHERE product_id=$1 AND archived_at IS NULL`
	err := database.DB.Get(&quantity, SQL, productId)
	if err != nil {
		logrus.Error("GetProductQuantity:Error in getting quantity of product %v", err)
		return quantity, err
	}

	return quantity, nil

}

func UpdateInventoryQuantity(inventoryId uuid.UUID, quantity int) error {
	SQL := `UPDATE product_inventory
	SET quantity = $1
	WHERE inventory_id =$2 and archived_at is null`
	_, err := database.DB.Exec(SQL, quantity, inventoryId)
	if err != nil {
		logrus.Error("UpdateInventoryQuantity:error in updating Product  inventory quantity%v", err)
		return err
	}

	return nil
}

func DeleteCart(sessionId uuid.UUID, tx *sqlx.Tx) error {
	SQL := `UPDATE cart_item
	SET archived_at = now()
	WHERE session_id = $1 and archived_at is null`
	_, err := tx.Exec(SQL, sessionId)
	if err != nil {
		logrus.Error("DeleteCar:error in deleting user Cart%v", err)
		return err
	}

	return nil
}

func AddOrderDetails(userId, paymentId uuid.UUID, total float64) error {
	SQL := `INSERT INTO order_details (user_id,total,payment_id) 	
			VALUES ($1, $2,$3)`
	_, err := database.DB.Exec(SQL, userId, total, paymentId)
	if err != nil {
		logrus.Error("AddOrderDetails:Error in inserting user order Details %v", err)
		return err
	}

	return nil

}

func ShowOrderDetails(userId uuid.UUID, tc *sqlx.Tx) (models.OrderDetails, error) {
	var orderDetails models.OrderDetails
	SQL := `SELECT order_id,user_id,total,payment_id FROM order_details 
            WHERE user_id=$1`
	err := tc.Get(&orderDetails, SQL, userId)
	if err != nil {
		logrus.Error("ShowOrderDetails:Error in order details %v", err)
		return orderDetails, err
	}

	return orderDetails, nil
}

func GetProductPrice(productId uuid.UUID) (float64, error) {
	var price float64
	SQL := `SELECT regular_price FROM product 
            WHERE product_id=$1`
	err := database.DB.Get(&price, SQL, productId)
	if err != nil {
		logrus.Error("GetProductPrice:Error in getting price of product %v", err)
		return price, err
	}

	return price, nil
}

func DeleteSession(sessionID string) error {
	currentTime := time.Now()
	SQL := `UPDATE sessions
			  SET archived_at= $1,
			      expired_at= now()
			  WHERE session_id= $2`
	_, err := database.DB.Exec(SQL, currentTime, sessionID)
	if err != nil {
		logrus.Error("DeleteSession: error in Deleting Session: %v", err)
		return err
	}
	return err
}

func CreateCartSession(userId string) (string, error) {
	SQL := `INSERT INTO cart_sessions (user_id) 	
			VALUES ($1)  
			returning session_id;`
	var sessID string
	err := database.DB.Get(&sessID, SQL, userId)
	if err != nil {
		logrus.Error("CreateCartSession:Error in creating Session %v", err)
		return "", err
	}
	return sessID, nil
}

func GetCartSessionId(productId uuid.UUID) (uuid.UUID, error) {
	var sessionId uuid.UUID
	SQL := `SELECT session_id FROM cart_item
            WHERE product_id=$1`
	err := database.DB.Get(&sessionId, SQL, productId)
	if err != nil {
		logrus.Error(" GetCartSessionId:Error in getting sessionId of cart %v", err)
		return sessionId, err
	}

	return sessionId, nil

}

func DeleteUserAccount(userId uuid.UUID) error {
	SQL := `UPDATE users
			  SET archived_at= now()
			  WHERE user_id= $1`
	_, err := database.DB.Exec(SQL, userId)
	if err != nil {
		logrus.Error("DeleteUserAccount: error in Deleting user: %v", err)
		return err
	}
	return err
}
