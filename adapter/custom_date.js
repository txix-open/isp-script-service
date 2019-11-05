var mdmId = "mdm_id";
var mdmVersion = "mdm_version";

var documentsPath = "documents";
var docTypePath = "id_tp_cd";

var relativesPath = "citizen_relatives";
var relativesTypePath = "rel_tp_cd";

var contactPath = "contacts";
var contactsTypePath = "cont_meth_cat_cd";

var usersPath = "users";
var usersTypePath = "id_type";

var esCredsPath = "escredentials";
var esCredsTypePath = "escred_tp_cd";

var addressesPath = "addresses";
var addressesTypePath = "addr_usage_tp_cd";

var allOtherValues = "*";

var userTypeMapping = {};
userTypeMapping["SSO"] = "SSO";
userTypeMapping["ЕРЛ"] = "ERL";
userTypeMapping["МФЦ"] = "MFC";
userTypeMapping[allOtherValues] = "NOT_SSO";

var relativesTypeMappng = {};
relativesTypeMappng["1000003"] = "1000003";
relativesTypeMappng["1000004"] = "1000004";
relativesTypeMappng[allOtherValues] = "OTHERS";

var id = arg["id"];
var version = arg["version"];
var data = arg["data"];

function expandArray(oldData, typeFieldPath) {
    var newData = {};
    for (var i = 0; i < oldData.length; i++) {
        var oldDataElement = oldData[i];
        var key = oldDataElement[typeFieldPath];
        var newDataElement = newData[key];
        if (newDataElement) {
            newDataElement.push(oldDataElement)
        } else {
            newData[key] = [oldDataElement]
        }
    }
    return newData;
}

function expandArrayWithMapping(oldData, typeFieldPath, mapping) {
    var newData = {};
    for (var i = 0; i < oldData.length; i++) {
        var oldDataElement = oldData[i];
        var key = oldData[i][typeFieldPath];
        if (!mapping[key]) {
            key = mapping[allOtherValues]
        }
        var newDataElement = newData[key];
        if (newDataElement) {
            newDataElement.push(oldDataElement)
        } else {
            newData[key] = [oldDataElement]
        }
    }
    return newData;
}

var response = {};
response[mdmId] = id;
response[mdmVersion] = version;

response[documentsPath] = expandArray(data[documentsPath], docTypePath);
response[contactPath] = expandArray(data[contactPath], contactsTypePath);
response[esCredsPath] = expandArray(data[esCredsPath], esCredsTypePath);
response[addressesPath] = expandArray(data[addressesPath], addressesTypePath);

response[usersPath] = expandArrayWithMapping(data[usersPath], usersTypePath, userTypeMapping);
response[relativesPath] = expandArrayWithMapping(data[relativesPath], relativesTypePath, relativesTypeMappng);

return response;
