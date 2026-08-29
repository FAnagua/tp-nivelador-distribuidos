import socket
import safe_socket

import lottery

CMD_CLIENT_AGENCY_ID = 0x01
CMD_CLIENT_BET = 0x02
CMD_SERVER_ACK = 0x03
CMD_CLIENT_FINISHED = 0x04
CMD_CLIENT_RESULTS = 0x05
CMD_SERVER_RESULTS = 0x06

class Protocol:
    def __init__(self):
        pass

    def serializeByte(self, value: int) -> bytes:
        return value.to_bytes(1)

    def serialize2Bytes(self, value: int) -> bytes:
        return value.to_bytes(2, byteorder="big")

    def serialize4Bytes(self, value: int) -> bytes:
        return value.to_bytes(4, byteorder="big")

    def serializeString(self, value: str) -> bytes:
        encoded = value.encode("utf-8")
        length = len(encoded)
        return self.serialize2Bytes(length) + encoded

    def serializeBet(self, bet: lottery.Bet) -> bytes:
        data = b""
        data += self.serializeString(bet.first_name)
        data += self.serializeString(bet.last_name)
        data += self.serialize4Bytes(bet.document)
        data += self.serializeString(bet.birthdate)
        data += self.serialize4Bytes(bet.number)
        return data

    def readByte(self, socket: socket.socket) -> int:
        data = safe_socket.recv_all(socket, 1)
        return int.from_bytes(data)

    def read2Bytes(self, socket: socket.socket) -> int:
        data = safe_socket.recv_all(socket, 2)
        return int.from_bytes(data, byteorder="big")

    def read4Bytes(self, socket: socket.socket) -> int:
        data = safe_socket.recv_all(socket, 4)
        return int.from_bytes(data, byteorder="big")

    def readString(self, socket: socket.socket) -> str:
        length = self.read2Bytes(socket)
        data = safe_socket.recv_all(socket, length)
        return data.decode("utf-8")

    def readBet(self, socket: socket.socket, agency_id: int) -> lottery.Bet:
        first_name = self.readString(socket)
        last_name = self.readString(socket)
        document = self.read4Bytes(socket)
        birth_date = self.readString(socket)
        number = self.read4Bytes(socket)
        return lottery.Bet(agency_id, first_name, last_name, document, birth_date, number)

    def readBets(self, socket: socket.socket, agency_id: int) -> list[lottery.Bet]:
        num_bets = self.read2Bytes(socket)
        bets = []
        for _ in range(num_bets):
            bet = self.readBet(socket, agency_id)
            bets.append(bet)
        return bets

    def sendAck(self, socket: socket.socket) -> None:
        data = self.serializeByte(CMD_SERVER_ACK)
        safe_socket.send_all(socket, data)

    def sendResults(self, socket: socket.socket, results: list[lottery.Bet]) -> None:
        data = self.serializeByte(CMD_SERVER_RESULTS)
        num_bets = len(results)
        data += self.serialize2Bytes(num_bets)
        socket.sendall(data)
    
        for bet in results:
            bet_winning_data = self.serializeBet(bet)
            socket.sendall(bet_winning_data)
