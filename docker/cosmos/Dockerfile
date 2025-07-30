FROM cosmwasm/wasmd:v0.50.0

# Install dependencies
RUN apk add --no-cache \
    curl \
    jq \
    bash

# Install Rust and Cargo for contract building
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"
RUN rustup target add wasm32-unknown-unknown

# Set working directory
WORKDIR /root

# Expose ports
EXPOSE 26657 1317 9090 9091

# Default command will be overridden by docker-compose
CMD ["wasmd", "start"]