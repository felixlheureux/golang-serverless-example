package domain

import "time"

type Challenge struct {
	ChallengeID        string    `db:"challenge_id"`
	EthereumAddressHex string    `db:"ethereum_address"`
	Challenge          string    `db:"challenge"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}
