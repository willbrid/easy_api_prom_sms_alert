package util

import "easy-api-prom-alert-sms/config"

// GetRecipientMembers get recipient members from recipient name and Recipient slice
func GetRecipientMembers(recipients []config.Recipient, recipientName string) []string {
	var recipient config.Recipient

	for _, recipientItem := range recipients {
		if recipientItem.Name == recipientName {
			recipient = recipientItem
			break
		}
	}

	return recipient.Members
}
