package main

type DraftTable struct {

}

func (d *DraftTable) GetLeftSeat(seatIndex int) *DraftSeat {
	//"left" is -1
	left := seatIndex - 1
	if left < 0 {
		left += len(d.Seats)
	}
	return d.Seats[left]
}

func (d *DraftTable) GetRightSeat(seatIndex int) *DraftSeat {
	//"right" is +1
	right := seatIndex + 1
	right = right % len(d.Seats)
	return d.Seats[right]
}

func (d *DraftTable) DealPacks() {
	for r:=0; r<d.Format.NumRounds; r++ {
		for i,seat := range d.Seats {
			round := new(Round)
			round.Number = r
			round.Pack = d.Format.GeneratePack(r)
			if r%2 == 0 { 
				round.PassTo = d.GetLeftSeat(i).Packs
				round.ReceiveFrom = d.GetRightSeat(i).Packs
			} else {
				round.PassTo = d.GetRightSeat(i).Packs
				round.ReceiveFrom = d.GetLeftSeat(i).Packs
			}
			seat.Rounds <- round
		}
	}

	for _,seat := range d.Seats {
		close(seat.Rounds)
	}
}

func (d *DraftTable) RunDraft() { 
	finished := 0
	for {
		select {
		case pick := <-d.Picks:
			player := d.Log.Players[pick.SeatNum]
			player.Picks = append(player.Picks, pick)
		case <-d.Finished:
			finished += 1
		}

		if finished == len(d.Seats) {
			break
		}
	}
	close(d.Picks)
	close(d.Finished)
	//any other cleanup to do?
}