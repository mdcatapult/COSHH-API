/*
* Copyright $today.year Medicines Discovery Catapult
* Licensed under the Apache License, Version 2.0 (the "Licence");
* you may not use this file except in compliance with the Licence.
* You may obtain a copy of the Licence at
*     http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the Licence is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the Licence for the specific language governing permissions and
* limitations under the Licence.
 */

package users

import (
	"github.com/go-ldap/ldap/v3"
	"log"
	"sort"
	"strings"
)

func GetUsers(ldapUser, ldapPassword string) ([]string, error) {

	// connect to ldap server
	ldapURL := "ldap://medcat.local"
	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// binding to ldap server
	err = l.Bind(ldapUser, ldapPassword)
	if err != nil {
		log.Fatal(err)
	}

	searchRequest := ldap.NewSearchRequest(
		"dc=medcat,dc=local",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectCategory=person))",
		[]string{},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var userNames []string
	for _, entry := range sr.Entries {
		if !strings.Contains(entry.DN, "Service Accounts") &&
			!strings.Contains(entry.DN, "Admin") &&
			!strings.Contains(entry.DN, "Leavers") &&
			!strings.Contains(entry.DN, "BindUsers") {
			userNames = append(userNames, entry.GetAttributeValue("cn"))
		}
	}

	sort.Strings(userNames)

	return userNames, nil
}
