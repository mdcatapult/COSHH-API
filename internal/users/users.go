package users

import (
	"fmt"
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
		fmt.Println(entry.DN)
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
