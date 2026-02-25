package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	MongoURI        string
	MongoDatabase   string
	JWTSecret       string
	RazorpayKey     string
	RazorpaySecret  string
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioServiceSID string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:            getEnv("PORT", ":8080"),
		MongoURI:        getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase:   getEnv("MONGODB_DATABASE", "ecommerce_store"),
		JWTSecret:       getEnv("JWT_SECRET", "supersecretkey"),
		RazorpayKey:     getEnv("RAZORPAY_KEY", ""),
		RazorpaySecret:  getEnv("RAZORPAY_SECRET", ""),
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioServiceSID: getEnv("TWILIO_SERVICE_SID", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
