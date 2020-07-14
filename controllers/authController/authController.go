package authController

import (
	"github.com/Badrouu17/go-postgresql-api-boilerplate/database"
	"github.com/Badrouu17/go-postgresql-api-boilerplate/queries"

	"github.com/Badrouu17/go-postgresql-api-boilerplate/utils/abort"
	"github.com/Badrouu17/go-postgresql-api-boilerplate/utils/jwt"
	"github.com/Badrouu17/go-postgresql-api-boilerplate/utils/password"
	"github.com/gofiber/fiber"
)

type user struct {
	ID                   int
	Name                 string
	Email                string
	Photo                string
	Password             string
	PasswordChangedAt    int32
	PasswordResetToken   string
	PasswordResetExpires int32
}

func createSendToken(u user, ctx *fiber.Ctx) {
	token, err := jwt.SignToken(u.ID)
	if err != nil {
		abort.Msg(500, "error making token", ctx)
		return
	}

	ctx.Status(201).JSON(&fiber.Map{
		"status": 201,
		"token":  token,
		"id":     u.ID,
		"name":   u.Name,
		"email":  u.Email,
	})
}

func Signup(ctx *fiber.Ctx) {
	// getting the body the right way
	type signupInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(signupInput)
	ctx.BodyParser(input)
	// check if all needed input is here
	if input.Name == "" || input.Email == "" || input.Password == "" {
		abort.Msg(400, "you need to provide email, name and password input.", ctx)
		return
	}
	// hash password
	hashed, hachingErr := password.HashPassword(input.Password)
	if hachingErr != nil {
		abort.Err(500, hachingErr, ctx)
		return
	}
	// saving the user into the db
	results := []user{}
	err := database.DB.Select(&results, queries.InsertUser(input.Name, input.Email, hashed))
	if err != nil {
		abort.Err(500, err, ctx)
		return
	}
	// create and send the response token
	createSendToken(results[0], ctx)
}

func Login(ctx *fiber.Ctx) {
	// getting the body the right way
	type loginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(loginInput)
	ctx.BodyParser(input)
	// 1) check if all needed input is here
	if input.Email == "" || input.Password == "" {
		abort.Msg(400, "you need to provide email and password input.", ctx)
		return
	}
	//  2) Check if user exists
	result := user{}
	err := database.DB.Get(&result, queries.GetUserWithEmail(input.Email))
	if err != nil {
		abort.Msg(400, "no user with this email", ctx)
		return
	}
	// 3) check if password is correct
	if !password.CheckPasswordHash(result.Password, input.Password) {
		abort.Msg(401, "the password you enterd is wrong, please try again", ctx)
		return
	}
	// 4) If everything ok
	// create and send the response token
	createSendToken(result, ctx)
}