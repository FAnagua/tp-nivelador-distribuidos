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

    def load_bets_winners(self, agency_id: int) -> list[lottery.Bet]:
        bets_winning = []
        with self.lock:
            bets = self.lottery.load_bets()
            for bet in bets:
                if self.lottery.has_won(bet) and bet.agency_id == agency_id:
                    bets_winning.append(bet)
        return bets_winning

    def close(self):
        self.barrier.abort()
