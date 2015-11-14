package matching

import (
	"fmt"
	"log"
	"strings"
)

type Transaction struct {
	Amount         string
	Currency       string
	Created        string
	Description    string
	Id             string
	VendorMatchKey string
}

type Email struct {
	MessageId         string
	Subject           string
	Sender            string
	BodyHtml          string
	BodyHtmlBase64Url string
	VendorMatchKey    string
}

type Match struct {
	Transaction *Transaction
	Email       *Email
}

type Matcher struct {
	transactions map[string]*Transaction
	emails       map[string]*Email
}

func NewMatcher() *Matcher {
	return &Matcher{
		transactions: make(map[string]*Transaction),
		emails:       make(map[string]*Email),
	}
}

func (m *Matcher) MatchTransaction(t *Transaction) *Match {
	for id, r := range m.emails {
		if isMatch(t, r) {
			delete(m.emails, id)
			log.Printf("Matched transaction to receipt: %v --> %v\n", t, r)
			return &Match{
				Transaction: t,
				Email:       r,
			}
		}
	}

	log.Printf("Transaction queued: %v\n", t)
	m.transactions[t.Id] = t
	return nil
}

func (m *Matcher) MatchEmail(e *Email) *Match {
	for id, t := range m.transactions {
		if isMatch(t, e) {
			delete(m.transactions, id)
			log.Printf("Matched email to transaction: %v --> %v\n", e, t)
			return &Match{
				Transaction: t,
				Email:       e,
			}
		}
	}

	log.Printf("Receipt queued: %v\n", e)
	m.emails[e.MessageId] = e
	return nil
}

func isMatch(t *Transaction, e *Email) bool {
	return e.VendorMatchKey == t.VendorMatchKey && strings.Contains(e.BodyHtml, t.Amount)
}

func (t *Transaction) String() string {
	return fmt.Sprintf("id=%v vendor=%v amount=%v description=%v", t.Id, t.VendorMatchKey, t.Amount, t.Description)
}

func (e *Email) String() string {
	return fmt.Sprintf("message id=%v vendor=%s sender=%v subject=%v", e.MessageId, e.VendorMatchKey, e.Sender, e.Subject)
}
