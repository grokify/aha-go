# Authentication

## Getting Your API Key

1. Log in to your Aha.io account
2. Go to **Settings** → **Account** → **Personal** → **Developer**
3. Click **Generate API key**
4. Copy and save your API key securely

!!! warning "Keep Your API Key Secret"
    Never commit your API key to version control. Use environment variables or a secrets manager.

## Configuration

### Environment Variables

Set these environment variables:

```bash
export AHA_SUBDOMAIN="yourcompany"    # Your Aha.io subdomain
export AHA_API_KEY="your-api-key"      # Your API key
```

The subdomain is the part before `.aha.io` in your Aha URL (e.g., `yourcompany.aha.io`).

### Using a .env File

Create a `.env` file (add to `.gitignore`):

```bash
AHA_SUBDOMAIN=yourcompany
AHA_API_KEY=your-api-key
```

Load with [direnv](https://direnv.net/) or your preferred method.

### Shell Configuration

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
export AHA_SUBDOMAIN="yourcompany"
export AHA_API_KEY="your-api-key"
```

## Programmatic Configuration

### Using Environment Variables (Recommended)

```go
client, err := aha.NewClient()
if err != nil {
    log.Fatal(err)
}
```

The client automatically reads from `AHA_SUBDOMAIN` and `AHA_API_KEY`.

### Explicit Configuration

```go
client, err := aha.NewClientWithConfig(aha.Config{
    Subdomain: "yourcompany",
    APIKey:    "your-api-key",
})
```

### From Config Struct

```go
cfg := aha.ConfigFromEnv()
// Modify if needed
cfg.Subdomain = "different-subdomain"

client, err := aha.NewClientWithConfig(cfg)
```

## CLI Configuration

The CLI uses the same environment variables:

```bash
# Set credentials
export AHA_SUBDOMAIN="yourcompany"
export AHA_API_KEY="your-api-key"

# Verify connection
aha user me
```

## Verifying Your Setup

### With Code

```go
ctx := context.Background()
user, err := client.GetCurrentUser(ctx)
if err != nil {
    log.Fatal("Authentication failed:", err)
}
fmt.Printf("Authenticated as: %s %s\n", user.FirstName, user.LastName)
```

### With CLI

```bash
aha user me
```

Expected output:

```
ID:    user-123
Name:  John Doe
Email: john@company.com
```

## Troubleshooting

### "unauthorized" Error

- Verify your API key is correct
- Check that the subdomain matches your Aha.io URL
- Ensure the API key hasn't been revoked

### "not found" Error

- Confirm you have access to the product/resource
- Verify the product key or ID is correct

### Environment Variables Not Loading

- Ensure variables are exported (`export VAR=value`, not just `VAR=value`)
- Check your shell loaded the configuration (`echo $AHA_SUBDOMAIN`)
- Restart your terminal or source your config file
