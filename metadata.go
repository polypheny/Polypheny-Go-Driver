package polypheny

import (
	"context"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
)

// SearchEntities returns metadata about the entities in namespace. A nil pattern
// requests all entities. Document collection metadata requires a server that
// implements the Prism 2.2 document metadata fields.
func (conn *PolyphenyConn) SearchEntities(ctx context.Context, namespace string, pattern *string) (*prism.EntitiesResponse, error) {
	response, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_EntitiesRequest{
			EntitiesRequest: &prism.EntitiesRequest{
				NamespaceName: namespace,
				EntityPattern: pattern,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	entities := response.GetEntitiesResponse()
	if entities == nil {
		return nil, &ClientError{message: "expected a Prism entities response"}
	}
	return entities, nil
}
