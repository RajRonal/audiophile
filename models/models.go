package models

import (
	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	"context"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

type CreateUser struct {
	FirstName     string `json:"firstName" db:"first_name"`
	LastName      string `json:"lastName" db:"last_name"`
	Email         string `json:"email" db:"email"`
	ContactNumber string `json:"contactNumber" db:"contact_number"`
	UserName      string `json:"userName" db:"user_name"`
	Password      string `json:"password" db:"password"`
}

type UserRole string
type ClaimsKey string

//type PaymentTypes string

const (
	UserRoleAdmin UserRole  = "admin"
	UserRoleUser  UserRole  = "user"
	ClaimKey      ClaimsKey = "claim"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type AddLogin struct {
	UserId   string `json:"userId" db:"user_id"`
	Password string `json:"password" db:"password"`
}
type ProductCategory struct {
	CategoryName        string `json:"categoryName" db:"category_name"`
	CategoryDescription string `json:"categoryDescription" db:"category_description"`
}
type ProductInventory struct {
	Quantity int `json:"quantity" db:"quantity"`
}
type Product struct {
	ProductName        string  `json:"productName" db:"product_name"`
	ProductDescription string  `json:"productDescription" db:"product_description"`
	RegularPrice       float64 `json:"regularPrice" db:"regular_price"`
	DiscountedPrice    float64 `json:"discountedPrice" db:"discounted_price"`
}
type BulkProduct struct {
	Products []Product `json:"products"`
}
type UpdateProduct struct {
	ProductName        string  `json:"productName" db:"product_name"`
	ProductDescription string  `json:"productDescription" db:"product_description"`
	RegularPrice       float64 `json:"regularPrice" db:"regular_price"`
	DiscountedPrice    float64 `json:"discountedPrice" db:"discounted_price"`
}

type Coupon struct {
	CouponName         string `json:"couponName" db:"coupon_name"`
	CouponDescription  string `json:"couponDescription" db:"coupon_description"`
	DiscountPercentage int    `json:"discountPercentage" db:"discount_percentage"`
	DiscountStatus     bool   `json:"discountStatus" db:"discount_status"`
}

type SearchProducts struct {
	ProductId          string  `json:"productId" db:"product_id"`
	ProductName        string  `json:"productName" db:"product_name"`
	ProductDescription string  `json:"productDescription" db:"product_description"`
	RegularPrice       float64 `json:"regularPrice" db:"regular_price"`
	DiscountedPrice    float64 `json:"discountedPrice" db:"discounted_price"`
	TotalCount         int     `json:"-" db:"total_count"`
}

type PaginatedProducts struct {
	Products   []SearchProducts `json:"products"`
	TotalCount int              `json:"totalCount"`
}

type App struct {
	Ctx     context.Context
	Client  *firestore.Client
	Storage *cloud.Client
}

type InventoryProductDetails struct {
	ProductId          string         `json:"productId"db:"product_id"`
	ProductName        string         `json:"productName" db:"product_name"`
	ProductDescription string         `json:"productDescription" db:"product_description"`
	RegularPrice       float64        `json:"regularPrice" db:"regular_price"`
	DiscountedPrice    float64        `json:"discountedPrice" db:"discounted_price"`
	ImageId            pq.StringArray `json:"imageId" db:"image_id"`
	TotalCount         int            `json:"-" db:"total_count"`
}

//type Image struct {
//	Images string `json:"images"`
//}
type PaginatedInventoryProductDetails struct {
	Details    []InventoryProductDetails `json:"details"`
	TotalCount int                       `json:"totalCount"`
}

type CartProduct struct {
	ProductId uuid.UUID `json:"productId" db:"product_id"`
	CouponId  uuid.UUID `json:"couponId" db:"coupon_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
}

type UserAddress struct {
	AddressLine1 string `json:"addressLine1" db:"address_line_1"`
	Landmark     string `json:"landmark" db:"landmark"`
	City         string `json:"city" db:"city"`
	PostalCode   int    `json:"postalCode" db:"postal_code"`
}

type PaymentDetails struct {
	PaymentType string `json:"paymentType" db:"payment_type"`
}

type UpdateInventory struct {
	ProductId string `json:"productId" db:"product_id"`
	Quantity  int    `json:"quantity" db:"quantity"`
}
type PaginatedUpdateInventory struct {
	Details []UpdateInventory `json:"details"`
}

type OrderDetails struct {
	OrderId   uuid.UUID `json:"orderId"db:"order_id"`
	UserId    uuid.UUID `json:"userId"db:"user_id"`
	Total     float64   `json:"total" db:"total"`
	PaymentId uuid.UUID `json:"paymentId" db:"payment_id"`
	//CreatedAt timestamp.Timestamp `json:"createdAt" db:"created_at"`
}
