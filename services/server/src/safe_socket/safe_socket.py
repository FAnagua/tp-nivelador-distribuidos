import socket

# TODO: Complete with a short-read/short-write tolerant implementation


def recv_all(socket: socket.socket, size):
    data = b""
    while len(data) < size:
        bytes_read = socket.recv(size - len(data))
        if not bytes_read:
            raise ConnectionError("Socket connection closed before receiving all data")
        data += bytes_read
    return data


def send_all(socket: socket.socket, bytes):
    total_bytes = len(bytes)

    while total_bytes > 0:
        bytes_write = socket.send(bytes)
        if bytes_write == 0:
            raise ConnectionError("Socket connection closed before sending all data")
        total_bytes -= bytes_write
        bytes = bytes[bytes_write:]
