FROM mcr.microsoft.com/devcontainers/go:bookworm

USER vscode

# Install air
RUN go install github.com/air-verse/air@latest

# Install uv and aider-chat
# We use the home directory of the vscode user to ensure they are available in the PATH
ENV UV_INSTALL_DIR="/home/vscode/.local/bin"
# Original command to instal UV
# RUN curl -LsSf https://astral.sh/uv/install.sh | sh
# Command to install Aider and UV
RUN curl -LsSf https://aider.chat/install.sh | sh
 
# Ensure the local bin is in the PATH
ENV PATH="${UV_INSTALL_DIR}:${PATH}"
