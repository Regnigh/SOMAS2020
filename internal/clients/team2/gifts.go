package team2

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type IslandTrust struct {
	island shared.ClientID
	trust  int
}

// Override functions to sort a list
func (p []IslandTrust) Len() int           { return len(p) }
func (p []IslandTrust) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p []IslandTrust) Less(i, j int) bool { return p[i].trust < p[j].trust }

// GetGiftRequests allows clients to signalize that they want a gift
// This information is fed to OfferGifts of all other clients.
func (c *client) GetGiftRequests() shared.GiftRequestDict {
	turn := c.gameState().Turn
	requests := shared.GiftRequestDict{}

	// check our critical and threshold - if either is off - request
	ourAgentCritical := status == shared.Critical
	requestAmount := internalThreshold(c) - c.gameState().ClientInfo.Resources

	// confidence[island] * requestAmount until -> target
	if ourAgentCritical || requestAmount > 0 {
		target := 1.5 * requestAmount
		var trustRank []IslandTrust

		for team, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {

			// get our confidence in the team
			if status != shared.Critical {
				islandConf := IslandTrust{
					island: team,
					trust:  c.confidence("GiftWeRequest", team),
				}
				trustRank = append(confidences, islandConf)
			}

		}

		// keep a ranked list of the teams
		sort.Sort(trustRank)

		// request confidence (0-1)*amountNeeded to consecutive islands in rank
		// until amountRequested = (factor eg 1.5) * amountNeeded (to accommodate for some islands not giving us a gift)
		for i := 0; i < len(trustRank); i++ {
			requestAmount := trustRank[i].trust * target
			requestedTo := trustRank[i].island
			requests[requestedTo] = shared.GiftRequest(requestAmount)

			// to keep track in our history
			newGiftRequest := GiftInfo{
				requested: requestAmount,
			}
			c.giftHist[requestedTo].OurRequest[turn] = newGiftRequest

			target -= requestAmount

			if target <= 0 {
				return requests
			}
		}
	}

	return requests
}

// GetGiftOffers allows clients to make offers in response to gift requests by other clients.
// It can offer multiple partial gifts.
// COMPULSORY, you need to implement this method. This placeholder implementation offers no gifts,
// unless another team is critical.
func (c *client) GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict {
	offers := shared.GiftOfferDict{}
	turn := c.gameState().Turn

	// if we are critical do not offer gifts-> there should be a way to see which other islands are critical
	// if we are not critical and another island is critical offer gift
	// do not offer more than proportion of total resources we have
	ourAgentCritical := status == shared.Critical

	// prioritize giving gifts to islands we trust (for now confidence)

	// Give no more than half of amount before we reach threshold
	maxToGive := (c.gameState().ClientInfo.Resources - internalThreshold(c)) / 2

	var trustRank []IslandTrust
	if ourAgentCritical || maxToGive <= 0 {
		for team, reqAmount := range receivedRequests {
			// max would be 200
			// c.confidence("ReceivedRequests", team) should reflect the status of an island and int,float64 requests hist
			// c.confidence("GiftWeRequest", team) should reflect how they respond to our requests
			confidenceMetric := c.confidence("ReceivedRequests", team) + c.confidence("GiftWeRequest", team)
			if confidenceMetric > 100 || shared.Critical {
				islandConf := IslandTrust{
					island: team,
					trust:  confidenceMetric,
				}
				trustRank = append(confidences, islandConf)
			}
		}

		sort.Sort(trustRank)

		for i := 0; i < len(trustRank); i++ {
			offeredTo := trustRank[i].island
			offeredAmount := receivedRequests[offeredTo]

			if offeredAmount >= maxToGive {
				offeredAmount = maxToGive
				maxToGive = 0
			}

			offers[offeredTo] = shared.GiftOffer(offeredAmount)

			// to keep track in our history
			newGiftRequest := GiftInfo{
				requested: receivedRequests[offeredTo],
				gifted:    shared.GiftOffer(offeredAmount),
			}
			c.giftHist[offeredTo].IslandRequest[turn] = newGiftRequest

			maxToGive -= offeredAmount

			if maxToGive <= 0 {
				return offers
			}
		}
	}

	return offers
}

// GetGiftResponses allows clients to accept gifts offered by other clients.
// It also needs to provide a reasoning should it not accept the full amount.
// COMPULSORY, you need to implement this method

func (c *client) GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict {
	responses := shared.GiftResponseDict{}
	turn := c.gameState().Turn

	for client, offer := range receivedOffers {
		responses[client] = shared.GiftResponse{
			AcceptedAmount: shared.Resources(offer),
			Reason:         shared.Accept,
		}
	}

	for island, response := range responses {
		newGiftRequest := GiftInfo{
			requested: c.giftHist[island].OurRequest[turn].requested,
			gifted:    shared.GiftOffer(response.AcceptedAmount),
			reason:    shared.AcceptReason(response.AcceptedAmount),
		}
		c.giftHist[island].OurRequest[turn] = newGiftRequest
	}
	return responses
}

// UpdateGiftInfo will be called by the server with all the responses you received
// from the gift session. It is up to you to complete the transactions with the methods
// that you will implement yourself below. This allows for opinion formation.
// COMPULSORY, you need to implement this method
func (c *client) UpdateGiftInfo(receivedResponses shared.GiftResponseDict) {
	turn := c.gameState().Turn

	for island, response := range receivedResponses {
		newGiftRequest := GiftInfo{
			requested: c.giftHist[island].IslandRequest[turn].requested,
			gifted:    shared.GiftOffer(response.AcceptedAmount),
			reason:    shared.AcceptReason(response.AcceptedAmount),
		}
		c.giftHist[island].IslandRequest[turn] = newGiftRequest
	}
}

// SentGift is executed at the end of each turn and notifies clients that
// their gift was successfully sent, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) SentGift(sent shared.Resources, to shared.ClientID) {
	// You can check your updated resources like this:
	// myResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

}

// ReceivedGift is executed at the end of each turn and notifies clients that
// their gift was successfully received, along with the offer details.
// COMPULSORY, you need to implement this method
func (c *client) ReceivedGift(received shared.Resources, from shared.ClientID) {
	// You can check your updated resources like this:
	// myResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

}
