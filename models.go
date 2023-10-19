package model

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)
   
   var (
	// JWT Secret Key
	jwtSecretKey = []byte("")
   
	// Database Connection URI
	dbURI = "mongodb://localhost:51.158.37.24232"
   
	// MongoDB Collection and Database Names
	dbName         = "english_learning"
	usersCollName  = "users"
	lessonsCollName = "lessons"
   
	// Token Expiry
	accessTokenExpiry  = time.Minute * 15
	refreshTokenExpiry = time.Hour * 24 * 7
   )
   
   type server struct {
	authInterface AuthServer
   }
   
   type AuthServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error)
   }
   
   type authService struct {
	db *mongo.Database
   }
   
   type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
   }
   
   type RegisterRequest struct {
	Email    string
	Password string
   }
   
   type RegisterResponse struct {}
   
   type LoginRequest struct {
	Email    string
	Password string
   }
   
   type LoginResponse struct {
	AccessToken  string
	RefreshToken string
   }
   
   type RefreshTokenRequest struct {
	RefreshToken string
   }
   
   type RefreshTokenResponse struct {
	AccessToken string
   }
   
   func main() {
	// Connect to the MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURI))
	if err != nil {
	 log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
   
	// Ping the MongoDB to check the connection
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
	 log.Fatal(err)
	}
   
	// Create a new gRPC server
	srv := grpc.NewServer()
   
	// Register the AuthServer implementation with the gRPC server
	authService := authService{db: client.Database(dbName)}
	authServer := &server{authInterface: &authService}
	srv.RegisterService(&server, authServer)
   
	// Start the gRPC server
	lis, err := net.Listen("tcp", ":24232")
	if err != nil {
	 log.Fatal(err)
	}
	log.Println("gRPC server is running on port 24232")
	err = srv.Serve(lis)
	if err != nil {
	 log.Fatal(err)
	}
   }
   
   func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	collection := s.db.Collection(usersCollName)
   
	// Check if the user already exists
	count, err := collection.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
	 return nil, err
	}
	if count > 0 {
	 return nil, fmt.Errorf("user already exists")
	}
   
	// Create a new user
	user := User{
	 ID:       primitive.NewObjectID(),
	 Email:    req.Email,
	 Password: req.Password,
	}
   
	// Insert the user into the database
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
	 return nil, err
	}
   
	return &RegisterResponse{}, nil
   }
   
   func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	collection := s.db.Collection(usersCollName)
   
	// Find the user by email
	var user User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
	 return nil, fmt.Errorf("invalid credentials")
	}
   
	// Verify the password
	if req.Password != user.Password {
	 return nil, fmt.Errorf("invalid credentials")
	}
   
	// Generate and sign a new JWT access token
	accessToken, err := generateAccessToken(user.ID.Hex())
	if err != nil {
	 return nil, err
	}
   
	// Generate and sign a new JWT refresh token
	refreshToken, err := generateRefreshToken(user.ID.Hex())
	if err != nil {
	 return nil, err
	}
   
	return &LoginResponse{
	 AccessToken:  accessToken,
	 RefreshToken: refreshToken,
	}, nil
   }
   
   func (s *authService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// Verify the JWT refresh token
	claims, err := verifyToken(req.RefreshToken)
	if err != nil {
	 return nil, fmt.Errorf("Invalid refresh token")
	}
   
	// Generate and sign a new JWT access token
	accessToken, err := generateAccessToken(claims.UserID)
	if err != nil {
	 return nil, err
	}
   
	return &RefreshTokenResponse{
	 AccessToken: accessToken,
	}, nil
   }
   
   func generateAccessToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
   

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(accessTokenExpiry).Unix()
   
	// Generate the token
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
	 return "", err
	}
   
	return tokenString, nil
   }
   
   func generateRefreshToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
   
	// Claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()
   
	
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
	 return "", err
	}
   
	return tokenString, nil
   }
   
   func verifyToken(tokenString string) (*jwt.MapClaims, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	 return jwtSecretKey, nil
	})
	if err != nil {
	 return nil, err
	}
   
	// Verify the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	 return &claims, nil
	}
   
	return nil, fmt.Errorf("Invalid token")
   }
   
