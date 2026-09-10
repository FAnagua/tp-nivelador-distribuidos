import socket
import logger
import protocol
import lottery
import threading
import lottery_monitor
import signal

class Server:
    def __init__(self, server_host: str, server_port: int, agency_quorum_min: int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.server_socket = None
        self.agency_quorum_min = agency_quorum_min
        self.protocol = protocol.Protocol()
        self.lottery = lottery.Lottery("bets.csv")
        self.threads : list[threading.Thread] = []
        self.lottery_monitor = lottery_monitor.Monitor(self.lottery, self.agency_quorum_min)
        self.client_sockets : list[socket.socket] = []
        self.alive = True
        signal.signal(signal.SIGTERM, self.close)

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
                    logger.info(
                        "waiting-agencies",
                        logger.LogResult.in_progress,
                        "agency-id",
                        agency_id,
                    )
                    if not self.lottery_monitor.wait_for_all_agencies():
                        logger.error(
                            "waiting-agencies",
                            logger.LogResult.fail,
                            "agency-id",
                            agency_id,
                            "all-agencies-ready",
                            False,
                        )
                        finish = True
                        continue
                    logger.info("waiting-agencies", logger.LogResult.success, "agency-id", agency_id, "all-agencies-ready", True)
                    bets_winning = self.lottery_monitor.load_bets_winners(agency_id)
                    self.protocol.sendResults(client_socket, bets_winning)
                        
                message_amount += 1
        except Exception as e:
            if not self.alive:
                logger.info("close-thread", logger.LogResult.success, "server-closed", True)
                finish = True
            else:
                logger.error(
                    action, logger.LogResult.fail, "messages-amount", message_amount
                )
                raise e
        finally:
            try:
                client_socket.shutdown(socket.SHUT_RDWR)
                client_socket.close()
                self.client_sockets.remove(client_socket)
            except Exception:
                pass


    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as self.server_socket:
            self.server_socket.bind((self.server_host, self.server_port))
            self.server_socket.listen()
            while self.alive:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = self.server_socket.accept()
                    self.client_sockets.append(client_socket)
                except Exception as e:
                    if not self.alive:
                        logger.info(action, logger.LogResult.success, "server-closed", True)
                        break
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                thread = threading.Thread(target=self._handle_client, args=(client_socket,))
                self.threads.append(thread)

                thread.start()
        
        self.wait_for_threads()

    def wait_for_threads(self):
        logger.info("wait-for-threads", logger.LogResult.in_progress)
        for thread in self.threads:
            thread.join()
            logger.info("wait-for-threads", logger.LogResult.success, "thread-id", thread.ident)
        logger.info("wait-for-threads", logger.LogResult.success, "threads-count", len(self.threads))

    def close(self, signum, frame):
        logger.info("signal-received", logger.LogResult.success, "signal", signum)
        self.alive = False
        self.lottery_monitor.close()
        self.server_socket.shutdown(socket.SHUT_RDWR)
        self.server_socket.close()
        for client_socket in self.client_sockets:
            client_socket.shutdown(socket.SHUT_RDWR)
            client_socket.close()
        
