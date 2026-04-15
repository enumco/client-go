# enum Go API Client

Go client for the [enum](https://enum.co) API.

## Install

```bash
go get github.com/enumco/client-go
```

## Usage

```go
package main

import (
 "context"
 "fmt"
 "log"

 client "github.com/enumco/client-go"
 apiv1 "github.com/enumco/proto-gen-go/gen/enum/api/v1"
)

func main() {
 c, err := client.New(client.WithToken("your-api-token"))
 if err != nil {
  log.Fatal(err)
 }
 defer c.Close()

 resp, err := c.Organizations.ListOrganizations(context.Background(), &apiv1.ListOrganizationsRequest{})
 if err != nil {
  log.Fatal(err)
 }

 for _, org := range resp.GetOrganizations() {
  fmt.Println(org.GetName())
 }
}
```

## Services

| Field | Service |
|---|---|
| `c.Organizations` | OrganizationService |
| `c.Projects` | ProjectService |
| `c.Users` | UserService |
| `c.Kubernetes.Clusters` | KubernetesClusterService |
| `c.ObjectStorage.Users` | ObjectStorageUserService |
| `c.ObjectStorage.AccessKeys` | ObjectStorageAccessKeyService |

## Options

| Option | Description |
|---|---|
| `WithToken(token string)` | Bearer token for authentication |
| `WithAddr(addr string)` | Override the API endpoint (default: `api.enum.co:443`) |
| `WithTLS(creds credentials.TransportCredentials)` | Custom TLS configuration |
| `WithInsecure()` | Disable TLS - for local development only |

The client uses TLS with system roots by default. No extra configuration needed for production use.

## Documentation

See [enum Docs](https://docs.enum.co/api/).

## License

[MIT License](LICENSE.md)
