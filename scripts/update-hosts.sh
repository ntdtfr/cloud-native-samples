#!/bin/bash

# insert/update hosts entry
ip_address=$1
domaine=$2
apps=("traefik" "portainer" "keycloak" "grafana" "prometheus" "alertmanager" "api")
for app in "${apps[@]}"; do
  host_name="${app}.${domaine}"
  host_entry="${ip_address} ${host_name}"

  # find existing ip adress
  matches_in_ips="$(grep -n $ip_address /etc/hosts | cut -f1 -d:)"
  if [ ! -z "$matches_in_ips" ]
  then
    # find existing instances in the host file and save the line numbers
    matches_in_hosts="$(grep -n $host_name /etc/hosts | cut -f1 -d:)"
    if [ ! -z "$matches_in_hosts" ]
    then
        echo "Updating existing hosts entry."
        # iterate over the line numbers on which matches were found
        while read -r line_number; do
            # replace the text of each line with the desired host entry
            sudo sed -i '' "${line_number}s/.*/${host_entry} /" /etc/hosts
        done <<< "$matches_in_hosts"
    else
        echo "Adding new hosts entry."
        echo "$host_entry" | sudo tee -a /etc/hosts > /dev/null
    fi    # find existing instances in the host file and save the line numbers
  else
    echo "Adding new hosts entry."
    echo "$host_entry" | sudo tee -a /etc/hosts > /dev/null
  fi
done
