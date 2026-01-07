# Build Configuration Notes

## Build Methods

The SAM template uses `BuildMethod: go1.x` with `CodeUri: cmd/lambda` to build from the correct main.go file.

### Method 1: SAM's go1.x Builder (Current - Recommended)

The template.yaml uses:
```yaml
Metadata:
  BuildMethod: go1.x
  GoBuildFlags: '-ldflags="-s -w"'
Properties:
  CodeUri: cmd/lambda
```

**Configuration:**
- `CodeUri: cmd/lambda` - Points to directory containing Lambda main.go
- `GoBuildFlags` - Attempts to use optimization flags (may not always apply)
- Builds from `cmd/lambda/main.go` (Lambda entry point)

**Pros:**
- Automatic build handling by SAM
- No Makefile required
- Works in containers
- Builds from correct main.go location

**Cons:**
- Binary size may be larger if flags don't apply (~67 MB vs ~48 MB)
- GoBuildFlags support may vary by SAM version

**Usage:**
```bash
sam build
```

### Method 2: Makefile Builder (Alternative)

To use Makefile builder, change template.yaml to:
```yaml
Metadata:
  BuildMethod: makefile
```

Then ensure Makefile has target: `build-ConnectraApiFunction`

**Pros:**
- Full control over build flags
- Can use `-ldflags="-s -w"` for smaller binaries
- Can specify exact main.go path

**Cons:**
- Requires `make` in build environment
- More configuration needed

**Usage:**
```bash
sam build
```

### Method 3: Manual Build (For Testing)

Build manually and then use SAM to package:
```bash
# Build manually
make build-lambda
# or
./scripts/build.sh

# Then package with SAM (SAM will skip build if binary exists)
sam build
```

## Binary Size Optimization

The current build produces a ~67.7 MB binary. To reduce size:

1. **Use Makefile builder** with `-ldflags="-s -w"`:
   - `-s`: Omit symbol table and debug information
   - `-w`: Omit DWARF symbol table

2. **Or add build flags to go1.x method** (if supported):
   - Check SAM documentation for GoBuildFlags

## Current Build Output

- **Location**: `.aws-sam/build/ConnectraApiFunction/bootstrap`
- **Size**: ~67.13 MB (with go1.x method from cmd/lambda)
- **Size**: ~48.68 MB (with manual build using -ldflags="-s -w")
- **Main File**: `cmd/lambda/main.go` (Lambda entry point)

## Troubleshooting

### Build fails with "make not found"
- Use `BuildMethod: go1.x` (current configuration)
- Or ensure `make` is available in build environment

### Build fails with "main.go not found"
- Ensure `CodeUri` points to directory containing main.go
- Or use Makefile builder to specify exact path: `./cmd/lambda/main.go`

### Binary too large
- Use Makefile builder with `-ldflags="-s -w"`
- Or manually build and copy to artifacts directory
