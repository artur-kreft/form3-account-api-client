package accounts

type EnumStatus string

const (
	Pending   EnumStatus = "pending"
	Confirmed EnumStatus = "confirmed"
	Failed    EnumStatus = "failed"
)
