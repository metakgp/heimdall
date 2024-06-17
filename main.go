package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/likexian/whois"
	"github.com/pquerna/otp/totp"
	"github.com/rs/cors"
)

const COOKIE_DOMAIN = ".metakgp.org"

var (
	ErrJwtSecretKeyNotFound                  = errors.New("ERROR: JWT SECRET KEY NOT FOUND")
	ErrJwtTokenExpired                       = errors.New("ERROR: JWT TOKEN EXPIRED")
	ErrJwtTokenInvalid                       = errors.New("ERROR: JWT TOKEN INVALID")
	usersMap                map[string]*User = make(map[string]*User)
)

type LoginJwtFields struct {
	Email string `json:"email"`
}

type LoginJwtClaims struct {
	LoginJwtFields
	jwt.RegisteredClaims
}

type User struct {
	Email    string `json:"email"`
	Secret   string `json:"secret"`
	LastUsed int64  `json:"last_used"`
}

type OtpResponse struct {
	Email     string `json:"email"`
	OtpStatus bool   `json:"otp_status"`
	Timestamp int    `json:"timestamp"`
}

func getJwtKey() (string, error) {
	jwtKey := os.Getenv("JWT_SECRET_KEY")

	if jwtKey == "" {
		return "", ErrJwtSecretKeyNotFound
	}

	return jwtKey, nil
}

func jwtKeyFunc(*jwt.Token) (interface{}, error) {
	key, err := getJwtKey()

	if err != nil {
		return nil, err
	}

	return []byte(key), err
}

func generateOtp(user User) (bool, error) {
	validPeriod, err := strconv.Atoi(os.Getenv("OTP_VALIDITY_PERIOD"))
	if err != nil || validPeriod < 30 { // keep 30s as minimum valid period
		fmt.Println("Invalid OTP_VALIDITY_PERIOD env set. Defaulting to 600 seconds (10 minutes)")
		validPeriod = 600
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Heimdall",
		AccountName: user.Email,
		Period:      uint(validPeriod),
	})
	if err != nil {
		fmt.Println(err)
		return false, errors.New("error generating OTP")
	}

	otp, err := totp.GenerateCode(secret.Secret(), time.Now())
	if err != nil {
		fmt.Println(err)
		return false, errors.New("error generating OTP")
	}

	otp_status, err := sendOTP(user.Email, otp)
	if err != nil || !otp_status {
		fmt.Println(err)
		return false, errors.New("error generating OTP")
	}

	currentTime := int(time.Now().Unix())
	user.Secret = secret.Secret()
	user.LastUsed = int64(currentTime)
	usersMap[user.Email] = &user
	return otp_status, nil
}

func handleCampusCheck(res http.ResponseWriter, req *http.Request) {
	clientIP := req.Header.Get("X-Forwarded-For")
	whoisResponse, err := whois.Whois(clientIP)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Define a regular expression pattern to match the netname
	pattern := `netname:\s+(.*)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find the netname using the regular expression
	match := re.FindStringSubmatch(whoisResponse)

	response := make(map[string]bool)

	if len(match) >= 2 {
		netname := match[1]
		if netname == "IITKGP-IN" {
			response["is_inside_kgp"] = true
		} else {
			response["is_inside_kgp"] = false
		}
	} else {
		fmt.Println("Netname not found in the whois response.")
		response["is_inside_kgp"] = false
	}

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
	}
	res.Write(jsonResp)
}

func handleGetOtp(res http.ResponseWriter, req *http.Request) {
	email := req.FormValue("email")
	if email == "" {
		http.Error(res, "Missing email parameter", http.StatusBadRequest)
		return
	}

	// check for KGPian email
	if !strings.HasSuffix(email, "@kgpian.iitkgp.ac.in") {
		http.Error(res, "Invalid email domain. Must be @kgpian.iitkgp.ac.in", http.StatusBadRequest)
		return
	}

	user, ok := usersMap[email]
	if ok {
		cooldown, err := strconv.Atoi(os.Getenv("RESEND_OTP_COOLDOWN"))
		if err != nil {
			fmt.Println("Invalid RESEND_OTP_COOLDOWN env set. Defaulting to 60 seconds (1 minute)")
			cooldown = 60 // keep 30s as minimum cooldown
		}
		cooldownDuration := time.Duration(cooldown) * time.Second
		if time.Now().Unix()-user.LastUsed < int64(cooldownDuration.Seconds()) {
			http.Error(res, fmt.Sprintf("You requested OTP recently. Please wait %d seconds before requesting again.", cooldown), http.StatusBadRequest)
			return
		} else {
			otp_status, err := generateOtp(*user)
			if err != nil || !otp_status {
				fmt.Println(err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}

			response := OtpResponse{
				Timestamp: int(user.LastUsed),
				Email:     email,
				OtpStatus: otp_status,
			}

			respJson, err := json.Marshal(response)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Return JSON response with OTP
			res.Header().Set("Content-Type", "application/json")
			res.Write(respJson)
			return
		}
	}

	var newUser User
	newUser.Email = email
	otp_status, err := generateOtp(newUser)
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	response := OtpResponse{
		Timestamp: int(newUser.LastUsed),
		Email:     email,
		OtpStatus: otp_status,
	}

	respJson, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return JSON response with OTP
	res.Header().Set("Content-Type", "application/json")
	res.Write(respJson)
}

func handleVerifyOtp(res http.ResponseWriter, req *http.Request) {
	email := req.FormValue("email")
	if email == "" {
		http.Error(res, "Missing email parameter", http.StatusBadRequest)
		return
	}

	otp := req.FormValue("otp")
	if otp == "" {
		http.Error(res, "Missing otp parameter", http.StatusBadRequest)
		return
	}

	user, ok := usersMap[email]
	if !ok {
		http.Error(res, "Please Request OTP first", http.StatusBadRequest)
		return
	}

	valid := totp.Validate(otp, user.Secret)
	if !valid {
		http.Error(res, "Invalid OTP", http.StatusBadRequest)
		return
	}

	signingKey, err := getJwtKey()
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	expiryDays, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_DAYS"))
	if err != nil || expiryDays < 1 { // keep 1 day as minimum valid period
		fmt.Println("Invalid JWT_EXPIRY_DAYS env set. Defaulting to 90 days (3 months)")
		expiryDays = 90 // Default to 90 days (3 months)
	}

	issueTime := time.Now()
	expiryTime := issueTime.AddDate(0, 0, expiryDays)

	claims := &LoginJwtClaims{
		LoginJwtFields: LoginJwtFields{Email: user.Email},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issueTime),
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		fmt.Println("Could not parse sign token")
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "heimdall",
		Value:    tokenString,
		Expires:  expiryTime,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Domain:   COOKIE_DOMAIN,
	}

	http.SetCookie(res, &cookie)
}

func handleValidateJwt(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("jwt")
	if err != nil {
		http.Error(res, "No JWT session token found.", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value

	var loginClaims = LoginJwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &loginClaims, jwtKeyFunc)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(res, "Invalid token signature", http.StatusBadRequest)
			return
		}
		if err.Error() == fmt.Sprintf("%s: %s", jwt.ErrTokenInvalidClaims.Error(), jwt.ErrTokenExpired.Error()) {
			http.Error(res, ErrJwtTokenExpired.Error(), http.StatusUnauthorized)
			return
		}
	}

	if !token.Valid {
		http.Error(res, ErrJwtTokenInvalid.Error(), http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*LoginJwtClaims)
	if !ok {
		http.Error(res, ErrJwtTokenInvalid.Error(), http.StatusBadRequest)
		return
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		fmt.Println(res, "Error marshalling claims to JSON: %v", err)
		http.Error(res, ErrJwtTokenInvalid.Error(), http.StatusUnauthorized)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(claimsJSON)
}
func main() {
	godotenv.Load()

	initMailer()

	generalCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://heimdall.metakgp.org"},
		AllowCredentials: true,
	})

	specialCors := cors.AllowAll()

	mux := http.NewServeMux()
	mux.Handle("/campus-check", specialCors.Handler(http.HandlerFunc(handleCampusCheck)))
	mux.Handle("/get-otp", generalCors.Handler(http.HandlerFunc(handleGetOtp)))
	mux.Handle("/verify-otp", generalCors.Handler(http.HandlerFunc(handleVerifyOtp)))
	mux.Handle("/validate-jwt", generalCors.Handler(http.HandlerFunc(handleValidateJwt)))

	handler := cors.AllowAll().Handler(mux)

	fmt.Println("Heimdall Server running on port : 3333")
	err := http.ListenAndServe(":3333", handler)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
