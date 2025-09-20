package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/api/health"
	orderV1API "github.com/YuraMishin/bigtechmicroservices/order/internal/api/order/v1"
	defClients "github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc"
	inventoryClientV1 "github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/YuraMishin/bigtechmicroservices/order/internal/client/grpc/payment/v1"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/config"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository"
	orderRepo "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/order"
	defService "github.com/YuraMishin/bigtechmicroservices/order/internal/service"
	orderService "github.com/YuraMishin/bigtechmicroservices/order/internal/service/order"
	"github.com/YuraMishin/bigtechmicroservices/platform/pkg/closer"
	pgmigrator "github.com/YuraMishin/bigtechmicroservices/platform/pkg/migrator/pg"
	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	router     *chi.Mux
	httpServer *http.Server

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	inventoryClientAdapter defClients.InventoryClient
	paymentClientAdapter   defClients.PaymentClient

	dbPool *pgxpool.Pool

	orderRepository repository.OrderRepository
	orderService    defService.OrderService
	openAPIServer   http.Handler
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Router(_ context.Context) *chi.Mux {
	if d.router == nil {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(config.AppConfig().OrderHTTP.ShutdownTimeout()))
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Get("/health", health.Handler())
		d.router = r
	}
	return d.router
}

func (d *diContainer) InventoryConn(_ context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		addr := net.JoinHostPort(config.AppConfig().OrderGRPC.InventoryGrpcHost(), config.AppConfig().OrderGRPC.InventoryGrpcPort())
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic("failed to connect to InventoryService: " + err.Error())
		}
		closer.AddNamed("Inventory gRPC conn", func(ctx context.Context) error { return conn.Close() })
		d.inventoryConn = conn
	}
	return d.inventoryConn
}

func (d *diContainer) PaymentConn(_ context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		addr := net.JoinHostPort(config.AppConfig().OrderGRPC.PaymentGrpcHost(), config.AppConfig().OrderGRPC.PaymentGrpcPort())
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic("failed to connect to PaymentService: " + err.Error())
		}
		closer.AddNamed("Payment gRPC conn", func(ctx context.Context) error { return conn.Close() })
		d.paymentConn = conn
	}
	return d.paymentConn
}

func (d *diContainer) InventoryClientAdapter(ctx context.Context) defClients.InventoryClient {
	if d.inventoryClientAdapter == nil {
		d.inventoryClientAdapter = inventoryClientV1.NewClient(inventoryV1.NewInventoryServiceClient(d.InventoryConn(ctx)))
	}
	return d.inventoryClientAdapter
}

func (d *diContainer) PaymentClientAdapter(ctx context.Context) defClients.PaymentClient {
	if d.paymentClientAdapter == nil {
		d.paymentClientAdapter = paymentClientV1.NewClient(paymentV1.NewPaymentServiceClient(d.PaymentConn(ctx)))
	}
	return d.paymentClientAdapter
}

func (d *diContainer) DBPool(ctx context.Context) *pgxpool.Pool {
	dbUser := config.AppConfig().OrderPostgresql.PostgresUser()
	dbPassword := config.AppConfig().OrderPostgresql.PostgresPassword()
	dbHost := config.AppConfig().OrderPostgresql.PostgresHost()
	dbPort := config.AppConfig().OrderPostgresql.PostgresPort()
	dbName := config.AppConfig().OrderPostgresql.PostgresDB()
	sslMode := config.AppConfig().OrderPostgresql.PostgresSSLMode()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		panic("failed to create pgxpool: " + err.Error())
	}

	closer.AddNamed("Postgres pool", func(ctx context.Context) error { d.dbPool.Close(); return nil })
	d.dbPool = pool

	return d.dbPool
}

func (d *diContainer) runMigrations(ctx context.Context) error {
	db := stdlib.OpenDB(*d.DBPool(ctx).Config().ConnConfig)
	m := pgmigrator.NewMigrator(db, config.AppConfig().OrderPostgresql.MigrationDirectory())
	return m.Up()
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		repo, err := orderRepo.NewPostgresRepository(d.DBPool(ctx))
		if err != nil {
			panic("failed to create orderRepository: " + err.Error())
		}
		d.orderRepository = repo
	}
	return d.orderRepository
}

func (d *diContainer) OrderService(ctx context.Context) defService.OrderService {
	if d.orderService == nil {
		svc, err := orderService.NewService(d.OrderRepository(ctx), d.InventoryClientAdapter(ctx), d.PaymentClientAdapter(ctx))
		if err != nil {
			panic("failed to create orderService: " + err.Error())
		}
		d.orderService = svc
	}
	return d.orderService
}

func (d *diContainer) OpenAPIServer(ctx context.Context) http.Handler {
	if d.openAPIServer == nil {
		api, err := orderV1API.NewAPI(d.InventoryConn(ctx), d.PaymentConn(ctx), d.OrderService(ctx))
		if err != nil {
			panic("failed to create orderV1API: " + err.Error())
		}
		srv, err := orderV1.NewServer(api)
		if err != nil {
			panic(err)
		}
		d.openAPIServer = srv
	}
	return d.openAPIServer
}

func (d *diContainer) HTTPServer(ctx context.Context) *http.Server {
	if d.httpServer == nil {
		addr := net.JoinHostPort(config.AppConfig().OrderHTTP.HttpHost(), config.AppConfig().OrderHTTP.HttpPort())

		r := d.Router(ctx)
		r.Mount("/", d.OpenAPIServer(ctx))

		server := &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadHeaderTimeout(),
		}

		closer.AddNamed("HTTP server", func(ctx context.Context) error { return server.Shutdown(ctx) })
		d.httpServer = server
	}

	return d.httpServer
}
