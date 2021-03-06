package cqrs_test

import (
	"fmt"
	"testing"

	"github.com/andrewwebber/cqrs"
)

func TestInMemoryEventStreamRepository(t *testing.T) {
	typeRegistry := cqrs.NewTypeRegistry()
	persistance := cqrs.NewInMemoryEventStreamRepository()
	repository := cqrs.NewRepository(persistance, typeRegistry)

	hashedPassword, err := GetHashForPassword("$ThisIsMyPassword1")
	accountID := "5058e029-d329-4c4b-b111-b042e48b0c5f"

	if err != nil {
		t.Fatal("Error: ", err)
	}

	cqrs.PackageLogger().Debugf("Get hash for user...")

	cqrs.PackageLogger().Debugf("Create new account...")
	account := NewAccount("John", "Snow", "john.snow@cqrs.example", hashedPassword, 0.0)
	account.SetID(accountID)
	err = account.ChangePassword("$ThisIsANOTHERPassword")
	if err != nil {
		t.Fatal(err)
	}
	_, err = repository.Save(account, "correlationID")
	if err != nil {
		t.Fatal(err)
	}

	accountFromHistory, err := NewAccountFromHistory(accountID, repository)
	if err != nil {
		t.Fatal(err)
	}

	if string(accountFromHistory.PasswordHash) != string(account.PasswordHash) {
		t.Fatal("Expected PasswordHash to match")
	}

	if events, err := persistance.AllIntegrationEventsEverPublished(); err != nil {
		t.Fatal(err)
	} else {
		cqrs.PackageLogger().Debugf(fmt.Sprintf("%+v", events))
	}

	correlationEvents, err := persistance.GetIntegrationEventsByCorrelationID("correlationID")
	if err != nil {
		t.Fatal(err)
	}

	if len(correlationEvents) == 0 {
		t.Fatal("Expeced correlation events")
	}

	cqrs.PackageLogger().Debugf("GetIntegrationEventsByCorrelationID")
	for _, correlationEvent := range correlationEvents {
		cqrs.PackageLogger().Debugf(fmt.Sprintf("%+v", correlationEvent))
	}
}
