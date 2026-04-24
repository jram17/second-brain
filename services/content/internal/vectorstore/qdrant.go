package vectorstore

import (
	"context"
	"time"

	"github.com/jram17/second-brain/services/content/pkg/breaker"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const collectionName = "content_embeddings"
const vectorSize = 768

type QdrantStore struct {
	client      pb.PointsClient
	collections pb.CollectionsClient
	conn        *grpc.ClientConn
}

func NewQdrantStore(addr string) (*QdrantStore, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	collectionsClient := pb.NewCollectionsClient(conn)
	pointsClient := pb.NewPointsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create collection if not exists
	_, _ = collectionsClient.Create(ctx, &pb.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     vectorSize,
					Distance: pb.Distance_Cosine,
				},
			},
		},
	})

	return &QdrantStore{
		client:      pointsClient,
		collections: collectionsClient,
		conn:        conn,
	}, nil
}

var cb = breaker.New("qdrant")

// store a vector point in Qdrant
func (q *QdrantStore) Upsert(contentId, userId string, embedding []float32, text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := cb.Execute(func() (interface{},error) {
		return q.client.Upsert(ctx, &pb.UpsertPoints{
			CollectionName: collectionName,
			Points: []*pb.PointStruct{
				{
					Id: &pb.PointId{
						PointIdOptions: &pb.PointId_Uuid{Uuid: contentId},
					},
					Vectors: &pb.Vectors{
						VectorsOptions: &pb.Vectors_Vector{
							Vector: &pb.Vector{Data: embedding},
						},
					},
					Payload: map[string]*pb.Value{
						"userId":    {Kind: &pb.Value_StringValue{StringValue: userId}},
						"contentId": {Kind: &pb.Value_StringValue{StringValue: contentId}},
						"text":      {Kind: &pb.Value_StringValue{StringValue: text}},
					},
				},
			},
		})
	})
	return err
}

// Delete removes a point
func (q *QdrantStore) Delete(contentId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := cb.Execute(func() (interface{}, error) {
		return q.client.Delete(ctx, &pb.DeletePoints{
			CollectionName: collectionName,
			Points: &pb.PointsSelector{
				PointsSelectorOneOf: &pb.PointsSelector_Points{
					Points: &pb.PointsIdsList{
						Ids: []*pb.PointId{
							{PointIdOptions: &pb.PointId_Uuid{Uuid: contentId}},
						},
					},
				},
			},
		})
	})
	return err
}

// Close - close the grpc connection
func (q *QdrantStore) Close() error {
	return q.conn.Close()
}
