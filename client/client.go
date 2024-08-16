package client

import (
	"golang.org/x/exp/slog"

	"github.com/mirjalilova/api-gateway/config"
	pb "github.com/mirjalilova/api-gateway/genproto/health_sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Clients contains all the clients including MinIOClient
type Clients struct {
	MedicalRecords        pb.MedicalRecordServiceClient
	GeneticData           pb.GeneticDataServiceClient
	WearableData          pb.WearableDataServiceClient
	LifeStyle             pb.LifestyleServiceClient
	HealthRecommendations pb.RecommendationServiceClient
	Notification          pb.NotificationServiceClient
}

func NewClients() *Clients {
	cnf := config.Load()

	conn, err := grpc.NewClient("localhost"+cnf.MED_TRACK, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to connect to live streaming service: %v", "err", err)
	}

	medicalrecord := pb.NewMedicalRecordServiceClient(conn)
	geneticdata := pb.NewGeneticDataServiceClient(conn)
	wearabledata := pb.NewWearableDataServiceClient(conn)
	lifestyle := pb.NewLifestyleServiceClient(conn)
	healthrecommendations := pb.NewRecommendationServiceClient(conn)
	notification := pb.NewNotificationServiceClient(conn)

	return &Clients{
		MedicalRecords:        medicalrecord,
		GeneticData:           geneticdata,
		WearableData:          wearabledata,
		LifeStyle:             lifestyle,
		HealthRecommendations: healthrecommendations,
		Notification:          notification,
	}
}