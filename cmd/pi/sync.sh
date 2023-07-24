#!/bin/bash

# Get the directory of the script
script_directory=$(dirname "$(readlink -f "$0")")

# Define the source and destination directories
source_directory="$script_directory/portal"
destination_directory="/home/pi"

rm $script_directory/portal/portal
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
            ssh -n "pi@$host" "cd $destination_directory && make portal"

            echo "Syncing binary back to your machine"
            scp "pi@$host:$destination_directory/portal" "$script_directory/portal/"
        fi
        #Stop existing portal process
        echo "Stopping existing portal process on $name ($host)"
        ssh -n "pi@$host" "sudo pkill portal"

        # Sync binary to other hosts
        echo "Syncing binary to $name ($host)"
        scp "$script_directory/portal/portal" "pi@$host:$destination_directory/"

        # Start new portal binary as root
        echo "Starting new portal binary on $name ($host)"
        ssh "pi@$host" << EOF > /dev/null 2>&1
    sudo nohup $destination_directory/portal > /dev/null 2>&1 &
    exit
EOF
        echo "Deployment completed on $name ($host)"
        echo
    else
        echo "$name ($host) is offline. Skipping..."
        echo
    fi
done < "$hosts_file"

rm -f "$script_directory/portal/portal"