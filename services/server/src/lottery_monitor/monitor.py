import threading
from collections.abc import Iterator
import lottery

class Monitor:
    def __init__(self, lottery: lottery.Lottery, agency_quorum_min: int):
        self.lottery = lottery
        self.lock = threading.Lock()
        self.barrier = threading.Barrier(agency_quorum_min)

    def wait_for_all_agencies(self) -> bool:
        try:
            self.barrier.wait()
        except threading.BrokenBarrierError:
            return False

        return True

    def store_bets(self, bets: list[lottery.Bet]):
        with self.lock:
            self.lottery.store_bets(bets)

    def close(self):
        self.barrier.abort()
