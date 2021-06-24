package license

import (
	"time"

	"github.com/rs/zerolog/log"
)

const CheckInterval = 12
const MaxConsecutiveErrors = 6

type LicenseChecker struct {
	ProductPermalink string
	LicenseKey       string
	Variant          string
	exitChan         chan bool
}

func NewLicenseChecker(product, key, variant string) *LicenseChecker {
	return &LicenseChecker{
		ProductPermalink: product,
		LicenseKey:       key,
		Variant:          variant,
		exitChan:         make(chan bool),
	}
}

func (checker *LicenseChecker) Start() chan bool {
	go checker.ScheduleCheck()
	return checker.exitChan
}

// ScheduleCheck verifies that the gumroad license is still active every 12 hours
func (checker *LicenseChecker) ScheduleCheck() {
	consecutiveErrors := 0
	timer := time.NewTicker(CheckInterval * time.Hour)

	for range timer.C {
		log.Debug().
			Str("section", "server").
			Str("action", "check_license").
			Int("remaining_retries", MaxConsecutiveErrors-consecutiveErrors).
			Msg("Checking license status")
		_, err := VerifyLicense(checker.ProductPermalink, checker.LicenseKey, checker.Variant, false)
		if err != nil {
			consecutiveErrors++
			log.Info().
				Str("section", "server").
				Str("action", "check_license").
				Int("remaining_retries", MaxConsecutiveErrors-consecutiveErrors).
				Msg(err.Error())
			if consecutiveErrors <= MaxConsecutiveErrors {
				continue
			}
			checker.exitChan <- true
			break
		}
		consecutiveErrors = 0
	}
}
