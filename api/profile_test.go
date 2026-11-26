package api

import (
	"context"
	"os"
	"reflect"
	"testing"
)

func TestDecodeUserProfile(t *testing.T) {
	input := "FAAQAAwACAAEAGYAZwBlAGMAXwBbAFcAUwBPAEsARwBDAD8AOwA3ADMALwArACcAIwAfABsAFwATAA8ACwAHAAMAaQBoAGIAXgBaAFYAUgBOAEoARgBCAD4AOgA2ADIALgAqACYAIgAeABoAFgASAA4ACgAGAAIAYQBdAFkAVQBRAE0ASQBFAEEAPQA5ADUAMQAtACkAJQAhAB0AGQAVABEADQAJAAUAAQA="
	output, err := decodeUserProfile(input)
	if err != nil {
		t.Fatalf("decodeUserProfile returned an error: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("decodeUserProfile returned an empty slice")
	}
	t.Logf("output length: %d \n", len(output))
	t.Logf("output: %+v", output)
}

func TestEncodeUserProfile(t *testing.T) {
	input := "bAVrBWoFaQVoBWQFYwViBV4FXQVcBVgFVwVWBVUFVAVTBU8FTgVKBUkFSAVEBUMFQgUyBTEFLQUsBSsFKgUpBSUFJAUjBR8FHgUdBRwFGwUXBRYFEgURBRAFDAULBQoFBgUFBQQFAAX/BPsE+gT5BNQE0wTPBM4EzQTJBMgExwTGBMIEwQTABLwEuwS6BLYEtQS0BLAErwSrBKoEpgSlBKQEowSiBJ4EnQScBJsElwSWBJIEkQSQBIwEiwSKBIkEhQSEBIMEggSBBIAEfgR9BHwEdgR1BGoEaQRlBGQEYARfBFsEWgRWBFUEUQRQBEwESwRKBEkERQREBEAEPwQ7BDoENgQ1BDEEMAQsBCsEJgQlBCEEIAQcBBsEGgQZBBUEFAQQBA8ECwQKBAIEAQQABPwD+wP3A/YD5gPlA+ED4APcA9sD1wPWA9ED0APMA8sDxwPGA70DvAO4A7cDswOyA7EDsAOsA6sDpwOmA6IDoQOdA5wDmwOaA5YDlQOUA5MDkgORA5ADZANjA18DXgNVA1QDUANPA0sDSgNGA0UDQQNAAz8DPgM9AzkDOAM0AzMDLwMuAy0DLAMoAycDIwMiAx4DHQMcAxsDFwMWAxIDEQMNAwwDAwMCA/4C/QL5AvgC9wL2AvIC8QLtAuwC6ALnAt4C3QLZAtgC1ALTAs8CzgLKAskCxQLEArQCswKvAq4CqgKpAqUCpAKgAp8CmwKaApYClQJzAnICcQJwAmwCawJnAmYCYgJhAl0CXAJYAlcCVgJVAlECUAJMAksCRwJGAkICQQI9AjwCOAI3Ai4CLQIpAigCJAIjAhACDwILAgoCBgIFAgQCAwL6AfkB9QH0AfAB7wHrAeoB5gHlAeEB4AHeAd0B2QHYAdQB0wHPAc4BygHJAcUBxAHAAb8BuwG6AbYBtQGxAbABrwGuAaoBqQGlAaQBoAGfAZsBmgGWAZUBlAGTAY8BjgGKAYkBgAF/AXsBegF2AXUBbAFrAU4BTQEzATIBLgEtASgBIwEeARkBFAEKAQUBAAHxAOwA4gC/ALUAqwCmAKEAnQCcAJgAlwCSAI0AiACDAH4AegB5AHYAbwBqAGYAFAAQAAwACAAEAG8FbQVnBWUFYQVgBV8FWwVZBVIFUAVNBUsFRwVFBTAFLgUoBSYFIgUhBSAFGgUYBRUFEwUNBQkFBwUDBQEF/gT8BNIE0ATMBMoExQTEBMMEvQS5BLgEtwSzBLEErgSsBKkEpwSfBJoEmASVBJMEjwSOBI0EhgR/BHsEegR5BHgEdwRrBGgEZgRjBGEEXgRdBFwEWQRXBFQEUgRPBE0ESARGBEEEPgQ9BDwEOQQ3BDQEMgQvBC0EKQQoBCcEIgQfBB4EHQQYBBYEEwQRBAwEBAQDBP8D/QP6A+QD4gPfA90D2gPYA9QD0wPSA80DygPJA8gDwAO+A7sDuQO2A7QDrwOuA60DqAOlA6MDoAOeA5kDmAOXA48DjgONA4wDiwNnA2UDYgNgA1gDUwNSA1EDTgNMA0kDRwNEA0IDPAM6AzcDNQMyAzADKwMpAyYDJQMkAyEDHwMaAxgDFQMTAxADDgMGAwQDAQP/AvwC+gL1AvMC8ALuAusC6QLhAt8C3ALaAtcC1QLSAtACzQLLAsgCxgKyArACrQKrAqgCpgKjAqECngKcApkClwJ2AnQCbwJuAm0CagJoAmUCYwJgAl4CWwJZAlQCUgJPAk0CSgJIAkUCQwJAAj4COwI5AjECLwIsAioCJwIlAhMCEQIMAgkCCAIHAv0B+wH2AfMB8gHxAe4B7AHpAecB5AHiAdwB2gHXAdYB1QHSAdABzQHLAcgBxwHGAcEBvgG8AbkBuAG3AbIBrQGsAasBpgGjAaEBngGcAZkBlwGSAZEBkAGNAYsBgwGBAX4BfAF5AXcBbwFtAVEBTwExAS8BKgEpASUBJAEgAR8BGwEaARYBFQEMAQsBBwEGAQIBAQHzAPIA7gDtAOQA4wDBAMAAtwC2AK0ArACoAKcAowCiAJ4AmQCUAJMAjwCOAIoAiQCFAIQAgAB/AHsAdQB0AHEAcABsAGsAZwBlAGMAXwBbAFcAUwBPAEsARwBDAD8AOwA3ADMALwArACcAIwAfABsAFwATAA8ACwAHAAMAbgVmBVoFUQVMBUYFLwUnBRkFFAUPBQ4FCAUCBf0E0QTLBL8EvgSyBK0EqAShBKAEmQSUBIgEhwRtBGwEZwRiBFgEUwROBEcEQwRCBDgEMwQuBCQEIwQXBBIEDgQNBP4D+QP4A+MD3gPZA88DzgPFA8QDwwPCA8EDvwO6A7UDqgOpA6QDnwNmA2EDVwNWA00DSANDAzsDNgMxAyoDIAMZAxQDDwMFAwAD+wL0Au8C6gLgAtsC1gLRAswCxwKxAqwCpwKiAp0CmAJ1AmkCZAJfAloCUwJOAkkCRAI/AjoCMAIrAiYCEgIOAg0CAgIBAgAC/wH+AfwB+AH3Ae0B6AHjAdsB0QHMAcMBwgG9AbQBswGoAacBogGdAZgBjAGCAX0BeAFuAWoBaQFoAWcBZgFlAWQBYwFiAWEBYAFfAV4BXQFcAVsBWgFZAVgBVwFWAVUBVAFTAVIBUAEwASwBKwEnASYBIgEhAR0BHAEYARcBDgENAQkBCAEEAQMB9QD0APAA7wDmAOUAwwDCALkAuACvAK4AqgCpAKUApACgAJ8AmwCaAJYAlQCRAJAAjACLAIcAhgCCAIEAfQB4AHcAcwByAG4AbQBpAGgAYgBeAFoAVgBSAE4ASgBGAEIAPgA6ADYAMgAuACoAJgAiAB4AGgAWABIADgAKAAYAAgB0BHMEcgRxBHAEbwRuBIUDhAODA4IDgQOAA38DfgN9A3wDewN6A3kDeAN3A3YDdQN0A3MDcgNxA3ADbwNuA20DbANrA2oDaQNoA10DXANbA1oDWQOUApMCkgKRApACjwKOAo0CjAKLAooCiQKIAocChgKFAoQCgwKCAoECgAJ/An4CfQJ8AjYCNQI0AjMCMgKIAYcBhgGFAYQB/wD+AP0A/AD7AOsA6gDpAOgA5wC0ALMAsgCxALAACQQIBAcEBgQFBIoDiQOIA4cDhgNhAF0AWQBVAFEATQBJAEUAQQA9ADkANQAxAC0AKQAlACEAHQAZABUAEQANAAkABQABAA=="
	output, err := decodeUserProfile(input)
	if err != nil {
		t.Fatalf("decodeUserProfile returned an error: %v", err)
	}
	encoded := encodeUserProfile(output)
	if encoded != input {
		t.Fatalf("encodeUserProfile did not return the original input. Expected: %s, Got: %s", input, encoded)
	}
}

func TestDecodeProfileRLE(t *testing.T) {
	got, err := decodeProfileRLE([]int{2, 7, 1, 3}, 3)
	if err != nil {
		t.Fatalf("decodeProfileRLE returned an error: %v", err)
	}
	if want := []int{7, 7, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("decodeProfileRLE = %v, want %v", got, want)
	}

	if _, err := decodeProfileRLE([]int{2, 7}, 3); err == nil {
		t.Fatal("decodeProfileRLE accepted RLE data shorter than the expected size")
	}
	if _, err := decodeProfileRLE([]int{2, 7, 2, 3}, 3); err == nil {
		t.Fatal("decodeProfileRLE accepted RLE data longer than the expected size")
	}
	if got, err := decodeProfileRLE(nil, 2); err != nil || !reflect.DeepEqual(got, []int{0, 0}) {
		t.Fatalf("decodeProfileRLE(nil, 2) = %v, %v; want [0 0], nil", got, err)
	}
}

func TestProfileSanitize(t *testing.T) {
	var profile Profile
	profile.Name = "test"
	profile.Server = 1
	profile.Data.Cards.Ids = encodeUserProfile([]uint16{1001, 1002, 1003})
	profile.Data.Cards.Eps = []int{2, 1, 1, 0}
	profile.Data.Cards.Arts = []int{3, 1}
	profile.Data.Cards.Levels = []int{1, 30, 2, 1}
	profile.Data.Cards.Skills = []int{3, 2}
	profile.Data.Cards.Trains = []int{1, 1, 2, 0}
	profile.Data.Cards.Masters = []int{2, 5, 1, 0}

	got, err := profile.Sanitize()
	if err != nil {
		t.Fatalf("Sanitize returned an error: %v", err)
	}
	want := []CardInfo{
		{ID: 1001, Eps: 1, Arts: 1, Level: 30, Skill: 2, Train: true, Master: 5},
		{ID: 1002, Eps: 1, Arts: 1, Level: 1, Skill: 2, Master: 5},
		{ID: 1003, Arts: 1, Level: 1, Skill: 2, Master: 0},
	}
	if !reflect.DeepEqual(got.Cards, want) {
		t.Fatalf("Sanitize cards = %+v, want %+v", got.Cards, want)
	}
}

// Before running this test, make sure to set the environment variables BESTDORI_API_USER and BESTDORI_API_PSK with valid credentials.
func TestProfileSanitizeFull(t *testing.T) {
	client, err := NewApiWithLogin(context.TODO(), nil, os.Getenv("BESTDORI_API_USER"), os.Getenv("BESTDORI_API_PSK"))
	if err != nil {
		t.Fatalf("NewApiWithLogin returned an error: %v", err)
	}
	profile, err := client.GetUserProfile(context.TODO())
	if err != nil {
		t.Fatalf("GetUserProfile returned an error: %v", err)
	}
	sanitized, err := profile.Profiles[0].Sanitize()
	if err != nil {
		t.Fatalf("Sanitize returned an error: %v", err)
	}
	err = saveToFile(sanitized, "sanitized_profile.json")
	if err != nil {
		t.Fatalf("saveToFile returned an error: %v", err)
	}
}
