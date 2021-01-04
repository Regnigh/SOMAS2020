package team5

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

//================================================================
/*	Types */
//=================================================================
const id = shared.Team5

//================ Common Pool =========================================

//cpRequestHistory history of CP Requests
type cpRequestHistory []shared.Resources

//cpAllocationHistory History of allocations
type cpAllocationHistory []shared.Resources

//================ Resource History =========================================

//resourceHistory OUR islands resources per turn
type resourceHistory map[uint]shared.Resources

//================ Foraging History =========================================

// forageOutcome records the ROI on a foraging session
type forageOutcome struct {
	turn   uint
	input  shared.Resources
	output shared.Resources
}

// forageHistory stores history of foraging outcomes
type forageHistory map[shared.ForageType][]forageOutcome

//================ Gifts ===========================================

//giftOutcome defines the gifts
type giftOutcome struct {
	occasions uint
	amount    shared.Resources
}

// type GiftRequest shared.Resources

// type GiftRequestDict map[shared.ClientID]GiftRequest

// giftResponse is a struct of the response and reason
type giftResponse struct {
	AcceptedAmount shared.Resources
	Reason         shared.AcceptReason
}

// giftInfo holds information about the gifts
type giftInfo struct {
	requested shared.GiftRequest
	gifted    shared.GiftOffer
	reason    shared.AcceptReason
}

//giftExchange Looks at the exchanges our island has made (Gifts)
type giftExchange struct {
	IslandRequest map[uint]giftInfo
	OurRequest    map[uint]giftInfo
}

//giftHistory is the history of our gifts
type giftHistory map[shared.ClientID]giftExchange

// Client Information */

type clientConfig struct {

	// Initial non planned foraging
	InitialForageTurns      uint
	MinimumForagePercentage float64
	NormalForagePercentage  float64
	JBForagePercentage      float64

	// Normal foraging
	NormalRandomIncrease float64
	MaxForagePercentage  float64
	SkipForage           uint // Skip for X turns if no positive RoI

	// If resources go above this limit we are balling with money
	JBThreshold shared.Resources
	// Middle class:  Middle < Jeff bezos
	MiddleThreshold shared.Resources
	// Poor: Imperial student < Middle
	ImperialThreshold shared.Resources

	// How much to request when we are dying
	DyingGiftRequestAmount shared.Resources
	// How much to request when we are at Imperial
	ImperialGiftRequestAmount shared.Resources
	// How much to request when we are dying
	MiddleGiftRequestAmount shared.Resources
}

// Client is the island number
type client struct {
	*baseclient.BaseClient

	// Roles
	team5President president

	// History
	resourceHistory     resourceHistory
	forageHistory       forageHistory
	giftHistory         giftHistory
	cpRequestHistory    cpRequestHistory
	cpAllocationHistory cpAllocationHistory

	taxAmount shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources

	config clientConfig
}

//================================================================
/*	Constants */
//=================================================================
// Wealth Tiers
type wealthTier int

const (
	dying wealthTier = iota
	imperialStudent
	middleClass
	jeffBezos
)

const (
	// Accept ...
	Accept shared.AcceptReason = iota
	// DeclineDontNeed ...
	DeclineDontNeed
	// DeclineDontLikeYou ...
	DeclineDontLikeYou
	// Ignored ...
	Ignored
)

//================================================================
/*	Functions */
//=================================================================
//String converts string to number
func (wt wealthTier) String() string {
	strings := [...]string{"dying", "imperialStudent", "middleClass", "jeffBezos"}
	if wt >= 0 && int(wt) < len(strings) {
		return strings[wt]
	}
	return fmt.Sprintf("Unkown wealth state '%v'", int(wt))
}

// GoString implements GoStringer
func (wt wealthTier) GoString() string {
	return wt.String()
}

// MarshalText implements TextMarshaler
func (wt wealthTier) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(wt.String())
}

// MarshalJSON implements RawMessage
func (wt wealthTier) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(wt.String())
}
