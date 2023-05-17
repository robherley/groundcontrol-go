# groundcontrol-go 🚀

An unoffical Go client for [GroundControl](https://groundcontrol.sh/)

## Setup

```
go get github.com/robherley/groundcontrol-go
```

## Usage

```go
client := groundcontrol.New("<project>", "<apikey>")

// check globally enabled
client.IsFeatureFlagEnabled(ctx, "my-feature-flag")

// check enablement for single actor
client.IsFeatureFlagEnabled(ctx, "my-feature-flag", groundcontrol.Actor("alice"))

// or multiple actors
client.IsFeatureFlagEnabled(ctx, "my-feature-flag", groundcontrol.Actor("alice"), groundcontrol.Actor("bob"))
```

### Options

#### `WithBaseURL`

Sets the base URL of the client:

```go
client := groundcontrol.New("<project>", "<apikey>", groundcontrol.WithBaseURL("http://localhost:8080"))
```

#### `WithHTTPClient`

Sets the underlying `net/http.Client`:

```go
// e.g. hashicorp/go-retryablehttp
retryClient := retryablehttp.NewClient().StandardClient()

client := groundcontrol.New("<project>", "<apikey>", groundcontrol.WithHTTPClient(retryClient))
```
