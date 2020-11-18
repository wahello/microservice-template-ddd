// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"context"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"robovoice-template/internal/billing/application"
	"robovoice-template/internal/billing/infrastructure/rpc"
	"robovoice-template/internal/book/application"
	"robovoice-template/internal/book/infrastructure/rpc"
	"robovoice-template/internal/book/infrastructure/store"
	"robovoice-template/internal/db"
	"robovoice-template/internal/user/application"
	"robovoice-template/internal/user/infrastructure/rpc"
	"robovoice-template/pkg/rpc"
	"robovoice-template/pkg/traicing"
	"time"
)

// Injectors from wire.go:

func InitializeAPIService(ctx context.Context) (*APIService, func(), error) {
	logger, err := InitLogger(ctx)
	if err != nil {
		return nil, nil, err
	}
	tracer, cleanup, err := InitTracer(ctx, logger)
	if err != nil {
		return nil, nil, err
	}
	clientConn, cleanup2, err := runGRPCClient(logger, tracer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	apiService, err := NewAPIService(logger, clientConn)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return apiService, func() {
		cleanup2()
		cleanup()
	}, nil
}

func InitializeUserService(ctx context.Context) (*UserService, func(), error) {
	logger, err := InitLogger(ctx)
	if err != nil {
		return nil, nil, err
	}
	service, err := NewUserApplication()
	if err != nil {
		return nil, nil, err
	}
	tracer, cleanup, err := InitTracer(ctx, logger)
	if err != nil {
		return nil, nil, err
	}
	rpcServer, cleanup2, err := runGRPCServer(logger, tracer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userServer, err := NewUserRPCServer(service, logger, rpcServer)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	userService, err := NewUserService(logger, userServer)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return userService, func() {
		cleanup2()
		cleanup()
	}, nil
}

func InitializeBillingService(ctx context.Context) (*BillingService, func(), error) {
	logger, err := InitLogger(ctx)
	if err != nil {
		return nil, nil, err
	}
	service, err := NewBillingApplication()
	if err != nil {
		return nil, nil, err
	}
	tracer, cleanup, err := InitTracer(ctx, logger)
	if err != nil {
		return nil, nil, err
	}
	rpcServer, cleanup2, err := runGRPCServer(logger, tracer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	billingServer, err := NewBillingRPCServer(service, logger, rpcServer)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	billingService, err := InitBillingService(logger, billingServer)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return billingService, func() {
		cleanup2()
		cleanup()
	}, nil
}

func InitializeBookService(ctx context.Context) (*BookService, func(), error) {
	logger, err := InitLogger(ctx)
	if err != nil {
		return nil, nil, err
	}
	store, cleanup, err := InitStore(ctx, logger)
	if err != nil {
		return nil, nil, err
	}
	bookStore, err := InitBookStore(ctx, logger, store)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	tracer, cleanup2, err := InitTracer(ctx, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	clientConn, cleanup3, err := runGRPCClient(logger, tracer)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	billingRPCClient, err := NewBillingRPCClient(ctx, logger, clientConn)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	userRPCClient, err := NewUserRPCClient(ctx, logger, clientConn)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	service, err := NewBookApplication(bookStore, billingRPCClient, userRPCClient)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	rpcServer, cleanup4, err := runGRPCServer(logger, tracer)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	bookServer, err := NewBookRPCServer(service, logger, rpcServer)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	bookService, err := NewBookService(logger, bookServer)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return bookService, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

// TODO: Move to inside package
// runGRPCServer ...
func runGRPCServer(log *zap.Logger, tracer opentracing.Tracer) (*rpc.RPCServer, func(), error) {
	return rpc.InitServer(log, tracer)
}

// TODO: Move to inside package
// runGRPCClient - set up a connection to the server.
func runGRPCClient(log *zap.Logger, tracer opentracing.Tracer) (*grpc.ClientConn, func(), error) {
	return rpc.InitClient(log, tracer)
}

// DefaultService ======================================================================================================
type DefaultService struct {
	Log *zap.Logger
}

var DefaultSet = wire.NewSet(InitLogger, InitTracer)

func InitLogger(ctx context.Context) (*zap.Logger, error) {
	viper.SetDefault("LOG_LEVEL", zap.InfoLevel)
	viper.SetDefault("LOG_TIME_FORMAT", time.RFC3339Nano)

	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return log, nil
}

func InitTracer(ctx context.Context, log *zap.Logger) (opentracing.Tracer, func(), error) {
	viper.SetDefault("TRACER_SERVICE_NAME", "ShortLink")
	viper.SetDefault("TRACER_URI", "localhost:6831")

	config := traicing.Config{
		ServiceName: viper.GetString("TRACER_SERVICE_NAME"),
		URI:         viper.GetString("TRACER_URI"),
	}

	tracer, tracerClose, err := traicing.Init(config)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := tracerClose.Close(); err != nil {
			log.Error(err.Error())
		}
	}

	return tracer, cleanup, nil
}

// InitStore return db
func InitStore(ctx context.Context, log *zap.Logger) (*db.Store, func(), error) {
	var st db.Store
	db2, err := st.Use(ctx, log)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := db2.Store.Close(); err != nil {
			log.Error(err.Error())
		}
	}

	return db2, cleanup, nil
}

// APIService ==========================================================================================================
type APIService struct {
	Log *zap.Logger

	ClientRPC *grpc.ClientConn
}

var APISet = wire.NewSet(DefaultSet, runGRPCClient, NewAPIService)

func NewAPIService(log *zap.Logger, clientRPC *grpc.ClientConn) (*APIService, error) {
	return &APIService{
		Log:       log,
		ClientRPC: clientRPC,
	}, nil
}

// UserService =========================================================================================================
type UserService struct {
	Log *zap.Logger

	userRPCServer *user_rpc.UserServer
}

var UserSet = wire.NewSet(DefaultSet, runGRPCServer, NewUserService, NewUserApplication, NewUserRPCServer)

func NewUserApplication() (*user.Service, error) {
	userService, err := user.New()
	if err != nil {
		return nil, err
	}

	return userService, nil
}

func NewUserRPCClient(ctx context.Context, log *zap.Logger, rpcClient *grpc.ClientConn) (user_rpc.UserRPCClient, error) {
	userService, err := user_rpc.Use(ctx, rpcClient)
	if err != nil {
		return nil, err
	}

	return userService, nil
}

func NewUserRPCServer(userService *user.Service, log *zap.Logger, serverRPC *rpc.RPCServer) (*user_rpc.UserServer, error) {
	userRPCServer, err := user_rpc.New(serverRPC, log, userService)
	if err != nil {
		return nil, err
	}

	return userRPCServer, nil
}

func NewUserService(log *zap.Logger, userRPCServer *user_rpc.UserServer) (*UserService, error) {
	return &UserService{
		Log: log,

		userRPCServer: userRPCServer,
	}, nil
}

// BillingService ======================================================================================================
type BillingService struct {
	Log *zap.Logger

	billingRPCServer *billing_rpc.BillingServer
}

var BillingSet = wire.NewSet(DefaultSet, runGRPCServer, InitBillingService, NewBillingApplication, NewBillingRPCServer)

func NewBillingApplication() (*billing.Service, error) {
	billingService, err := billing.New()
	if err != nil {
		return nil, err
	}

	return billingService, nil
}

func NewBillingRPCClient(ctx context.Context, log *zap.Logger, rpcClient *grpc.ClientConn) (billing_rpc.BillingRPCClient, error) {
	billingService, err := billing_rpc.Use(ctx, rpcClient)
	if err != nil {
		return nil, err
	}

	return billingService, nil
}

func NewBillingRPCServer(billingService *billing.Service, log *zap.Logger, serverRPC *rpc.RPCServer) (*billing_rpc.BillingServer, error) {
	billingRPCServer, err := billing_rpc.New(serverRPC, log, billingService)
	if err != nil {
		return nil, err
	}

	return billingRPCServer, nil
}

func InitBillingService(log *zap.Logger, billingRPCServer *billing_rpc.BillingServer) (*BillingService, error) {
	return &BillingService{
		Log: log,

		billingRPCServer: billingRPCServer,
	}, nil
}

// BookService =========================================================================================================
type BookService struct {
	Log *zap.Logger

	bookRPCServer *book_rpc.BookServer
}

var BookSet = wire.NewSet(DefaultSet, runGRPCServer, runGRPCClient, NewBookApplication, NewBillingRPCClient, NewUserRPCClient, InitStore, InitBookStore, NewBookRPCServer, NewBookService)

// InitMetaStore
func InitBookStore(ctx context.Context, log *zap.Logger, conn *db.Store) (*store.BookStore, error) {
	st := store.BookStore{}
	bookStore, err := st.Use(ctx, log, conn)
	if err != nil {
		return nil, err
	}

	return bookStore, nil
}

func NewBookApplication(store2 *store.BookStore, billingRPC billing_rpc.BillingRPCClient, userRPC user_rpc.UserRPCClient) (*book.Service, error) {
	bookService, err := book.New(store2, userRPC, billingRPC)
	if err != nil {
		return nil, err
	}

	return bookService, nil
}

func NewBookRPCClient(ctx context.Context, log *zap.Logger, rpcClient *grpc.ClientConn) (book_rpc.BookRPCClient, error) {
	bookService, err := book_rpc.Use(ctx, rpcClient)
	if err != nil {
		return nil, err
	}

	return bookService, nil
}

func NewBookRPCServer(bookService *book.Service, log *zap.Logger, serverRPC *rpc.RPCServer) (*book_rpc.BookServer, error) {
	bookRPCServer, err := book_rpc.New(serverRPC, log, bookService)
	if err != nil {
		return nil, err
	}

	return bookRPCServer, nil
}

func NewBookService(log *zap.Logger, bookRPCServer *book_rpc.BookServer) (*BookService, error) {
	return &BookService{
		Log: log,

		bookRPCServer: bookRPCServer,
	}, nil
}
