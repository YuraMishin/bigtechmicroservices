//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		_ = env.ClearInventoryCollection(ctx)
		cancel()
	})

	Describe("ListParts", func() {
		It("должен возвращать непустой список деталей после вставки тестовой детали", func() {
			_, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).ToNot(BeNil())
			Expect(resp.GetParts()).ToNot(BeEmpty())
		})
	})

	Describe("GetPart", func() {
		It("должен возвращать деталь по UUID после вставки тестовой детали", func() {
			uuid, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())

			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{Uuid: uuid})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp).ToNot(BeNil())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().GetUuid()).To(Equal(uuid))
		})
	})
})
