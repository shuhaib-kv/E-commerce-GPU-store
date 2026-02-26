package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collectionNames = []string{
	"users", "addresses", "admins", "categories", "products",
	"product_attribute_definitions", "orders", "order_items",
	"carts", "cart_products", "coupons", "discounts",
	"wallets", "wallet_history", "payments",
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

// Drop removes all data from every collection.
func Drop(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, name := range collectionNames {
		if _, err := db.Collection(name).DeleteMany(ctx, bson.M{}); err != nil {
			log.Printf("[DROP] Failed to clear %s: %v", name, err)
		} else {
			fmt.Printf("[DROP] Cleared collection: %s\n", name)
		}
	}
	fmt.Println("[DROP] All collections cleared.")
}

// Seed inserts sample data into all collections.
func Seed(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()

	// --- Admins ---
	admin1ID := primitive.NewObjectID()
	admins := []interface{}{
		bson.M{
			"_id": admin1ID, "name": "Super Admin",
			"email": "admin@gpu.store", "password": hashPassword("admin123"),
		},
	}
	insertMany(ctx, db, "admins", admins)

	// --- Categories ---
	catGaming := primitive.NewObjectID()
	catWorkstation := primitive.NewObjectID()
	catMining := primitive.NewObjectID()
	catBudget := primitive.NewObjectID()
	categories := []interface{}{
		bson.M{"_id": catGaming, "name": "Gaming"},
		bson.M{"_id": catWorkstation, "name": "Workstation"},
		bson.M{"_id": catMining, "name": "Mining"},
		bson.M{"_id": catBudget, "name": "Budget"},
	}
	insertMany(ctx, db, "categories", categories)

	// --- Discounts ---
	disc10 := primitive.NewObjectID()
	disc15 := primitive.NewObjectID()
	disc20 := primitive.NewObjectID()
	discounts := []interface{}{
		bson.M{"_id": disc10, "discount_name": "Summer Sale", "discount_percentage": uint(10)},
		bson.M{"_id": disc15, "discount_name": "Flash Deal", "discount_percentage": uint(15)},
		bson.M{"_id": disc20, "discount_name": "Clearance", "discount_percentage": uint(20)},
	}
	insertMany(ctx, db, "discounts", discounts)

	// --- Products ---
	prod1 := primitive.NewObjectID()
	prod2 := primitive.NewObjectID()
	prod3 := primitive.NewObjectID()
	prod4 := primitive.NewObjectID()
	prod5 := primitive.NewObjectID()
	prod6 := primitive.NewObjectID()
	prod7 := primitive.NewObjectID()
	prod8 := primitive.NewObjectID()
	prod9 := primitive.NewObjectID()
	prod10 := primitive.NewObjectID()
	products := []interface{}{
		bson.M{
			"_id": prod1, "name": "NVIDIA GeForce RTX 4090", "price": uint(159999), "model_no": uint(4090),
			"stock": uint(25), "category_id": catGaming, "brand": "NVIDIA", "discount_id": disc10,
			"description": "The ultimate gaming GPU with 24GB GDDR6X, ray tracing, and DLSS 3.",
			"specifications": bson.M{"vram": "24GB GDDR6X", "cuda_cores": "16384", "boost_clock": "2520 MHz", "tdp": "450W", "interface": "PCIe 4.0 x16"},
			"image1": "https://placehold.co/400x300/1a1a2e/10b981?text=RTX+4090", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod2, "name": "NVIDIA GeForce RTX 4080 Super", "price": uint(109999), "model_no": uint(4080),
			"stock": uint(40), "category_id": catGaming, "brand": "NVIDIA", "discount_id": disc15,
			"description": "High-performance gaming with 16GB GDDR6X and Ada Lovelace architecture.",
			"specifications": bson.M{"vram": "16GB GDDR6X", "cuda_cores": "10240", "boost_clock": "2550 MHz", "tdp": "320W", "interface": "PCIe 4.0 x16"},
			"image1": "https://placehold.co/400x300/1a1a2e/10b981?text=RTX+4080S", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod3, "name": "AMD Radeon RX 7900 XTX", "price": uint(94999), "model_no": uint(7900),
			"stock": uint(35), "category_id": catGaming, "brand": "AMD",
			"description": "AMD's flagship GPU with 24GB GDDR6 for 4K gaming.",
			"specifications": bson.M{"vram": "24GB GDDR6", "stream_processors": "6144", "boost_clock": "2500 MHz", "tdp": "355W", "interface": "PCIe 4.0 x16"},
			"image1": "https://placehold.co/400x300/1a1a2e/e74c3c?text=RX+7900+XTX", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod4, "name": "NVIDIA RTX A6000", "price": uint(349999), "model_no": uint(6000),
			"stock": uint(10), "category_id": catWorkstation, "brand": "NVIDIA",
			"description": "Professional workstation GPU with 48GB GDDR6 for AI, rendering, and simulation.",
			"specifications": bson.M{"vram": "48GB GDDR6", "cuda_cores": "10752", "boost_clock": "1800 MHz", "tdp": "300W", "ecc_memory": "Yes"},
			"image1": "https://placehold.co/400x300/1a1a2e/3498db?text=RTX+A6000", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod5, "name": "NVIDIA RTX 4000 Ada", "price": uint(124999), "model_no": uint(4000),
			"stock": uint(20), "category_id": catWorkstation, "brand": "NVIDIA", "discount_id": disc20,
			"description": "Compact professional GPU with 20GB GDDR6 for CAD and 3D workflows.",
			"specifications": bson.M{"vram": "20GB GDDR6", "cuda_cores": "6144", "boost_clock": "2175 MHz", "tdp": "130W", "form_factor": "Single Slot"},
			"image1": "https://placehold.co/400x300/1a1a2e/3498db?text=RTX+4000+Ada", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod6, "name": "NVIDIA CMP 170HX", "price": uint(249999), "model_no": uint(170),
			"stock": uint(5), "category_id": catMining, "brand": "NVIDIA",
			"description": "Dedicated cryptocurrency mining processor, no video output.",
			"specifications": bson.M{"hash_rate": "164 MH/s", "memory": "8GB HBM2e", "tdp": "250W", "video_output": "None"},
			"image1": "https://placehold.co/400x300/1a1a2e/f39c12?text=CMP+170HX", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod7, "name": "NVIDIA GeForce RTX 4060 Ti", "price": uint(41999), "model_no": uint(4060),
			"stock": uint(60), "category_id": catBudget, "brand": "NVIDIA", "discount_id": disc10,
			"description": "Great 1080p and 1440p gaming with 8GB GDDR6 and DLSS 3 support.",
			"specifications": bson.M{"vram": "8GB GDDR6", "cuda_cores": "4352", "boost_clock": "2535 MHz", "tdp": "160W", "interface": "PCIe 4.0 x8"},
			"image1": "https://placehold.co/400x300/1a1a2e/10b981?text=RTX+4060+Ti", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod8, "name": "AMD Radeon RX 7600", "price": uint(27999), "model_no": uint(7600),
			"stock": uint(80), "category_id": catBudget, "brand": "AMD",
			"description": "Affordable 1080p gaming GPU with 8GB GDDR6.",
			"specifications": bson.M{"vram": "8GB GDDR6", "stream_processors": "2048", "boost_clock": "2655 MHz", "tdp": "165W", "interface": "PCIe 4.0 x8"},
			"image1": "https://placehold.co/400x300/1a1a2e/e74c3c?text=RX+7600", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod9, "name": "Intel Arc A770", "price": uint(32999), "model_no": uint(770),
			"stock": uint(45), "category_id": catBudget, "brand": "Intel",
			"description": "Intel's flagship discrete GPU with 16GB GDDR6 and ray tracing.",
			"specifications": bson.M{"vram": "16GB GDDR6", "xe_cores": "32", "boost_clock": "2100 MHz", "tdp": "225W", "interface": "PCIe 4.0 x16"},
			"image1": "https://placehold.co/400x300/1a1a2e/2980b9?text=Arc+A770", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": prod10, "name": "NVIDIA GeForce RTX 4070 Ti Super", "price": uint(79999), "model_no": uint(4070),
			"stock": uint(30), "category_id": catGaming, "brand": "NVIDIA", "discount_id": disc15,
			"description": "Sweet spot for 1440p and 4K gaming with 16GB GDDR6X.",
			"specifications": bson.M{"vram": "16GB GDDR6X", "cuda_cores": "8448", "boost_clock": "2610 MHz", "tdp": "285W", "interface": "PCIe 4.0 x16"},
			"image1": "https://placehold.co/400x300/1a1a2e/10b981?text=RTX+4070+TiS", "image2": "", "image3": "", "created_at": now, "updated_at": now,
		},
	}
	insertMany(ctx, db, "products", products)

	// --- Attribute Definitions (for Gaming category) ---
	attrDefs := []interface{}{
		bson.M{"_id": primitive.NewObjectID(), "category_id": catGaming, "attribute_key": "vram", "display_name": "VRAM", "attribute_type": "string", "required": true},
		bson.M{"_id": primitive.NewObjectID(), "category_id": catGaming, "attribute_key": "boost_clock", "display_name": "Boost Clock", "attribute_type": "string", "required": true},
		bson.M{"_id": primitive.NewObjectID(), "category_id": catGaming, "attribute_key": "tdp", "display_name": "TDP", "attribute_type": "string", "required": false},
		bson.M{"_id": primitive.NewObjectID(), "category_id": catWorkstation, "attribute_key": "ecc_memory", "display_name": "ECC Memory", "attribute_type": "enum", "required": false, "options": []string{"Yes", "No"}},
	}
	insertMany(ctx, db, "product_attribute_definitions", attrDefs)

	// --- Users ---
	user1ID := primitive.NewObjectID()
	user2ID := primitive.NewObjectID()
	user3ID := primitive.NewObjectID()
	users := []interface{}{
		bson.M{
			"_id": user1ID, "first_name": "John", "last_name": "Doe", "user_name": "johndoe",
			"email": "john@example.com", "password": hashPassword("password123"),
			"phone": "9876543210", "block_status": false, "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": user2ID, "first_name": "Jane", "last_name": "Smith", "user_name": "janesmith",
			"email": "jane@example.com", "password": hashPassword("password123"),
			"phone": "9876543211", "block_status": false, "created_at": now, "updated_at": now,
		},
		bson.M{
			"_id": user3ID, "first_name": "Bob", "last_name": "Wilson", "user_name": "bobwilson",
			"email": "bob@example.com", "password": hashPassword("password123"),
			"phone": "9876543212", "block_status": true, "created_at": now, "updated_at": now,
		},
	}
	insertMany(ctx, db, "users", users)

	// --- Addresses ---
	addr1ID := primitive.NewObjectID()
	addr2ID := primitive.NewObjectID()
	addresses := []interface{}{
		bson.M{
			"_id": addr1ID, "user_id": user1ID, "name": "John Doe", "phone_number": "9876543210",
			"pincode": "682001", "house": "Flat 4B, Tower 2", "area": "MG Road",
			"landmark": "Near Metro Station", "city": "Kochi",
		},
		bson.M{
			"_id": addr2ID, "user_id": user2ID, "name": "Jane Smith", "phone_number": "9876543211",
			"pincode": "560001", "house": "12, 3rd Floor", "area": "Brigade Road",
			"landmark": "Opposite Mall", "city": "Bangalore",
		},
	}
	insertMany(ctx, db, "addresses", addresses)

	// --- Wallets ---
	wallets := []interface{}{
		bson.M{"_id": primitive.NewObjectID(), "user_id": user1ID, "balance": uint(5000)},
		bson.M{"_id": primitive.NewObjectID(), "user_id": user2ID, "balance": uint(12000)},
		bson.M{"_id": primitive.NewObjectID(), "user_id": user3ID, "balance": uint(0)},
	}
	insertMany(ctx, db, "wallets", wallets)

	// --- Wallet History ---
	walletHistory := []interface{}{
		bson.M{"_id": primitive.NewObjectID(), "user_id": user1ID, "credit": uint(5000), "debit": uint(0)},
		bson.M{"_id": primitive.NewObjectID(), "user_id": user2ID, "credit": uint(15000), "debit": uint(0)},
		bson.M{"_id": primitive.NewObjectID(), "user_id": user2ID, "credit": uint(0), "debit": uint(3000)},
	}
	insertMany(ctx, db, "wallet_history", walletHistory)

	// --- Coupons ---
	coupons := []interface{}{
		bson.M{
			"_id": primitive.NewObjectID(), "coupon_name": "WELCOME10", "coupon_code": "WELCOME10",
			"coupon_percentage": uint(10), "expiry_date": now.AddDate(0, 1, 0),
		},
		bson.M{
			"_id": primitive.NewObjectID(), "coupon_name": "GPU20", "coupon_code": "GPU20",
			"coupon_percentage": uint(20), "expiry_date": now.AddDate(0, 2, 0),
		},
		bson.M{
			"_id": primitive.NewObjectID(), "coupon_name": "EXPIRED5", "coupon_code": "EXPIRED5",
			"coupon_percentage": uint(5), "expiry_date": now.AddDate(0, 0, -1),
		},
	}
	insertMany(ctx, db, "coupons", coupons)

	// --- Carts ---
	cart1ID := primitive.NewObjectID()
	carts := []interface{}{
		bson.M{"_id": cart1ID, "user_id": user1ID},
	}
	insertMany(ctx, db, "carts", carts)

	// --- Cart Products ---
	cartProducts := []interface{}{
		bson.M{"_id": primitive.NewObjectID(), "cart_id": cart1ID, "product_id": prod7, "product_name": "NVIDIA GeForce RTX 4060 Ti", "quantity": uint(1), "product_price": uint(41999)},
		bson.M{"_id": primitive.NewObjectID(), "cart_id": cart1ID, "product_id": prod8, "product_name": "AMD Radeon RX 7600", "quantity": uint(2), "product_price": uint(27999)},
	}
	insertMany(ctx, db, "cart_products", cartProducts)

	// --- Orders ---
	order1ID := "ORD-" + primitive.NewObjectID().Hex()[:8]
	order2ID := "ORD-" + primitive.NewObjectID().Hex()[:8]
	order3ID := "ORD-" + primitive.NewObjectID().Hex()[:8]
	orders := []interface{}{
		bson.M{
			"_id": primitive.NewObjectID(), "user_id": user1ID, "address_id": addr1ID,
			"order_id": order1ID, "payment_method": "razorpay", "total_amount": uint(159999),
			"status": true, "payment_status": true,
			"expected_delivery_date": now.AddDate(0, 0, 7), "created_at": now.AddDate(0, 0, -10),
		},
		bson.M{
			"_id": primitive.NewObjectID(), "user_id": user2ID, "address_id": addr2ID,
			"order_id": order2ID, "payment_method": "wallet", "total_amount": uint(41999),
			"status": false, "payment_status": true,
			"expected_delivery_date": now.AddDate(0, 0, 5), "created_at": now.AddDate(0, 0, -3),
		},
		bson.M{
			"_id": primitive.NewObjectID(), "user_id": user1ID, "address_id": addr1ID,
			"order_id": order3ID, "payment_method": "razorpay", "total_amount": uint(79999),
			"status": false, "payment_status": false,
			"expected_delivery_date": now.AddDate(0, 0, 7), "created_at": now.AddDate(0, 0, -1),
		},
	}
	insertMany(ctx, db, "orders", orders)

	// --- Order Items ---
	orderItems := []interface{}{
		bson.M{"_id": primitive.NewObjectID(), "order_id": order1ID, "product_id": prod1, "product_name": "NVIDIA GeForce RTX 4090", "quantity": uint(1), "price": uint(159999)},
		bson.M{"_id": primitive.NewObjectID(), "order_id": order2ID, "product_id": prod7, "product_name": "NVIDIA GeForce RTX 4060 Ti", "quantity": uint(1), "price": uint(41999)},
		bson.M{"_id": primitive.NewObjectID(), "order_id": order3ID, "product_id": prod10, "product_name": "NVIDIA GeForce RTX 4070 Ti Super", "quantity": uint(1), "price": uint(79999)},
	}
	insertMany(ctx, db, "order_items", orderItems)

	// --- Payments ---
	payments := []interface{}{
		bson.M{
			"_id": primitive.NewObjectID(), "user_id": user1ID,
			"razor_payment_id": "pay_seed_001", "razor_pay_order_id": "order_seed_001",
			"signature": "sig_seed_001", "amount_paid": uint(159999),
		},
	}
	insertMany(ctx, db, "payments", payments)

	fmt.Println("[SEED] All collections seeded successfully.")
	fmt.Println("")
	fmt.Println("  Test Accounts:")
	fmt.Println("  ┌──────────────────────────────────────────────────┐")
	fmt.Println("  │ Admin:  admin@gpu.store     / admin123           │")
	fmt.Println("  │ User:   john@example.com    / password123        │")
	fmt.Println("  │ User:   jane@example.com    / password123        │")
	fmt.Println("  │ User:   bob@example.com     / password123 (blocked) │")
	fmt.Println("  └──────────────────────────────────────────────────┘")
}

func insertMany(ctx context.Context, db *mongo.Database, collection string, docs []interface{}) {
	if _, err := db.Collection(collection).InsertMany(ctx, docs); err != nil {
		log.Printf("[SEED] Failed to insert into %s: %v", collection, err)
	} else {
		fmt.Printf("[SEED] Inserted %d documents into %s\n", len(docs), collection)
	}
}
