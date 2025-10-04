#!/bin/bash

# A script to easily view the goHome application logs on a remote Raspberry Pi.

# --- Configuration ---
# Set the username and hostname for your Raspberry Pi.
REMOTE_USER="admin"
REMOTE_HOST="raspberrypi.local" # Or use the IP address

# Set the name of the service on the remote machine.
SERVICE_NAME="gohome.service"

# --- Script Logic ---

# Function to display help/usage information
usage() {
    echo "Usage: $0 [option]"
    echo ""
    echo "A simple wrapper to view logs for the ${SERVICE_NAME} service on ${REMOTE_HOST}."
    echo ""
    echo "Options:"
    echo "  -f, follow       Follow logs in real-time (most common)."
    echo "  -n <lines>       Show the last <lines> of the log."
    echo "  -g, grep <pattern> Filter logs for a specific pattern (e.g., 'AUTH_FAILURE')."
    echo "  -s, since <time> Show logs since a relative time (e.g., '10 min ago', '1h')."
    echo "  (no option)      View all logs in a scrollable pager (less)."
    exit 1
}

# --- Build the remote command based on user input ---

# Base command to view logs for our specific service unit
BASE_CMD="sudo journalctl -u ${SERVICE_NAME}"

REMOTE_CMD=${BASE_CMD}

# Parse command-line arguments
if [ -z "$1" ]; then
    # Default behavior: show all logs in a pager
    REMOTE_CMD="${BASE_CMD} --no-pager"
else
    case "$1" in
        -f|follow)
            REMOTE_CMD="${BASE_CMD} -f"
            ;;
        -n)
            if [[ -z "$2" || ! "$2" =~ ^[0-9]+$ ]]; then
                echo "Error: Please provide a number of lines for the -n option."
                usage
            fi
            REMOTE_CMD="${BASE_CMD} -n $2 --no-pager"
            ;;
        -g|grep)
            if [ -z "$2" ]; then
                echo "Error: Please provide a pattern to grep for."
                usage
            fi
            # We use --no-pager to ensure we get a clean stream for grep
            REMOTE_CMD="${BASE_CMD} --no-pager | grep --color=always -i '$2'"
            ;;
        -s|since)
            if [ -z "$2" ]; then
                echo "Error: Please provide a time specifier (e.g., '1h', '30 min ago')."
                usage
            fi
            REMOTE_CMD="${BASE_CMD} --since \"$2\" --no-pager"
            ;;
        *)
            echo "Error: Invalid option '$1'"
            usage
            ;;
    esac
fi

# --- Execute the command ---

echo "--- Connecting to ${REMOTE_HOST} to view logs..."
echo "--- Running command: ${REMOTE_CMD}"
echo "--- (Press Ctrl+C to exit)"
echo ""

# The -t flag allocates a pseudo-terminal, which is necessary for interactive
# commands like 'follow' and for preserving color output from grep.
ssh -t ${REMOTE_USER}@${REMOTE_HOST} "${REMOTE_CMD}"