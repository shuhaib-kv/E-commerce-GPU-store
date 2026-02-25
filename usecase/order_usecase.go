package usecase

import (
	"context"
	"errors"
	"ga/domain"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderUsecase struct {
	orderRepo    domain.OrderRepository
	cartRepo     domain.CartRepository
	couponRepo   domain.CouponRepository
	walletRepo   domain.WalletRepository
	discountRepo domain.DiscountRepository
	productRepo  domain.ProductRepository
}

func NewOrderUsecase(
	or domain.OrderRepository,
	cr domain.CartRepository,
	couponRepo domain.CouponRepository,
	wr domain.WalletRepository,
	dr domain.DiscountRepository,
	pr domain.ProductRepository,
) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: or, cartRepo: cr, couponRepo: couponRepo,
		walletRepo: wr, discountRepo: dr, productRepo: pr,
	}
}

func (uc *OrderUsecase) CreateOrder(ctx context.Context, userID, addressID primitive.ObjectID, paymentMethod, couponCode string, applyWallet bool) (string, uint, error) {
	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return "", 0, errors.New("cart not found")
	}

	items, err := uc.cartRepo.FindCartProducts(ctx, cart.ID)
	if err != nil || len(items) == 0 {
		return "", 0, errors.New("cart is empty")
	}

	// Calculate total with discounts
	var totalAmount uint
	for _, item := range items {
		product, err := uc.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			totalAmount += item.ProductPrice
			continue
		}
		price := product.Price * item.Quantity
		if !product.DiscountID.IsZero() {
			discount, err := uc.discountRepo.FindByID(ctx, product.DiscountID)
			if err == nil {
				discountAmt := uint(float64(price) * float64(discount.DiscountPercentage) / 100.0)
				price -= discountAmt
			}
		}
		totalAmount += price
	}

	// Apply coupon
	if couponCode != "" {
		coupon, err := uc.couponRepo.FindByCode(ctx, couponCode)
		if err != nil {
			return "", 0, errors.New("invalid coupon code")
		}
		if time.Now().After(coupon.ExpiryDate) {
			return "", 0, errors.New("coupon has expired")
		}
		couponDiscount := uint(float64(totalAmount) * float64(coupon.CouponPercentage) / 100.0)
		totalAmount -= couponDiscount
	}

	// Apply wallet
	if applyWallet {
		wallet, err := uc.walletRepo.FindByUserID(ctx, userID)
		if err == nil && wallet.Balance > 0 {
			if wallet.Balance >= totalAmount {
				uc.walletRepo.UpdateBalance(ctx, userID, wallet.Balance-totalAmount)
				uc.walletRepo.AddHistory(ctx, &domain.WalletHistory{UserID: userID, Debit: totalAmount})
				totalAmount = 0
			} else {
				totalAmount -= wallet.Balance
				uc.walletRepo.AddHistory(ctx, &domain.WalletHistory{UserID: userID, Debit: wallet.Balance})
				uc.walletRepo.UpdateBalance(ctx, userID, 0)
			}
		}
	}

	orderID := uuid.New().String()
	order := &domain.Order{
		UserID:               userID,
		AddressID:            addressID,
		OrderID:              orderID,
		PaymentMethod:        paymentMethod,
		TotalAmount:          totalAmount,
		Status:               true,
		PaymentStatus:        paymentMethod == "cod",
		ExpectedDeliveryDate: time.Now().AddDate(0, 0, 12),
		CreatedAt:            time.Now(),
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return "", 0, err
	}

	for _, item := range items {
		uc.orderRepo.CreateItem(ctx, &domain.OrderItem{
			OrderID:     orderID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.ProductPrice,
		})
	}

	uc.cartRepo.DeleteCartProducts(ctx, cart.ID)

	return orderID, totalAmount, nil
}

func (uc *OrderUsecase) ListOrders(ctx context.Context, userID primitive.ObjectID) ([]map[string]interface{}, error) {
	orders, err := uc.orderRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, order := range orders {
		items, _ := uc.orderRepo.FindItemsByOrderID(ctx, order.OrderID)
		result = append(result, map[string]interface{}{
			"order":  order,
			"items":  items,
		})
	}
	return result, nil
}

func (uc *OrderUsecase) CancelOrder(ctx context.Context, userID primitive.ObjectID, orderID string) error {
	order, err := uc.orderRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if err := uc.orderRepo.UpdateStatus(ctx, orderID, false, order.PaymentStatus); err != nil {
		return err
	}

	// Refund to wallet if paid
	if order.PaymentStatus {
		wallet, err := uc.walletRepo.FindByUserID(ctx, userID)
		if err == nil {
			uc.walletRepo.UpdateBalance(ctx, userID, wallet.Balance+order.TotalAmount)
			uc.walletRepo.AddHistory(ctx, &domain.WalletHistory{
				UserID: userID,
				Credit: order.TotalAmount,
			})
		}
	}
	return nil
}

// Admin
func (uc *OrderUsecase) ListAllOrders(ctx context.Context, filter map[string]interface{}, page, pageSize int) ([]domain.Order, int64, error) {
	return uc.orderRepo.FindAll(ctx, filter, page, pageSize)
}

func (uc *OrderUsecase) EditOrder(ctx context.Context, orderID string, status, paymentStatus bool) error {
	return uc.orderRepo.UpdateStatus(ctx, orderID, status, paymentStatus)
}
