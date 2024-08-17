package main

import (
	"database/sql"

	"mime/multipart"

	"io"

	"os"

	"path/filepath"

	"strings"

	"github.com/go-sql-driver/mysql"

	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"

	"fmt"

	"math/rand"

	"github.com/dgrijalva/jwt-go"

	"time"

	"github.com/gin-contrib/cors"
)

/* <Structs init> */

type User struct {
	UserName     string `json:"UserName"`
	UserPassword string `json:"UserPassword"`
	UserEmail    string `json:"UserEmail"`
}

type UserName struct {
	UserName string `json:"UserName"`
}

type UserPassword struct {
	UserPassword string `json:"UserPassword"`
}

type UserEmail struct {
	UserEmail string `json:"UserEmail"`
}

type UserLogin struct {
	UserName     string `json:"UserName"`
	UserPassword string `json:"UserPassword"`
}

type UserDispatch struct {
	DispatchID int    `json:"DispatchID"`
	UserName   string `json:"UserName"`
}

type Dispatch struct {
	File        *multipart.FileHeader `form:"File"`
	Title       string                `form:"Title"`
	Description string                `form:"Description"`
	Text        string                `form:"Text"`
}

type DispatchAnswer struct {
	File          []byte `form:"File"`
	Title         string `form:"Title"`
	Description   string `form:"Description"`
	Text          string `form:"Text"`
	FileExtension string `form:"FileExtension"`
}

type Claims struct {
	UserName string `json:"UserName"`
	jwt.StandardClaims
}

type EmailClaims struct {
	UserEmail string `json:"UserEmail"`
	jwt.StandardClaims
}

/* </Structs init> */

/* <Interfaces init> */

type UserMethods interface {
	addUserToDB() bool
	getUserDispatches() ([]int, bool)
	setUserDispatch() (int, bool)
}

/* </Interfaces init> */

/* <UserMethods> */

func (u User) addUserToDB() bool { // Adding User to DB by his credentials
	_, err := db.Exec("INSERT INTO Users (UserName, UserPassword, UserEmail) VALUES (?, ?, ?)", u.UserName, u.UserPassword, u.UserEmail) // Adding to DB

	if err != nil { // If error appeared
		fmt.Println("Error in |func addUserToDB|:", err.Error())
		return true
	}

	return false // If no errors
}

func (u User) getUserDispatches() ([]int, bool) { // Getting all Dispatch ids made by User
	var usersDispatches []int

	rows, err := db.Query("SELECT * FROM Users_Dispatches WHERE UserName = ?", u.UserName) // Getting them from DB

	if err != nil { // If error appeared
		fmt.Println("Error in |func getUserDispatches|:", err.Error())
		return []int{}, true
	}

	defer rows.Close() // Don't forget to close DB connection

	for rows.Next() { // Do while we have rows to read

		var dispatch UserDispatch // Create UserDispatch that will read row info

		if err := rows.Scan(&dispatch.DispatchID, &dispatch.UserName); err != nil { // Read row info and write into created UserDispatch by rows.Scan(),rows.Scan() returns error if it has appeared
			fmt.Println("Error in |func getUserDispatches|:", err.Error())
			return []int{}, true
		}
		usersDispatches = append(usersDispatches, dispatch.DispatchID) // If no errors
	}

	if err := rows.Err(); err != nil { // If error appeared
		fmt.Println("Error in |func getUserDispatches|:", err.Error())
		return []int{}, true
	}

	// If no errors:
	fmt.Printf("Successfully got UserDispatches for User: %s\n", u.UserName)
	return usersDispatches, false
}

func (u User) setUserDispatch() (int, bool) { // Insert into DB info about author of new dispatch
	var id, wasErr = generateDispatchId() // Creating dispatch id

	if !wasErr {
		_, err := db.Exec("INSERT INTO Users_Dispatches (DispatchID, UserName) VALUES (?, ?)", id, u.UserName) // Inserting into DB

		if err != nil { // If error appeared
			fmt.Println("Error in |func addUserToDB|:", err.Error())
			return 0, true
		}

		// If no errors:
		fmt.Printf("Successfully added Dispatch from User: %s\n", u.UserName)
		return id, false

	} else { // If error appeared
		fmt.Println("Can't complete |func setUserDispatch| because of ERROR")
		return 0, true
	}
}

func generateDispatchId() (id int, wasErrors bool) { // Generate unique ID
	for { // Do, while we don't found unique ID
		var dispatch UserDispatch // Create UserDispatch

		id = rand.Intn(100000-10000) + 10000 // Generate random ID in range 10000-99999

		row := db.QueryRow("SELECT * FROM Users_Dispatches WHERE DispatchID = ?", id) // Check if that ID already exist in DB

		if err := row.Scan(&dispatch.DispatchID, &dispatch.UserName); err != nil { // If ID don't exist error will appear, else - ID exist
			if err == sql.ErrNoRows { // Check if error caused by that ID already exist
				return id, false
			} else { // If error appeared because of something else
				fmt.Println("Error in |func generateDispatchId|:", err.Error())
				return id, true
			}
		}

	}
}

/* </UserMethods> */

/* <JWT Tokens functions> */

func GenerateJWTToken(username string) (string, bool) { // Generate new JWT token for auth
	expirationTime := time.Now().Add(1440 * time.Minute) // Time when token won't be valid

	claims := &Claims{ // Create Claims that will be need in JWT creation
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Create JWT token

	tokenString, err := token.SignedString(jwtKey) // Sign JWT token with our secret key

	if err != nil { // If error appeared
		fmt.Println("Error in |func GenerateJWTToken|: ", err.Error())
		return "", true
	}

	return tokenString, false // If no errors
}

func DecodeJWTToken(tokenToDecode string) (string, bool) { // Decode JWT token for auth
	claims := &Claims{} // Create Claims

	tokenDecoded, err := jwt.ParseWithClaims(tokenToDecode, claims, func(token *jwt.Token) (interface{}, error) { return jwtKey, nil }) // Decode JWT token for auth

	if err != nil { // If error appeared

		if strings.Split(err.Error(), " ")[2] == "expired" { // If JWT token expired
			fmt.Println("Error in |func DecodeToken| because JWT Token is expired and now invalid: ", err)
			return "InValid", true
		}

		fmt.Println("Error in |func DecodeToken|: ", err)
		return "", true
	}

	if !tokenDecoded.Valid { // If JWT token isn't valid
		fmt.Println("Error in |func DecodeToken| because JWT Token is InValid: ", err)
		return "InValid", true
	}

	return claims.UserName, false // If no errors
}

func GenerateEmailJWTToken(email string) (string, bool) { // Generate new JWT token for mailAuth
	expirationTime := time.Now().Add(1440 * time.Minute) // Time when token won't be valid

	claims := &EmailClaims{ // Create EmailClaims
		UserEmail: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Create JWT token

	tokenString, err := token.SignedString(jwtEmailKey) // Sign JWT token with our another secret key

	if err != nil { // If error appeared
		fmt.Println("Error in |func GenerateEmailJWTToken|: ", err.Error())
		return "", true
	}

	return tokenString, false // If no errors
}

func DecodeEmailJWTToken(tokenToDecode string) (string, bool) { // Decode JWT token for mailAuth
	claims := &EmailClaims{} // Create EmailClaims

	tokenDecoded, err := jwt.ParseWithClaims(tokenToDecode, claims, func(token *jwt.Token) (interface{}, error) { return jwtEmailKey, nil }) // Decode JWT token for mailAuth

	if err != nil { // If error appeared

		if strings.Split(err.Error(), " ")[2] == "expired" { // If JWT token expired
			fmt.Println("Error in |func DecodeToken| because JWT Token is expired and now invalid: ", err)
			return "InValid", true
		}

		fmt.Println("Error in |func DecodeEmailJWTToken|: ", err)
		return "", true
	}

	if !tokenDecoded.Valid { // If JWT token isn't valid
		fmt.Println("Error in |func DecodeEmailJWTToken| because JWT Token is InValid: ", err)
		return "InValid", true
	}

	return claims.UserEmail, false // If no errors
}

/* </JWT Tokens functions> */

/* <User functions> */

func getUserFromDB(username string) (User, bool) { // Getting User credentials from DB
	var u User // Create User

	row := db.QueryRow("SELECT * FROM Users WHERE UserName = ?", username) // Find row with User info in DB

	if err := row.Scan(&u.UserName, &u.UserPassword, &u.UserEmail); err != nil { // Read row into User
		if err == sql.ErrNoRows { // If appeared error caused by ErrNoRows
			fmt.Println("Error in |func getUserFromDB|:", err.Error())
			u.UserName = "No such User"
			return u, true
		}
		// If other error appeared
		fmt.Println("Error in |func getUserFromDB|:", err.Error())
		return u, true
	}

	return u, false // If no errors
}

func signUp(c *gin.Context) { // Endpoint to sign up User
	var user User                 // Create User
	var userInterface UserMethods // Create UserMethods

	if err := c.BindJSON(&user); err != nil { // Read user data from request body, c.BindJSON() returns error
		// If error appeared:
		fmt.Println("Error in |func signUp|: ", err)
		c.JSON(http.StatusBadRequest, user)
		return
	}

	userInterface = user                    // Connect interface with struct
	wasError := userInterface.addUserToDB() // Add User to DB by func, addUserToDB() returns error

	if wasError { // If error appered
		row := db.QueryRow("SELECT * FROM Users WHERE UserName = ?", user.UserName) // Find User in DD

		if err := row.Scan(&user.UserName, &user.UserPassword, &user.UserEmail); err != nil { // Read User info from DB

			// If error appeared while reading:

			if err == sql.ErrNoRows { // If error caused by ErrNoRows

				row_email := db.QueryRow("SELECT * FROM Users WHERE UserEmail = ?", user.UserEmail) // Find User with that email in DB

				if err := row_email.Scan(&user.UserName, &user.UserPassword, &user.UserEmail); err != nil { // Read User info from DB

					// If error appeared while reading:
					fmt.Println("Error in |func signUp| while checking UserEmail: ", err.Error())
					c.JSON(http.StatusServiceUnavailable, err.Error())
					return

				}

				// If error appeared because of User with that email already exist:
				fmt.Println("Error in |func signUp| because User already exist with same UserEmail")
				c.JSON(http.StatusConflict, "User already exist with same UserEmail")
				return

			}

			// If error appeared because of something else
			fmt.Println("Error in |func signUp| while checking UserName: ", err.Error())
			c.JSON(http.StatusServiceUnavailable, err.Error())
			return
		}

		// If error appeared because User already exist
		fmt.Println("Error in |func signUp| because User already exist with same UserName")
		c.JSON(http.StatusConflict, "User already exist with same UserName")
		return
	}

	c.JSON(http.StatusCreated, user) // If no errors
}

func logIn(c *gin.Context) { // Endpoint to log in User
	var user_login UserLogin // Create UserLogin

	if err := c.BindJSON(&user_login); err != nil { // Read user data from request body
		// If error appeared:
		fmt.Println("Error in |func logIn|: ", err)
		c.JSON(http.StatusBadRequest, user_login)
		return
	}

	user_in_db, wasError := getUserFromDB(user_login.UserName) // Get User with that username from DB

	if wasError { // If error appeared

		if user_in_db.UserName == "No such User" { // If error appeared because User don't exists
			fmt.Println("Error in |func logIn| while fetching User from DB: No such User")
			c.JSON(http.StatusForbidden, "User doesn't exist")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func logIn| while fetching User from DB: DB error")
		c.JSON(http.StatusServiceUnavailable, "DB error")
		return
	}

	if !(user_login.UserPassword == user_in_db.UserPassword) { // Check: is password right
		// If passwords aren't the same
		c.JSON(http.StatusForbidden, "Passwords don't compare")
		return
	}

	jwtToken, wasError := GenerateJWTToken(user_login.UserName) // Generate JWT token for auth

	if wasError { // If error appeared
		fmt.Println("Error in |func logIn|: something went wrong while generating JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token generating error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken}) // If no errors
}

func auth(c *gin.Context) { // Endpoint to check User's token
	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared
		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func authenticate|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func authenticate|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	user, wasError := getUserFromDB(userName) // Get User with that username from DB

	if wasError { // If error appeared

		if user.UserName == "No such User" { // If error appeared because User don't exists
			fmt.Println("Error in |func authenticate| while fetching User from DB: No such User")
			c.JSON(http.StatusUnauthorized, "User doesn't exist")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func authenticate| while fetching User from DB: DB error")
		c.JSON(http.StatusServiceUnavailable, "DB error")
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"UserName": userName, "UserEmail": user.UserEmail}) // If no errors
}

func changeUsername(c *gin.Context) { // Endpoint to change username
	var newUserName UserName // Create UserName

	if err := c.BindJSON(&newUserName); err != nil { // Read username from request body
		// If error appeared:
		fmt.Println("Error in |func changeUsername|: ", err)
		c.JSON(http.StatusBadRequest, newUserName)
		return
	}

	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared

		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func changeUsername|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func changeUsername|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	_, err := db.Query("UPDATE Users SET UserName = ? WHERE UserName = ?", newUserName.UserName, userName) // Change username in DB

	if err != nil { // If error appeared

		fmt.Println("Error in |func changeUsername| while changing UserName: ", err.Error())
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return

	}

	c.JSON(http.StatusAccepted, "UserName has changed") // If no errors
}

func changeEmail(c *gin.Context) { // Endpoint to change email
	var newUserEmail UserEmail // Create UserEmail

	if err := c.BindJSON(&newUserEmail); err != nil { // Read user's email from request body
		fmt.Println("Error in |func changeEmail|: ", err)
		c.JSON(http.StatusBadRequest, newUserEmail)
		return
	}

	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared

		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func changeEmail|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func changeEmail|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	_, err := db.Query("UPDATE Users SET UserEmail = ? WHERE UserName = ?", newUserEmail.UserEmail, userName) // Change user's email in DB

	if err != nil { // If error appeared

		fmt.Println("Error in |func changeEmail| while changing UserEmail: ", err.Error())
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return

	}

	c.JSON(http.StatusAccepted, "UserEmail has changed") // If no errors
}

func deleteUser(c *gin.Context) { // Endpoint to delete user
	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared

		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func deleteUser|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func deleteUser|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	_, err := db.Query("DELETE FROM Users WHERE UserName=?", userName) // Delete user from DB

	if err != nil { // If error appeared

		fmt.Println("Error in |func deleteUser| while changing UserEmail: ", err.Error())
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return

	}

	c.JSON(http.StatusAccepted, "User deleted") // If no errors
}

func getUserDispatches(c *gin.Context) { // Endpoint to get user's dispatches ids
	var user User                 // Create User
	var userInterface UserMethods // Create UserMethods

	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared

		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func getUserDispatches|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func getUserDispatches|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	user.UserName = userName // Set username in struct
	userInterface = user     // Connect interface with struct

	ids, wasError := userInterface.getUserDispatches() // Get user Dispatches from DB

	if wasError { // If error appeared
		fmt.Println("Error in |func getUserDispatches|: something went wrong while getting User Dispatches from DB")
		c.JSON(http.StatusServiceUnavailable, "Error while getting User Dispatches from DB")
		return
	}

	c.JSON(http.StatusAccepted, ids) // If no errors
}

func getAnyDispatches() ([]int, bool) { // Get all dispatches ids
	var usersDispatches []int // Create int slice

	rows, err := db.Query("SELECT * FROM Users_Dispatches") // Get dispatches from DB

	if err != nil { // If error appeared
		fmt.Println("Error in |func getAnyDispatches|:", err.Error())
		return []int{}, true
	}

	defer rows.Close() // Don't forget to close DB connection

	for rows.Next() { // For each row
		var dispatch UserDispatch                                                   // Create UserDispatch
		if err := rows.Scan(&dispatch.DispatchID, &dispatch.UserName); err != nil { // Read row info into

			// If error appeared:
			fmt.Println("Error in |func getAnyDispatches|:", err.Error())
			return []int{}, true
		}

		usersDispatches = append(usersDispatches, dispatch.DispatchID) // If no errors
	}

	if err := rows.Err(); err != nil { // If error appeared while getting rows from DB
		fmt.Println("Error in |func getAnyDispatches|:", err.Error())
		return []int{}, true
	}

	return usersDispatches, false // If no errors
}

func getDispatches(c *gin.Context) { // Endpoint to get list of all ids of dispatches

	ids, wasError := getAnyDispatches() // Get ids list

	if wasError { // If error appeared
		fmt.Println("Error in |func getDispatches|: something went wrong while getting User Dispatches from DB")
		c.JSON(http.StatusServiceUnavailable, "Error while getting User Dispatches from DB")
		return
	}

	c.JSON(http.StatusAccepted, ids) // If no errors
}

/* </User functions> */

/* <Dispatch functions> */

func createDispatch(c *gin.Context) { // Endpoint to create Dispatch
	var dispatch Dispatch         // Create Dispatch
	var user User                 // Create User
	var userInterface UserMethods // Create UserMethods

	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared

		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func createDispatch|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func createDispatch|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	if err := c.ShouldBind(&dispatch); err != nil { // Read dispatch info from request body
		// If error appeared:
		fmt.Println("Error in |func createDispatch| while binding dispatch: ", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user.UserName = userName                        // Set username in struct
	userInterface = user                            // Connect interface with struct
	id, wasError := userInterface.setUserDispatch() // Add dispatch to DB

	os.MkdirAll(fmt.Sprintf("./%d", id), os.ModePerm) // Create dir where dispatch data will be storage

	title, err := os.Create(fmt.Sprintf("./%d/title.txt", id)) // Create title.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func createDispatch| while saving title:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while saving title")
		return
	}

	description, err := os.Create(fmt.Sprintf("./%d/description.txt", id)) // Create description.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func createDispatch| while saving description:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while saving description")
		return
	}

	text, err := os.Create(fmt.Sprintf("./%d/text.txt", id)) // Create text.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func createDispatch| while saving text:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while saving text")
		return
	}

	// Don't forget to close all files
	defer title.Close()
	defer description.Close()
	defer text.Close()

	// Write data to it's files
	title.WriteString(dispatch.Title)
	description.WriteString(dispatch.Description)
	text.WriteString(dispatch.Text)

	imageFormat := fmt.Sprintf(".%s", strings.Split(dispatch.File.Header.Get("Content-Type"), "/")[len(strings.Split(dispatch.File.Header.Get("Content-Type"), "/"))-1]) // Get image format from File header sended with image in body
	c.SaveUploadedFile(dispatch.File, fmt.Sprintf("./%d/image%s", id, imageFormat))                                                                                      // Save image from request to dispatch folder

	if wasError { // If error appeared
		fmt.Println("Error in |func createDispatch| while creating dispatch")
		c.JSON(http.StatusServiceUnavailable, "Error while creating dispatch")
		return
	}

	c.JSON(http.StatusCreated, id) // If no errors
}

func saveDispatch(c *gin.Context) { // Endpoint to save dispatch (if dispatch already exist)
	var dispatch Dispatch // Create Dispatch

	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header
	id := c.GetHeader("ID")                       // Get dispatch Id from ID header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared

		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func saveDispatch|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func saveDispatch|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	if err := c.ShouldBind(&dispatch); err != nil { // Read dispatch data from request
		fmt.Println("Error in |func saveDispatch| while binding dispatch: ", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	title, err := os.Create(fmt.Sprintf("./%s/title.txt", id)) // Rewrite title.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func saveDispatch| while saving title:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while saving title")
		return
	}

	description, err := os.Create(fmt.Sprintf("./%s/description.txt", id)) // Rewrite description.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func saveDispatch| while saving description:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while saving description")
		return
	}

	text, err := os.Create(fmt.Sprintf("./%s/text.txt", id)) // Rewrite text.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func saveDispatch| while saving text:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while saving text")
		return
	}

	// Don't forget to close all files
	defer title.Close()
	defer description.Close()
	defer text.Close()

	// Write data to it's files
	title.WriteString(dispatch.Title)
	description.WriteString(dispatch.Description)
	text.WriteString(dispatch.Text)

	imageFormat := fmt.Sprintf(".%s", strings.Split(dispatch.File.Header.Get("Content-Type"), "/")[len(strings.Split(dispatch.File.Header.Get("Content-Type"), "/"))-1]) // Get image format from File header sended with image in body
	c.SaveUploadedFile(dispatch.File, fmt.Sprintf("./%s/image%s", id, imageFormat))                                                                                      // Save image from request to dispatch folder

	if wasError { // If error appeared

		fmt.Println("Error in |func saveDispatch| while creating dispatch")
		c.JSON(http.StatusServiceUnavailable, "Error while creating dispatch")
		return

	}

	c.JSON(http.StatusAccepted, id) // If no errors
}

func getDispatch(c *gin.Context) { // Endpoint to get dispatch data from folder
	var dispatch DispatchAnswer // Create DispatchAnswer

	id := c.GetHeader("ID") // Get dispatch Id from ID header

	title_raw, err := os.Open(fmt.Sprintf("./%s/title.txt", id)) // Open title.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while opening title.txt:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while opening title.txt")
		return
	}

	title, err := io.ReadAll(title_raw) // Read title.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while reading title.txt:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while reading title.txt")
		return
	}

	description_raw, err := os.Open(fmt.Sprintf("./%s/description.txt", id)) // Open description.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while opening description.txt:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while opening description.txt")
		return
	}

	description, err := io.ReadAll(description_raw) // Read description.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while reading description.txt:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while reading description.txt")
		return
	}

	text_raw, err := os.Open(fmt.Sprintf("./%s/text.txt", id)) // Open text.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while opening text.txt:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while opening text.txt")
		return
	}

	text, err := io.ReadAll(text_raw) // Read text.txt

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while reading text.txt:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while reading text.txt")
		return
	}

	paths, err := filepath.Glob(fmt.Sprintf("./%s/image", id) + ".*") // Get path to image

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while searching image:")
		c.JSON(http.StatusServiceUnavailable, "Error while searching image")
		return
	}

	if len(paths) == 0 { // If no paths found
		fmt.Println("Error in |func getDispatch| while opening image: no image files founded")
		c.JSON(http.StatusServiceUnavailable, "Error while opening image")
		return
	}

	if len(paths) > 1 { // If there are more than 1 image in dir
		fmt.Println("Error in |func getDispatch| while opening image: too many image files founded")
		c.JSON(http.StatusServiceUnavailable, "Error while opening image")
		return
	}

	image_raw, err := os.Open(paths[0]) // Open image

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while opening image:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while opening image")
		return
	}

	image, err := io.ReadAll(image_raw) // Read image

	if err != nil { // If error appeared
		fmt.Println("Error in |func getDispatch| while reading image:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while reading image")
		return
	}

	fileExtension := strings.Split(paths[0], ".")[len(strings.Split(paths[0], "."))-1] // Get image extension
	imageExtension := fmt.Sprintf("image/%s", fileExtension)                           // Add image extension to make File Header

	// Don't forget to close all files
	defer title_raw.Close()
	defer description_raw.Close()
	defer text_raw.Close()
	defer image_raw.Close()

	// Create response
	dispatch.Title = string(title)
	dispatch.Description = string(description)
	dispatch.Text = string(text)
	dispatch.File = image
	dispatch.FileExtension = imageExtension

	c.JSON(http.StatusAccepted, dispatch) // If no errors
}

func deleteDispatch(c *gin.Context) { // Endpoint to delete existing dispatch
	var dispatch UserDispatch // Create UserDispatch

	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header
	id := c.GetHeader("ID")                       // Get dispatch Id from ID header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared

		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func deleteDispatch|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func deleteDispatch|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	userDispatch := db.QueryRow("SELECT * FROM Users_Dispatches WHERE DispatchID = ?", id) // Find dispatch in DB

	if err := userDispatch.Scan(&dispatch.DispatchID, &dispatch.UserName); err != nil { // Read dispatch info
		// If error appeared:
		fmt.Println("Error in |func deleteDispatch| while searching for UserDispatch in DB:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while searching for UserDispatch")
		return
	}

	if dispatch.UserName == userName { // If user has rights to delete dispatch

		_, err := db.Query("DELETE FROM Users_Dispatches WHERE DispatchID = ?", id) // Delete dispatch row from DB

		if err != nil { // If error appeared

			fmt.Println("Error in |func deleteDispatch| while deleting Dispatch from DB: ", err.Error())
			c.JSON(http.StatusServiceUnavailable, err.Error())
			return

		}

		err = os.RemoveAll(fmt.Sprintf("./%s", id)) // Remove dispatch dir

		if err != nil { // If error appeared

			fmt.Println("Error in |func deleteDispatch| while deleting Dispatch: ", err.Error())
			c.JSON(http.StatusServiceUnavailable, err.Error())
			return

		}

		c.JSON(http.StatusAccepted, id) // If no errors

	} else {
		c.JSON(http.StatusForbidden, id) // If User don't have rights to remove this dispatch
	}
}

func canEdit(c *gin.Context) { // Endpoint to check: is requested user - author of dispatch (Can user edit the dispatch)
	var dispatch UserDispatch // Create UserDispatch

	requiredToken := c.GetHeader("Authorization") // Get JWT token from Authorization header
	id := c.GetHeader("ID")                       // Get dispatch Id from ID header

	if len(requiredToken) == 0 { // If no token
		c.JSON(http.StatusUnauthorized, "UnAuthorized")
		return
	}

	userName, wasError := DecodeJWTToken(requiredToken) // Decode JWT token

	if wasError { // If error appeared

		if userName == "InValid" { // If error appeared because JWT token isn't valid
			fmt.Println("Error in |func getAuthor|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func getAuthor|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	userDispatch := db.QueryRow("SELECT * FROM Users_Dispatches WHERE DispatchID = ?", id) // Get dispatch info from DB

	if err := userDispatch.Scan(&dispatch.DispatchID, &dispatch.UserName); err != nil { // Read dispatch row
		// If error appeared:
		fmt.Println("Error in |func getAuthor| while searching for UserDispatch in DB:", err)
		c.JSON(http.StatusServiceUnavailable, "Error while searching for UserDispatch")
		return
	}

	if dispatch.UserName == userName { // If requested User - Author
		c.JSON(http.StatusAccepted, id)
	} else { // If not
		c.JSON(http.StatusForbidden, id)
	}
}

/* </Dispatch functions> */

/* <Password recovery> */

func sendPasswordRecovery(c *gin.Context) { // Endpoint to send letter with link to reset password
	var email UserEmail // Create UserEmail

	auth := smtp.PlainAuth("", smtp_username, smtp_password, smtp_host_ip) // Create plain auth for user in SMTP server, from who's acccount we will send emails
	port := "25"                                                           // SMTP port

	if err := c.BindJSON(&email); err != nil { // Get email from request body
		// If error appeared:
		fmt.Println("Error in |func sendPasswordRecovery|: ", err)
		c.JSON(http.StatusBadRequest, email)
		return
	}

	jwtToken, wasError := GenerateEmailJWTToken(email.UserEmail) // Generate mail JWT token

	if wasError { // If error appeared:
		fmt.Println("Error in |func sendPasswordRecovery|: something went wrong while generating JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token generating error")
		return
	}

	to := []string{ // Make string slice from request email
		email.UserEmail,
	}

	subject := "Reset password on Read&Dispatch site"                                                     // Header of email
	body := fmt.Sprintf("Click on the link to change your password: %s?token=%s", frontend_url, jwtToken) // Body of email

	// Create message
	message := fmt.Sprintf("From: %s\r\n", smtp_sender_mail)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += fmt.Sprintf("\r\n%s\r\n", body)

	err := smtp.SendMail(smtp_host_ip+":"+port, auth, smtp_sender_mail, to, []byte(message)) // Send password-reset email

	if err != nil { // If error appeared
		fmt.Println("Error in |func sendPasswordRecovery|: something went wrong while sending Email:", err.Error())
		c.JSON(http.StatusServiceUnavailable, "Email sending error")
		return
	}

	c.JSON(http.StatusAccepted, "Sended email") // If no errors

}

func resetPassword(c *gin.Context) { // Endpoint to reset password
	var password UserPassword // Create UserPassword

	requiredToken := c.GetHeader("Authorization") // Get mail JWT token from Authorization header

	if err := c.BindJSON(&password); err != nil { // Read password from request body
		// If error appeared:
		fmt.Println("Error in |func resetPassword|: ", err)
		c.JSON(http.StatusBadRequest, password)
		return
	}

	userEmail, wasError := DecodeEmailJWTToken(requiredToken) // Decode mail JWT token

	if wasError { // If error appeared

		if userEmail == "InValid" { // If error appeared because mail JWT token isn't valid
			fmt.Println("Error in |func resetPassword|: got InValid JWT token")
			c.JSON(http.StatusUnauthorized, "Invalid JWT token")
			return
		}

		// If error appeared because of something else
		fmt.Println("Error in |func resetPassword|: something went wrong while decoding JWT token")
		c.JSON(http.StatusServiceUnavailable, "JWT token decoding error")
		return
	}

	_, err := db.Query("UPDATE Users UserPassword = ? WHERE UserEmail = ?", password.UserPassword, userEmail) // Update User's password in DB

	if err != nil { // If error appeared
		fmt.Println("Error in |func resetPassword|:", err.Error())
		c.JSON(http.StatusServiceUnavailable, password)
		return
	}

	c.JSON(http.StatusAccepted, "Changed password") // If no errors

}

/* </Password recovery> */

/* <Variables init> */

// JWT vars (We will need them to generate and decode JWT tokens)
var jwtKey = []byte(os.Getenv("SECRET_KEY"))
var jwtEmailKey = []byte(os.Getenv("ANOTHER_SECRET_KEY"))

//"thats_mine_dumb_secret_key"
//"thats_mine_second_dumb_secret_key"

// SMTP vars (We will need them to send email with password-reset link)
var smtp_username = os.Getenv("SMTPUSER")
var smtp_password = os.Getenv("SMTPPASS")
var smtp_host_ip = os.Getenv("SMTPHOST")
var smtp_sender_mail = "from@example.com"

// Fronted url (We will need it to create password-reset link)
var frontend_url = "http://localhost:5173/password-change" // Don't forget to add pathname to your password-reset page

// MySql vars (We will need them to get access to MySql DB)
var db *sql.DB
var db_user = os.Getenv("DBUSER")
var db_password = os.Getenv("DBPASS")
var db_host_ip = "127.0.0.1:3306"

/* </Variables init> */

func main() {
	// Start MySql connection
	cfg := mysql.Config{
		User:   db_user,
		Passwd: db_password,
		Net:    "tcp",
		Addr:   db_host_ip,
		DBName: "Users_DB",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN()) // Open connection to DB

	if err != nil { // If error appeared
		fmt.Println("Couldn't open DB! Error:", err.Error())
		return
	}

	defer db.Close() // Don't forget to close DB connection after stop API service

	pingErr := db.Ping() // Ping DB to check connection

	if pingErr != nil {
		fmt.Println("Couldn't open DB! Error:", pingErr.Error())
		return
	}

	fmt.Println("Connected!") // If we connected to DB

	// Create gin router
	router := gin.Default()

	// Add CORS policy
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	// User endpoints
	router.POST("/signup", signUp)
	router.POST("/login", logIn)

	router.GET("/auth", auth)

	router.POST("/changeUN", changeUsername)
	router.POST("/changeUE", changeEmail)
	router.GET("/delUser", deleteUser)

	router.GET("/getUserDispatches", getUserDispatches)
	router.GET("/getDispatches", getDispatches)

	// Dispatch endpoints
	router.POST("/createDispatch", createDispatch)
	router.POST("/saveDispatch", saveDispatch)
	router.GET("/getDispatch", getDispatch)
	router.GET("/deleteDispatch", deleteDispatch)

	router.GET("/canEdit", canEdit)

	// Password recovery endpoints
	router.POST("/sendPasswordRecovery", sendPasswordRecovery)
	router.POST("/resetPassword", resetPassword)

	// Run API
	router.Run("localhost:8080")
}
