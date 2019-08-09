package engine

/*
 * Copyright © 2006-2019 Around25 SRL <office@around25.com>
 *
 * Licensed under the Around25 Exchange License Agreement (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.around25.com/licenses/EXCHANGE_LICENSE
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Cosmin Harangus <cosmin@around25.com>
 * @copyright 2006-2019 Around25 SRL <office@around25.com>
 * @license 	EXCHANGE_LICENSE
 */

import (
	"gitlab.com/around25/products/matching-engine/model"
)

func (s *SkipList) addOrder(price uint64, order model.Order) {
	if pricePoint, ok := s.Get(price); ok {
		pricePoint.Entries = append(pricePoint.Entries, order)
		return
	}
	s.Set(price, &PricePoint{
		Entries: []model.Order{order},
	})
}

// Remove an order from the list of price points
// The method will also remove the price point entry if both book entry lists are empty
func (s *SkipList) removeEntryByPriceAndIndex(price uint64, pricePoint *PricePoint, index int) {
	pricePoint.Entries = append(pricePoint.Entries[:index], pricePoint.Entries[index+1:]...)
	if len(pricePoint.Entries) == 0 {
		s.Delete(price)
	}
}
