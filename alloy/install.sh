#!/bin/bash
LOKI_URL="http://localhost"
curl https://raw.githubusercontent.com/letwir/config/refs/heads/main/alloy/linux.alloy | sed "s@LOKI_URL@$LOKI_URL@g" | sudo tee /etc/alloy/config.alloy
exit 0
