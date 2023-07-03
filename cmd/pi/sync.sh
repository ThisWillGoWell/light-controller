#!/bin/bash

# Get the directory of the script
script_directory=$(dirname "$(readlink -f "$0")")

# Define the source and destination directories
source_directory="$script_directory/portal"
destination_directory="/home/pi"

# Check if the hosts.txt file exists
hosts_file="$script_directory/hosts.txt"
if [[ ! -f "$hosts_file" ]]; then
    echo "Error: hosts.txt file not found in the script directory."
    exit 1
fi

# Read the hosts from the hosts.txt file
hosts=()
while read -r name host
do
    echo "Checking availability of $name ($host)..."
    if ping -c 1 -W 1 "$host" >/dev/null; then
        echo "$name ($host) is online."

        if [[ -z $first_host ]]; then
            first_host=$host
            echo "Syncing $source_directory to $host:$destination_directory"
            rsync -avz "$source_directory/" "pi@$host:$destination_directory/"

            echo "Building binary on $name ($host)"
            ssh "pi@$host" "cd $destination_directory && make portal"

            echo "Syncing binary back to your machine"
            scp "pi@$host:$destination_directory/portal" "$script_directory/portal/"
        fi

        # Sync binary to other hosts
        echo "Syncing binary to $name ($host)"
        scp "$script_directory/portal/portal" "pi@$host:$destination_directory/"

        # Stop existing portal process
        echo "Stopping existing portal process on $name ($host)"
        ssh "pi@$host" "sudo pkill portal"

        # Start new portal binary as root
        echo "Starting new portal binary on $name ($host)"
        ssh "pi@$host" "sudo $destination_directory/portal &"

        echo "Deployment completed on $name ($host)"
        echo
    else
        echo "$name ($host) is offline. Skipping..."
        echo
    fi
done < "$hosts_file"
