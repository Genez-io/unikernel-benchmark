# Chech if script has at least 4 arguments
if [ $# -lt 4 ]; then
    #            $0      $1         $2               $3             $4
    echo "Usage: $0 <docker_ip> <unikernel_ip> <unikernel_port> <command>"
    exit 1
fi

# Open local server that will wait until a connection is made
# and then close it
nc -l -p 25565

# # Enable IP forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward && \
iptables -F && \
iptables -t nat -F && \
iptables -X && \
iptables -t nat -A PREROUTING -p tcp --dport $3 -j DNAT --to-destination $2:$3 && \
iptables -t nat -A POSTROUTING -p tcp -d $2 --dport $3 -j SNAT --to-source $1 && \
iptables -t nat -A PREROUTING -p udp --dport $3 -j DNAT --to-destination $2:$3 && \
iptables -t nat -A POSTROUTING -p udp -d $2 --dport $3 -j SNAT --to-source $1

# # Run command passed as argument
/bin/bash -c "$4"