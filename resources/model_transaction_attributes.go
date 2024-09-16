/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type TransactionAttributes struct {
	// transaction timestamp
	CreatedAt time.Time `json:"created_at"`
	// user address
	FromAddress string `json:"from_address"`
	// id deal
	Id string `json:"id"`
	// user address
	ToAddress string `json:"to_address"`
	// exchange value
	Value int `json:"value"`
}
