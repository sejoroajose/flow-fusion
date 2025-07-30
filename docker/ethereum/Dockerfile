FROM ghcr.io/foundry-rs/foundry:latest

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /workspace

# Copy foundry configuration
COPY foundry.toml ./

# Install forge dependencies when container starts
CMD ["anvil"]