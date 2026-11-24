package polypheny

import (
	"context"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
)

// SetClientInfoPropertiesContext sets Prism client properties for this session.
// Set CypherResultFormat to relational to retain columns and row boundaries
// for graph-only Cypher projections on supported servers.
func (conn *PolyphenyConn) SetClientInfoPropertiesContext(ctx context.Context, properties map[string]string) error {
	response, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_SetClientInfoPropertiesRequest{
			SetClientInfoPropertiesRequest: &prism.ClientInfoProperties{Properties: properties},
		},
	})
	if err != nil {
		return err
	}
	if response.GetSetClientInfoPropertiesResponse() == nil {
		return &ClientError{message: "Prism server did not confirm client information properties"}
	}
	return nil
}

// SearchNamespaces lists namespaces matching the optional name pattern and
// model. Common model values are "RELATIONAL", "DOCUMENT", and "GRAPH".
// A nil filter matches all values for that filter.
func (conn *PolyphenyConn) SearchNamespaces(ctx context.Context, pattern, namespaceType *string) (*prism.NamespacesResponse, error) {
	response, err := conn.call(ctx, &prism.Request{
		Type: &prism.Request_NamespacesRequest{
			NamespacesRequest: &prism.NamespacesRequest{
				NamespacePattern: pattern,
				NamespaceType:    namespaceType,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	namespaces := response.GetNamespacesResponse()
	if namespaces == nil {
		return nil, &ClientError{message: "expected a Prism namespaces response"}
	}
	return namespaces, nil
}

// SearchEntities returns metadata about the entities in namespace. A nil pattern
// requests all entities. Document collection metadata requires server support
// for the Prism 2.2 document fields. A graph entity is an empty marker on a
// server that implements graph entity discovery.
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
