package client

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	apiv1 "github.com/enumco/proto-gen-go/gen/enum/api/v1"
)

const defaultAddr = "api.enum.co:443"

type KubernetesClient struct {
	Clusters apiv1.KubernetesClusterServiceClient
}

type ObjectStorageClient struct {
	Users      apiv1.ObjectStorageUserServiceClient
	AccessKeys apiv1.ObjectStorageAccessKeyServiceClient
}

type Client struct {
	Organizations apiv1.OrganizationServiceClient
	Projects      apiv1.ProjectServiceClient
	Users         apiv1.UserServiceClient
	Kubernetes    KubernetesClient
	ObjectStorage ObjectStorageClient

	conn *grpc.ClientConn
}

// c, err := client.New(client.WithToken("..."))
func New(opts ...ClientOption) (*Client, error) {
	o := &options{addr: defaultAddr}
	for _, opt := range opts {
		opt(o)
	}

	if o.creds == nil {
		o.creds = credentials.NewTLS(&tls.Config{})
	}

	dialOpts := append(o.dialOpts, grpc.WithTransportCredentials(o.creds))

	conn, err := grpc.NewClient(o.addr, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Organizations: apiv1.NewOrganizationServiceClient(conn),
		Projects:      apiv1.NewProjectServiceClient(conn),
		Users:         apiv1.NewUserServiceClient(conn),
		Kubernetes: KubernetesClient{
			Clusters: apiv1.NewKubernetesClusterServiceClient(conn),
		},
		ObjectStorage: ObjectStorageClient{
			Users:      apiv1.NewObjectStorageUserServiceClient(conn),
			AccessKeys: apiv1.NewObjectStorageAccessKeyServiceClient(conn),
		},
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
