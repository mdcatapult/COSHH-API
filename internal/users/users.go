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
		"(&(objectClass=organizationalPerson)(!(objectClass=computer)))",
		[]string{"dn", "cn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var users []string
	for _, entry := range sr.Entries {
		if !strings.Contains(entry.DN, "Service Accounts") &&
			!strings.Contains(entry.DN, "Admin") &&
			!strings.Contains(entry.DN, "Leavers") &&
			!strings.Contains(entry.DN, "BindUsers") &&
			!strings.Contains(entry.DN, "Project") &&
			// filter out specific users that are not required
			!strings.Contains(entry.DN, "VPN") &&
			!strings.Contains(entry.DN, "intune") &&
			!strings.Contains(entry.DN, "vbind") &&
			!strings.Contains(entry.DN, "krbtgt") &&
			!strings.Contains(entry.DN, "Services") &&
			!strings.Contains(entry.DN, "SageBI") &&
			!strings.Contains(entry.DN, "Orium 001") &&
			!strings.Contains(entry.DN, "Guest") &&
			!strings.Contains(entry.DN, "AAD_") &&
			!strings.Contains(entry.DN, "MSOL") &&
			!strings.Contains(entry.DN, "MSI Processing") {
			users = append(users, entry.GetAttributeValue("cn"))
		}
	}

	sort.Strings(users)

	return users, nil
}
