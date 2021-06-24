// https://help.gumroad.com/article/76-license-keys
package license

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type LicenseResponse struct {
	Success  bool   `json:"success"`
	Uses     int    `json:"uses"`
	Message  string `json:"message"`
	Purchase struct {
		ID               string `json:"id"`
		ProductName      string `json:"product_name"`
		ProductID        string `json:"product_id"`
		SellerID         string `json:"seller_id"`
		Permalink        string `json:"permalink"`
		ProductPermalink string `json:"product_permalink"`
		CreatedAt        string `json:"created_at"`
		FullName         string `json:"full_name"`
		Variants         string `json:"variants"`
		// purchase was refunded, non-subscription product only
		Refunded bool `json:"refunded"`
		// purchase was refunded, non-subscription product only
		Chargebacked bool `json:"chargebacked"`

		// subscription product only
		// subscription was cancelled,
		SubscriptionCancelledAt *string `json:"subscription_cancelled_at"`
		// we were unable to charge the subscriber's card
		SubscriptionFailedAt *string `json:"subscription_failed_at"`

		CustomFields []string `json:"custom_fields"`
		Email        string   `json:"email"`
	}
}

// productPermalink (the unique permalink of the product)
// licensekey (the license key provided by your customer)
// incrementUsesCount ("true"/"false", optional, default: "true")
func VerifyLicense(productPermalink, licenseKey, variants string, incrementUsesCount bool) (int, error) {
	var (
		response LicenseResponse
	)

	resp, err := http.PostForm("https://api.gumroad.com/v2/licenses/verify",
		url.Values{
			"product_permalink":    {productPermalink},
			"license_key":          {licenseKey},
			"increment_uses_count": {strconv.FormatBool(incrementUsesCount)},
		})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == 404 {
		return 0, errors.New("License Check: " + response.Message)
	}

	if !response.Success {
		return 0, errors.New("License Check: not a successful response")
	}

	if response.Purchase.Variants != variants {
		return response.Uses, errors.New("License Check: subscription does not match " + variants + " != " + response.Purchase.Variants)
	}

	if response.Purchase.Refunded {
		return response.Uses, errors.New("License Check: product was refunded")
	}

	if response.Purchase.Chargebacked {
		return response.Uses, errors.New("License Check: product was chargebacked")
	}

	if response.Purchase.SubscriptionCancelledAt != nil {
		return response.Uses, errors.New("License Check: subscription cancelled")
	}

	if response.Purchase.SubscriptionFailedAt != nil {
		return response.Uses, errors.New("License Check: subscription failed")
	}

	return response.Uses, nil
}
