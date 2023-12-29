package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"example.com/m/v2/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserPost struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type OtpPost struct {
	PhoneNumber string `json:"phone_number"`
}

type OtpPostVerify struct {
	PhoneNumber string `json:"phone_number"`
	Otp         string `json:"otp"`
}

func postUser(c *gin.Context) {
	body := UserPost{}
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "User is not defined")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.AbortWithStatusJSON(400, "Bad Input")
		return
	}

	conn, err := pgx.Connect(c, getDbConnectionString())
	if err != nil {
		c.AbortWithStatusJSON(400, "Database error.")
	}
	defer conn.Close(c)

	queries := user.New(conn)

	insertedUser, err1 := queries.CreateUser(c, user.CreateUserParams{
		Name:        body.Name,
		PhoneNumber: body.PhoneNumber,
	})
	if err1 != nil {
		c.AbortWithStatusJSON(400, "Phone number already exists or input not valid"+err1.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, insertedUser)
}

func generateOtp(c *gin.Context) {
	body := OtpPost{}
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "Phone number is required")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		c.AbortWithStatusJSON(400, "Bad Input")
		return
	}

	conn, err := pgx.Connect(c, getDbConnectionString())
	if err != nil {
		c.AbortWithStatusJSON(400, "Database connection error.")
	}
	defer conn.Close(c)

	queries := user.New(conn)
	phoneNumber, err1 := queries.GetPhonenumber(c, body.PhoneNumber)
	if err1 != nil {
		c.AbortWithStatusJSON(404, "Phone number not found")
		return
	}

	otp := fmt.Sprintf("%04d", rand.Intn(9999))

	otpExpiry := pgtype.Timestamp{
		Time:  time.Now().Add(time.Minute * 1),
		Valid: true,
	}

	err2 := queries.UpdateOtp(c, user.UpdateOtpParams{
		PhoneNumber:       phoneNumber,
		Otp:               pgtype.Text{String: otp, Valid: true},
		OtpExpirationTime: otpExpiry,
	})

	if err2 != nil {
		fmt.Println("error occured while saving", err2.Error())
	}
}

func verifyOtp(c *gin.Context) {
	body := OtpPostVerify{}
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "Phone number & otp are required")
		return
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		c.AbortWithStatusJSON(400, "Bad Input")
		return
	}

	conn, err := pgx.Connect(c, getDbConnectionString())
	if err != nil {
		c.AbortWithStatusJSON(400, "Database connection error.")
	}
	defer conn.Close(c)

	queries := user.New(conn)
	record, err2 := queries.VerifyOtp(c, user.VerifyOtpParams{
		PhoneNumber: body.PhoneNumber,
		Otp:         pgtype.Text{String: body.Otp, Valid: true},
	})

	if err2 != nil {
		c.AbortWithStatusJSON(404, "Phone number or otp not found")
		return
	}

	if (record.OtpExpirationTime).Time.After(time.Now().UTC()) {
		c.IndentedJSON(http.StatusOK, "otp is valid")
	} else if (record.OtpExpirationTime).Time.Before(time.Now()) {
		c.AbortWithStatusJSON(404, "Otp has expired ")
		return
	} else {
		c.IndentedJSON(http.StatusOK, "otp is valid")
	}
}

func getDbConnectionString() string {
	return "user=postgres dbname=testdb host=localhost password=mysecretpassword port=5432 sslmode=disable timezone=UTC"
}

func main() {
	router := gin.Default()

	router.POST("/api/users/generateotp", generateOtp)
	router.POST("/api/users/verifyotp", verifyOtp)
	router.POST("/api/users", postUser)

	router.Run("localhost:8080")
}
