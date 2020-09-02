package adapter

import (
	"github.com/jmoiron/jsonq"
)

const (
	mdmId        = "mdm_id"
	mdmVersion   = "mdm_version"
	mdmDocuments = "documents"
	mdmRelatives = "citizen_relatives"
	mdmContacts  = "contacts"
	mdmUsers     = "users"
	mdmEsCreds   = "escredentials"
	mdmAddresses = "addresses"

	documentsPath = "documents"
	docTypePath   = "id_tp_cd"

	relativesPath     = "citizen_relatives"
	relativesTypePath = "rel_tp_cd"

	contactPath      = "contacts"
	contactsTypePath = "cont_meth_cat_cd"

	usersPath     = "users"
	usersTypePath = "id_type"

	esCredsPath     = "escredentials"
	esCredsTypePath = "escred_tp_cd"

	addressesPath     = "addresses"
	addressesTypePath = "addr_usage_tp_cd"

	allOtherValues = "*"
)

var (
	usersTypeMapping = map[string]string{
		"SSO":          "SSO",
		"ЕРЛ":          "ERL",
		"МФЦ":          "MFC",
		allOtherValues: "NOT_SSO",
	}
	relativesTypeMapping = map[string]string{
		"1000003":      "1000003",
		"1000004":      "1000004",
		allOtherValues: "OTHERS",
	}
)

func GetCustomData(data map[string]interface{}, id string, version int64) map[string]interface{} {
	jq := jsonq.NewQuery(data)

	m := make(map[string]interface{}, 6)
	m[mdmId] = id
	m[mdmVersion] = version

	expandArray(jq, m, documentsPath, docTypePath, mdmDocuments)
	expandArray(jq, m, contactPath, contactsTypePath, mdmContacts)
	expandArray(jq, m, esCredsPath, esCredsTypePath, mdmEsCreds)
	expandArray(jq, m, addressesPath, addressesTypePath, mdmAddresses)

	// Custom transform fields
	expandArrayWithMapping(jq, m, usersPath, usersTypePath, mdmUsers, usersTypeMapping)
	expandArrayWithMapping(jq, m, relativesPath, relativesTypePath, mdmRelatives, relativesTypeMapping)

	return m
}

func expandArray(src *jsonq.JsonQuery, dst map[string]interface{}, arrayPath, typeFieldPath, dstKey string) {
	docs, err := src.ArrayOfObjects(arrayPath)
	l := len(docs)
	if err == nil && l > 0 {
		documents := make(map[string][]interface{}, l)
		for _, d := range docs {
			if docType, ok := d[typeFieldPath]; ok {
				if docTypeStr, ok := docType.(string); ok {
					putDoc(documents, docTypeStr, d)
				}
			}
		}
		if len(documents) > 0 {
			dst[dstKey] = documents
		}
	}
}

func expandArrayWithMapping(
	src *jsonq.JsonQuery,
	dst map[string]interface{},
	arrayPath, typeFieldPath, dstKey string,
	mapping map[string]string,
) {
	var documents map[string][]interface{}
	docs, err := src.ArrayOfObjects(arrayPath)
	l := len(docs)
	if err == nil && l > 0 {
		documents = make(map[string][]interface{}, l)
		for _, d := range docs {
			if docType, ok := d[typeFieldPath]; !ok {
				continue
			} else if docTypeStr, ok := docType.(string); !ok {
				continue
			} else if v, ok := mapping[docTypeStr]; ok {
				putDoc(documents, v, d)
			} else if v, ok := mapping[allOtherValues]; ok {
				putDoc(documents, v, d)
			}
		}
	}
	if len(documents) > 0 {
		dst[dstKey] = documents
	}
}

func putDoc(dst map[string][]interface{}, key string, value interface{}) {
	if arr, ok := dst[key]; ok {
		dst[key] = append(arr, value)
	} else {
		dst[key] = []interface{}{value}
	}
}
