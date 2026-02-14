package domain

import (
	trip_types "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	UserID            string `bson:"userID"`
	PackageSlug       string `bson:"packageSlug"`
	TotalPriceInCents float64 `bson:"totalPriceInCents"`
	Route             *trip_types.OsrmApiResponse `bson:"route"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (r *RideFareModel) toProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		TotalPriceInCents: r.TotalPriceInCents,
		PackageSlug:       r.PackageSlug,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare
	for _, f := range fares {
		protoFares = append(protoFares, f.toProto())
	}
	return protoFares
}
