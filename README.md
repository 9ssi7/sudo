### sudo

The `sudo` package provides functionality for managing sudo-like operations, including user verification, token generation, and verification code validation.

### Installation

To install the `sudo` package, use the following `go get` command:

```bash
go get -u github.com/9ssi7/sudo
```

### Usage

Import the `sudo` package into your code:

```go
import "github.com/9ssi7/sudo"
```

### Configuration

Create a `sudo.Config` instance to configure the `sudo` service:

```go
config := sudo.Config{
    Redis:        // Your Redis service instance,
    NotifySender: // Your NotifySender function,
    Expire:       // Optional: Set the expiration time for verification codes. Default is 5 minutes.
}

sudoService := sudo.New(config)
```

### Service Methods

#### `Check`

Check verifies the validity of a given token.

```go
cmd := sudo.CheckCommand{
    UserId:   "user123",
    DeviceId: "device456",
    Token:    "your_token",
}

err := sudoService.Check(context.Background(), cmd)
```

#### `Start`

Start initiates the sudo process by generating a verification code and notifying the user.

```go
cmd := sudo.StartCommand{
    UserId:   "user123",
    DeviceId: "device456",
    Phone:    "+1234567890",
    Email:    "user@example.com",
    Locale:   "en_US",
}

token, err := sudoService.Start(context.Background(), cmd)
```

#### `Verify`

Verify validates a user's input against the generated verification code.

```go
cmd := sudo.VerifyCommand{
    UserId:      "user123",
    DeviceId:    "device456",
    VerifyToken: "generated_verify_token",
    Code:        "user_input_code",
}

accessToken, err := sudoService.Verify(context.Background(), cmd)
```

### Notifications

The package requires a notification sender function (`NotifySender`) to notify users during the sudo process.

```go
notifyFunc := func(cmd sudo.NotifyCommand) {
    // Implement your notification logic here
}

config.NotifySender = notifyFunc
```

### Error Messages

The package provides the following error messages:

- `sudo_redis_fetch_failed`
- `sudo_redis_set_failed`
- `sudo_json_marshal_failed`
- `sudo_not_found`
- `sudo_invalid_token`
- `sudo_invalid_code`
- `sudo_expired_code`
- `sudo_exceed_try_count`
- `sudo_unknown`
- `sudo_verify_started`

### Example

For a complete working examples, refer to the [recipes](github.com/9ssi7/recipes) repository.

### License

This package is licensed under the [Apache-2.0 License](LICENSE).