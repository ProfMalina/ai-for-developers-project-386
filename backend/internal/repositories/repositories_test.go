package repositories

import "testing"

func TestRepositoryConstructorsReturnInstances(t *testing.T) {
	if NewBookingRepository() == nil {
		t.Fatal("expected booking repository instance")
	}
	if NewEventTypeRepository() == nil {
		t.Fatal("expected event type repository instance")
	}
	if NewOwnerRepository() == nil {
		t.Fatal("expected owner repository instance")
	}
	if NewSlotGenerationConfigRepository() == nil {
		t.Fatal("expected slot generation config repository instance")
	}
	if NewTimeSlotRepository() == nil {
		t.Fatal("expected time slot repository instance")
	}
}
