# Code Signing Docker Image

This Docker image provides cross-platform code signing capabilities for Windows binaries using Azure Trusted Signing.

## What's Included

- **Azure CLI**: For authenticating with Azure Trusted Signing service
- **Jsign 7.0**: Cross-platform code signing tool with Azure Trusted Signing support
- **osslsigncode**: For verifying signatures

## Building the Image

```bash
# Build locally
docker build -t maestro-signing:latest -f docker/signing/Dockerfile docker/signing

# Or use make
make signing-image
```

## Usage

### Standalone Signing

```bash
# Set your Azure credentials
export AZURE_CLIENT_ID="your-client-id"
export AZURE_CLIENT_SECRET="your-secret"
export AZURE_TENANT_ID="your-tenant-id"
export AZURE_CODE_SIGNING_ACCOUNT_NAME="YourAccount"
export AZURE_CERTIFICATE_PROFILE_NAME="YourProfile"

# Sign a binary
docker run --rm \
  -v $(pwd):/workspace \
  -e AZURE_CLIENT_ID \
  -e AZURE_CLIENT_SECRET \
  -e AZURE_TENANT_ID \
  maestro-signing:latest \
  -c "
    az login --service-principal -u \$AZURE_CLIENT_ID -p \$AZURE_CLIENT_SECRET --tenant \$AZURE_TENANT_ID
    TOKEN=\$(az account get-access-token --resource https://codesigning.azure.net --query accessToken -o tsv)
    jsign --storetype TRUSTEDSIGNING \\
          --keystore neu.codesigning.azure.net \\
          --storepass \$TOKEN \\
          --alias $AZURE_CODE_SIGNING_ACCOUNT_NAME/$AZURE_CERTIFICATE_PROFILE_NAME \\
          /workspace/your-binary.exe
  "
```

### GoReleaser Integration

This image is automatically used by GoReleaser when building releases. See `.goreleaser.yml` for the configuration.

## Why Docker?

Azure Trusted Signing traditionally requires Windows and SignTool.exe. This Docker image enables:

- **Cross-platform signing**: Sign Windows binaries from macOS/Linux
- **CI/CD friendly**: Works in any Docker-enabled environment
- **Reproducible**: Pinned versions of all tools
- **Faster builds**: No repeated downloads

## Technical Details

- **Base Image**: `azul/zulu-openjdk-debian:21-latest` (Jsign requires Java)
- **Jsign Version**: 7.0 (first version with Azure Trusted Signing support)
- **Certificate Lifetime**: 3 days (automatically timestamped by Jsign)
- **Signature Type**: Authenticode with SHA256

## Verifying Signatures

On Windows:
```powershell
# Right-click the .exe → Properties → Digital Signatures tab
```

On Linux/macOS:
```bash
docker run --rm -v $(pwd):/workspace maestro-signing:latest \
  -c "osslsigncode verify /workspace/your-binary.exe"
```

## Troubleshooting

### "unable to get local issuer certificate"

This is expected when using `osslsigncode verify` in the Docker container. The signature is valid; the container just doesn't have the full Microsoft root certificate chain in its trust store. Windows will verify correctly.

### "certificate has expired"

Azure Trusted Signing issues 3-day certificates. Ensure the binary has a valid timestamp (Jsign does this automatically).

### Authentication failures

Ensure your Azure service principal has the "Code Signing Certificate Profile Signer" role on your Azure Trusted Signing account.

## References

- [Azure Trusted Signing Docs](https://learn.microsoft.com/en-us/azure/trusted-signing/)
- [Jsign Documentation](https://ebourg.github.io/jsign/)
- [GoReleaser Signing Docs](https://goreleaser.com/customization/sign/)
