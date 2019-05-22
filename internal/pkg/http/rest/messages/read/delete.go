// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package read

import (
	"net/http"

	"github.com/mailchain/mailchain/internal/pkg/stores"
)

// Delete returns a handler
func Delete(store stores.State) func(w http.ResponseWriter, r *http.Request) {
	// Delete swagger:route Delete /messages/{message_id}/read Messages DeleteRead
	//
	// Mark message as unread
	//
	// Responses:
	//   200: StatusOK
	//   404: NotFoundError
	//   422: ValidationError
	return func(w http.ResponseWriter, r *http.Request) {
		doRead(store.DeleteMessageRead, w, r)
	}
}

// DeleteRequest open api documentation
// swagger:parameters DeleteRead
type DeleteRequest struct {
	// Unique id of the message
	//
	// in: path
	// required: true
	MessageID string `json:"message_id"`
}
