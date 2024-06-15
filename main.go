package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
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

var ErrJwtSecretKeyNotFound = errors.New("ERROR: JWT SECRET KEY NOT FOUND")
var ErrJwtTokenExpired = errors.New("ERROR: JWT TOKEN EXPIRED")
var ErrJwtTokenInvalid = errors.New("ERROR: JWT TOKEN INVALID")

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

var usersMap map[string]*User = make(map[string]*User)

type OtpResponse struct {
	Email     string `json:"email"`
	OtpStatus bool   `json:"otp_status"`
	Timestamp int    `json:"timestamp"`
}

type VerifiedResponse struct {
	Jwt string `json:"jwt"`
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

func sendMail(emailTo, subject, body string) (bool, error) {
	fromEmail := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_APP_PASSWORD")

	if fromEmail == "" || password == "" {
		return false, fmt.Errorf("SMTP_EMAIL or SMTP_APP_PASSWORD environment variables not set")
	}

	message := "To: " + emailTo + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	auth := smtp.PlainAuth("", fromEmail, password, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, fromEmail, []string{emailTo}, []byte(message))
	if err != nil {
		return false, err
	}

	return true, nil
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

	user, ok := usersMap[email]
	if ok {
		// user already tried
		const cooldown = 120 * time.Second
		if time.Now().Unix()-user.LastUsed < int64(cooldown.Seconds()) {
			http.Error(res, "You requested OTP recently. Please wait 2min before requesting again.", http.StatusBadRequest)
			return
		}
	}

	// check for KGPian email
	if !strings.HasSuffix(email, "@kgpian.iitkgp.ac.in") {
		http.Error(res, "Invalid email domain. Must be @kgpian.iitkgp.ac.in", http.StatusBadRequest)
		return
	}

	validPeriod, err := strconv.Atoi(os.Getenv("OTP_VALIDITY_PERIOD"))
	if err != nil || validPeriod < 30 { // keep 30s as minimum valid period
		validPeriod = 600
	}

	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Heimdall",
		AccountName: email,
		Period:      uint(validPeriod),
	})
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Error generating OTP", http.StatusInternalServerError)
		return
	}

	otp, err := totp.GenerateCode(secret.Secret(), time.Now())
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Error generating OTP", http.StatusInternalServerError)
		return
	}

	var newUser User
	newUser.Email = email
	currentTime := int(time.Now().Unix())
	newUser.Secret = secret.Secret()
	newUser.LastUsed = int64(currentTime)
	usersMap[email] = &newUser

	otp_status, err := sendMail(email, "OTP for Sign In into Heimdall Portal of MetaKGP, IIT Kharagpur is "+otp, "OTP for Sign In into Heimdall Portal of MetaKGP, IIT Kharagpur is "+otp)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Error generating OTP", http.StatusInternalServerError)
		return
	}

	response := OtpResponse{
		Timestamp: currentTime,
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
		http.Error(res, "Please Request OTP before verifying", http.StatusBadRequest)
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
	}

	issue_time := time.Now()
	claims := &LoginJwtClaims{
		LoginJwtFields: LoginJwtFields{Email: user.Email},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(issue_time),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		fmt.Println("Could not parse sign token")
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := VerifiedResponse{Jwt: tokenString}
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

func handleValidateJwt(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(res, "Authorization header missing", http.StatusBadRequest)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(res, "Invalid Authorization header", http.StatusBadRequest)
		return
	}

	tokenString := parts[1]

	if tokenString == "" {
		http.Error(res, "No JWT session token found.", http.StatusBadRequest)
		return
	}

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
		http.Error(res, "Failed to parse claims", http.StatusBadRequest)
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
	mux := http.NewServeMux()
	mux.HandleFunc("/campus_check", handleCampusCheck)
	mux.HandleFunc("/get-otp", handleGetOtp)
	mux.HandleFunc("/verify-otp", handleVerifyOtp)
	mux.HandleFunc("/validate-jwt", handleValidateJwt)

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
