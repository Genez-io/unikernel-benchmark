import socket
import sys

PORT = 25565  # Port to listen
METRICS_FILE = sys.argv[1]  # Path to the file containing the metrics

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind(("", PORT))
    s.listen()
    conn, addr = s.accept()
    metrics = {}
    with conn:
        # Read contents from metrics file and send it to the client
        with open(METRICS_FILE, "r") as f:
            for line in f:
                key, value = line.strip().split("=")
                metrics[key] = value

        conn.sendall(str(metrics).encode("utf-8"))