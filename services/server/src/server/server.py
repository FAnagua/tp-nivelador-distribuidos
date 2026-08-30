import socket
import logger
import protocol
import lottery
import threading
import lottery_monitor

_ECHO_SERVER_MESSAGE_SIZE = 1024


class Server:
    def __init__(self, server_host: str, server_port: int, agency_quorum_min: int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.agency_quorum_min = agency_quorum_min
        self.protocol = protocol.Protocol()
        self.lottery = lottery.Lottery("bets.csv")
        self.threads : list[threading.Thread] = []
        self.lottery_monitor = lottery_monitor.Monitor(self.lottery, self.agency_quorum_min)

    def _handle_client(self, client_socket):
        action = "handle-client"
        message_amount = 0
        agency_id = None
        finish = False

        try:
            logger.info(action, logger.LogResult.in_progress)
            while not finish:
                cmd = self.protocol.readByte(client_socket)
                if cmd == protocol.CMD_CLIENT_AGENCY_ID:
                    agency_id = self.protocol.readByte(client_socket)
                    logger.info(
                        "received-agency-id",
                        logger.LogResult.success,
                        "agency-id",
                        agency_id,
                    )
                    #finish = True
                elif cmd == protocol.CMD_CLIENT_BET:
                    bets = self.protocol.readBets(client_socket, agency_id)
                    logger.info(
                        "received-bets",
                        logger.LogResult.success,
                        "bets",
                        str(bets),
                    )
                    self.protocol.sendAck(client_socket)
                    self.lottery_monitor.store_bets(bets)
                elif cmd == protocol.CMD_CLIENT_FINISHED:
                    logger.info(
                        "received-finished",
                        logger.LogResult.success,
                        "agency-id",
                        agency_id,
                    )
                    finish = True
                elif cmd == protocol.CMD_CLIENT_RESULTS:
                    self.lottery_monitor.wait_for_all_agencies()
                    bets = self.lottery.load_bets()
                    bets_winning = []
                    for bet in bets:
                        if self.lottery.has_won(bet) and bet.agency_id == agency_id:
                            logger.info(
                                "bet-winner",
                                logger.LogResult.success,
                                "bet",
                                str(bet),
                            )
                            bets_winning.append(bet)
                    self.protocol.sendResults(client_socket, bets_winning)
                        
                message_amount += 1
        except Exception as e:
            logger.error(
                action, logger.LogResult.fail, "messages-amount", message_amount
            )
            raise e

    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                thread = threading.Thread(target=self._handle_client, args=(client_socket,))
                self.threads.append(thread)


                thread.start()
        
        self.wait_for_threads()

    def wait_for_threads(self):
        for thread in self.threads:
            thread.join()
