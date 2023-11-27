#!/bin/bash

cleanup() {
	echo "Container stopped. Removing nginx configuration."
	rm /etc/nginx/sites-enabled/heimdall.metaploy.conf
}

trap 'cleanup' SIGQUIT

"${@}" &

cp /heimdall.metaploy.conf /etc/nginx/sites-enabled

wait $!
