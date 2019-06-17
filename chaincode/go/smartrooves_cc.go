package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "strconv"
  "strings"
  "time"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)

// SmartRoovesChaincode example simple Chaincode implementation
type SmartRoovesChaincode struct {
}

type tenant struct {
  ObjectType         string `json:"docType"`      //docType is used to distinguish the various types of objects in state database
  PPSNumber          string `json:"ppsNumber"`    //the fieldtags are needed to keep case from bouncing around
  FirstName          string `json:"firstName"`
  LastName           string `json:"lastName"`
  DateOfBirth        string `json:"dateOfBirth"`
  MaritalStatus      string `json:"maritalStatus"`
  FamilyComponents   int    `json:"familyComponents"`
  AnnualIncome       int    `json:"annualIncome"`
  MaxRentSpend       int    `json:"maxRentSpend"`
  Disability         bool   `json:"disability"`
  ApartmentId        string `json:"apartmentId"` // FK to apartment id
}

type apartment struct {
  ObjectType         string `json:"docType"`       //docType is used to distinguish the various types of objects in state database
  ApartmentId        string `json:"apartmentId"`   //the fieldtags are needed to keep case from bouncing around
  Capacity           int    `json:"capacity"`
  Address            string `json:"address"`
  AnnualRent         int    `json:"annualRent"`
  GovAnnualRentAllow int    `json:"govAnnualRentAllow"`
  Assigned           bool   `json:"assigned"`
  RentExtended       bool   `json:"rentExtended"`
  Accessibility      bool   `json:"accessibility"`
  StartLeaseDate     string `json:"startLeaseDate"`
  Owner              string `json:"owner"`
}


// ===================================================================================
// Main
// ===================================================================================
func main() {
  err := shim.Start(new(SmartRoovesChaincode))
  if err != nil {
    fmt.Printf("Error starting Parts Trace chaincode: %s", err)
  }
}

// ===================================================================================
//  init...initializes chaincode
// ===================================================================================
func (t *SmartRoovesChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
  return shim.Success(nil)
}

// ===================================================================================
// Invoke - Our entry point for Invocations
// ===================================================================================
func (t *SmartRoovesChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
  function, args := stub.GetFunctionAndParameters()
  fmt.Println("invoke is running " + function)

  // Handle different functions
  if function == "initApartment" { //create a new apartment
    return t.initApartment(stub, args)
  } else if function == "initTenant" { //create a new tenant
    return t.initTenant(stub, args)
  } else if function == "transferApartmentToGov" { //transfer ownership of an apartment to the government
    return t.transferApartmentToGov(stub, args)
  } else if function == "assignApartmentToTenant" { //assign an apartment to a tenant
    return t.assignApartmentToTenant(stub, args)
  } else if function == "recallApartmentFromTenant" { //recall an apartment from a tenant
    return t.recallApartmentFromTenant(stub, args)
  } else if function == "getAvailableApartments" { //return all available apartments
    return t.getAvailableApartments(stub, args)
  } else if function == "getAvailableTenants" { //return all the tenants that don't have an apartment assigned
    return t.getAvailableTenants(stub, args)
  } else if function == "querySmartRooves" { //return results based on input query string
    return t.querySmartRooves(stub, args)
  } else if function == "getHistoryForRecord" { //get history of values for a record
    return t.getHistoryForRecord(stub, args)
  }

  fmt.Println("invoke did not find func: " + function) //error
  return shim.Error("Received unknown function invocation")
}

// ============================================================
// initApartment - create a new apartment, store into chaincode state
// ============================================================
func (t *SmartRoovesChaincode) initApartment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error

  if len(args) != 10 {
    return shim.Error("Incorrect number of arguments. Expecting 10")
  }

  // ==== Input sanitation ====
  fmt.Println("- start init apartment")
  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string")
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string")
  }
  if len(args[3]) <= 0 {
    return shim.Error("4th argument must be a non-empty string")
  }
  if len(args[4]) <= 0 {
    return shim.Error("5th argument must be a non-empty string")
  }
  if len(args[5]) <= 0 {
    return shim.Error("6th argument must be a non-empty string")
  }
  if len(args[6]) <= 0 {
    return shim.Error("7th argument must be a non-empty string")
  }
  if len(args[7]) <= 0 {
    return shim.Error("8th argument must be a non-empty string")
  }
  if len(args[8]) <= 0 {
    return shim.Error("9th argument must be a non-empty string")
  }
  if len(args[9]) <= 0 {
    return shim.Error("10th argument must be a non-empty string")
  }

  apartmentId := strings.ToLower(args[0])
  capacity, err := strconv.Atoi(args[1])
  if err != nil {
    return shim.Error("2nd argument must be a numeric string")
  }
  address := args[2]
  annualRent, err := strconv.Atoi(args[3])
  if err != nil {
    return shim.Error("4th argument must be a numeric string")
  }
  govAnnualRentAllow, err := strconv.Atoi(args[4])
  if err != nil {
    return shim.Error("5th argument must be a numeric string")
  }
  assigned, err := strconv.ParseBool(args[5])
  if err != nil {
    return shim.Error("6th argument must be a boolean string")
  }
  rentExtended, err := strconv.ParseBool(args[6])
  if err != nil {
    return shim.Error("7th argument must be a boolean string")
  }
  accessibility, err := strconv.ParseBool(args[7])
  if err != nil {
    return shim.Error("8th argument must be a boolean string")
  }
  startLeaseDate := strings.ToLower(args[8])
  owner := strings.ToLower(args[9])
  
  if annualRent < govAnnualRentAllow {
    return shim.Error("Annual rent must be greater or equal than Goverment allowance")
  }

  // ==== Check if apartment already exists ====
  apartmentAsBytes, err := stub.GetState(apartmentId)
  if err != nil {
    return shim.Error("Failed to get apartment: " + err.Error())
  } else if apartmentAsBytes != nil {
    fmt.Println("This apartment already exists: " + apartmentId)
    return shim.Error("This apartment already exists: " + apartmentId)
  }

  // ==== Create apartment object and marshal to JSON ====
  objectType := "apartment"
  apartment := &apartment{objectType, apartmentId, capacity, address, annualRent, govAnnualRentAllow, assigned, rentExtended, accessibility, startLeaseDate, owner}

  apartmentJSONasBytes, err := json.Marshal(apartment)
  if err != nil {
    return shim.Error(err.Error())
  }

  // === Save apartment to state ===
  err = stub.PutState(apartmentId, apartmentJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Apartment saved and indexed. Return success ====
  fmt.Println("- end init apartment")
  return shim.Success(nil)
}

// ============================================================
// initTenant - create a new apartment, store into chaincode state
// ============================================================
func (t *SmartRoovesChaincode) initTenant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error

  if len(args) != 10 {
    return shim.Error("Incorrect number of arguments. Expecting 10")
  }

  // ==== Input sanitation ====
  fmt.Println("- start init apartment")
  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string")
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string")
  }
  if len(args[3]) <= 0 {
    return shim.Error("4th argument must be a non-empty string")
  }
  if len(args[4]) <= 0 {
    return shim.Error("5th argument must be a non-empty string")
  }
  if len(args[5]) <= 0 {
    return shim.Error("6th argument must be a non-empty string")
  }
  if len(args[6]) <= 0 {
    return shim.Error("7th argument must be a non-empty string")
  }
  if len(args[7]) <= 0 {
    return shim.Error("8th argument must be a non-empty string")
  }
  if len(args[8]) <= 0 {
    return shim.Error("9th argument must be a non-empty string")
  }
  if len(args[9]) <= 0 {
    return shim.Error("10th argument must be a non-empty string")
  }

  ppsNumber := strings.ToLower(args[0])
  firstName := args[1]
  lastName := args[2]
  dateOfBirth := strings.ToLower(args[3])
  maritalStatus := strings.ToLower(args[4])
  familyComponents, err := strconv.Atoi(args[5])
  if err != nil {
    return shim.Error("6th argument must be a numeric string")
  }
  annualIncome, err := strconv.Atoi(args[6])
  if err != nil {
    return shim.Error("7th argument must be a numeric string")
  }
  maxRentSpend, err := strconv.Atoi(args[7])
  if err != nil {
    return shim.Error("8th argument must be a numeric string")
  }
  disability, err := strconv.ParseBool(args[8])
  if err != nil {
    return shim.Error("9th argument must be a boolean string")
  }
  apartmentId := strings.ToLower(args[9])
  
  if annualIncome < maxRentSpend {
    return shim.Error("Annual income must be greater or equal than maximum rent spend")
  }

  // ==== Check if tenant already exists ====
  tenantAsBytes, err := stub.GetState(ppsNumber)
  if err != nil {
    return shim.Error("Failed to get tenant: " + err.Error())
  } else if tenantAsBytes != nil {
    fmt.Println("This tenant already exists: " + ppsNumber)
    return shim.Error("This tenant already exists: " + ppsNumber)
  }

  // ==== Create tenant object and marshal to JSON ====
  objectType := "tenant"
  tenant := &tenant{objectType, ppsNumber, firstName, lastName, dateOfBirth, maritalStatus, familyComponents, annualIncome, maxRentSpend, disability, apartmentId}

  tenantJSONasBytes, err := json.Marshal(tenant)
  if err != nil {
    return shim.Error(err.Error())
  }

  // === Save tenant to state ===
  err = stub.PutState(ppsNumber, tenantJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Tenant saved and indexed. Return success ====
  fmt.Println("- end init tenant")
  return shim.Success(nil)
}

// ============================================================
// transferApartmentToGov  - transfer apartment to government
// ============================================================
func (t *SmartRoovesChaincode) transferApartmentToGov(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error

  // 0
  // "apartmentId"
  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  // ==== Input sanitation ====
  fmt.Println("- start transfer apartment to government")
  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }

  apartmentId := strings.ToLower(args[0])

  // ==== Check if apartment already exists ====
  apartmentAsBytes, err := stub.GetState(apartmentId)
  if err != nil {
    return shim.Error("Failed to get apartment: " + err.Error())
  } else if apartmentAsBytes == nil {
    fmt.Println("This apartment does not exist: " + apartmentId)
    return shim.Error("This apartment does not exist: " + apartmentId)
  }

  apartment := apartment{}
  err = json.Unmarshal(apartmentAsBytes, &apartment) //unmarshal it aka JSON.parse()
  if err != nil {
    return shim.Error(err.Error())
  }

  if apartment.Assigned {
    return shim.Error("This apartment is alread assigned to a tenant: " + apartmentId)
  }
  
  if apartment.Owner == "gov" {
    return shim.Error("This apartment is alread assigned to gov: " + apartmentId)
  }

  apartment.Owner = "gov"

  apartmentJSONasBytes, err := json.Marshal(apartment)
  if err != nil {
    return shim.Error(err.Error())
  }

  // === Save apartment to state ===
  err = stub.PutState(apartmentId, apartmentJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Tenant saved and indexed. Return success ====
  fmt.Println("- end transfer apartment to government")
  return shim.Success(nil)
}

// ============================================================
// assignApartmentToTenant  - assign an apartment to a tenant
// ============================================================
func (t *SmartRoovesChaincode) assignApartmentToTenant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error

  // 0              1            2
  // "apartmentId", "ppsNumber", "startLeaseDate"
  if len(args) != 3 {
    return shim.Error("Incorrect number of arguments. Expecting 3")
  }

  // ==== Input sanitation ====
  fmt.Println("- start assign apartment to tenant")
  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string")
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string")
  }

  apartmentId := strings.ToLower(args[0])
  ppsNumber := strings.ToLower(args[1])
  startLeaseDate := strings.ToLower(args[2])

  // ==== Check if apartment already exists ====
  apartmentAsBytes, err := stub.GetState(apartmentId)
  if err != nil {
    return shim.Error("Failed to get apartment: " + err.Error())
  } else if apartmentAsBytes == nil {
    fmt.Println("This apartment does not exist: " + apartmentId)
    return shim.Error("This apartment does not exist: " + apartmentId)
  }

  // ==== Check if tenant already exists ====
  tenantAsBytes, err := stub.GetState(ppsNumber)
  if err != nil {
    return shim.Error("Failed to get tenant: " + err.Error())
  } else if tenantAsBytes == nil {
    fmt.Println("This tenant does not exist: " + ppsNumber)
    return shim.Error("This tenant does not exist: " + ppsNumber)
  }

  apartment := apartment{}
  err = json.Unmarshal(apartmentAsBytes, &apartment) //unmarshal it aka JSON.parse()
  if err != nil {
    return shim.Error(err.Error())
  }

  tenant := tenant{}
  err = json.Unmarshal(tenantAsBytes, &tenant) //unmarshal it aka JSON.parse()
  if err != nil {
    return shim.Error(err.Error())
  }

  if apartment.Assigned {
    return shim.Error("This apartment is alread assigned to a tenant: " + apartmentId)
  }
  if tenant.ApartmentId != "null" {
    return shim.Error("This tenant is already assigned to an apartment: " + ppsNumber)
  }
  if apartment.Capacity < tenant.FamilyComponents {
    return shim.Error("This apartment does not have enough capacity for this tenant: " + ppsNumber)
  }
  if apartment.AnnualRent - apartment.GovAnnualRentAllow > tenant.MaxRentSpend {
    return shim.Error("This apartment is too expensive for this tenant: " + ppsNumber)
  }

  apartment.Assigned = true
  apartment.StartLeaseDate = startLeaseDate
  apartment.RentExtended = false
  tenant.ApartmentId = apartmentId

  apartmentJSONasBytes, err := json.Marshal(apartment)
  if err != nil {
    return shim.Error(err.Error())
  }

  // === Save apartment to state ===
  err = stub.PutState(apartmentId, apartmentJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  tenantJSONasBytes, err := json.Marshal(tenant)
  if err != nil {
    return shim.Error(err.Error())
  }

  // === Save tenant to state ===
  err = stub.PutState(ppsNumber, tenantJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Tenant saved and indexed. Return success ====
  fmt.Println("- end assign apartment to tenant")
  return shim.Success(nil)
}

// ============================================================
// recallApartmentFromTenant - recall an apartment from a tenant
// ============================================================
func (t *SmartRoovesChaincode) recallApartmentFromTenant(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error

  // 0              1
  // "apartmentId", "ppsNumber"
  if len(args) != 2 {
    return shim.Error("Incorrect number of arguments. Expecting 2")
  }

  // ==== Input sanitation ====
  fmt.Println("- start recall apartment from tenant")
  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string")
  }

  apartmentId := strings.ToLower(args[0])
  ppsNumber := strings.ToLower(args[1])

  // ==== Check if apartment already exists ====
  apartmentAsBytes, err := stub.GetState(apartmentId)
  if err != nil {
    return shim.Error("Failed to get apartment: " + err.Error())
  } else if apartmentAsBytes == nil {
    fmt.Println("This apartment does not exist: " + apartmentId)
    return shim.Error("This apartment does not exist: " + apartmentId)
  }

  // ==== Check if tenant already exists ====
  tenantAsBytes, err := stub.GetState(ppsNumber)
  if err != nil {
    return shim.Error("Failed to get tenant: " + err.Error())
  } else if tenantAsBytes == nil {
    fmt.Println("This tenant does not exist: " + ppsNumber)
    return shim.Error("This tenant does not exist: " + ppsNumber)
  }

  apartment := apartment{}
  err = json.Unmarshal(apartmentAsBytes, &apartment) //unmarshal it aka JSON.parse()
  if err != nil {
    return shim.Error(err.Error())
  }

  tenant := tenant{}
  err = json.Unmarshal(tenantAsBytes, &tenant) //unmarshal it aka JSON.parse()
  if err != nil {
    return shim.Error(err.Error())
  }

  if !apartment.Assigned {
    return shim.Error("This apartment is not assigned to a tenant: " + apartmentId)
  }
  if tenant.ApartmentId != apartmentId {
    return shim.Error("This tenant is not assigned to this apartment: " + ppsNumber)
  }

  apartment.Assigned = false
  apartment.StartLeaseDate = "null"
  apartment.RentExtended = false
  tenant.ApartmentId = "null"

  apartmentJSONasBytes, err := json.Marshal(apartment)
  if err != nil {
    return shim.Error(err.Error())
  }

  // === Save apartment to state ===
  err = stub.PutState(apartmentId, apartmentJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  tenantJSONasBytes, err := json.Marshal(tenant)
  if err != nil {
    return shim.Error(err.Error())
  }

  // === Save tenant to state ===
  err = stub.PutState(ppsNumber, tenantJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Tenant saved and indexed. Return success ====
  fmt.Println("- end recall apartment from tenant")
  return shim.Success(nil)
}

// =========================================================================================
// getAvailableApartments return all the apartments that are not assigned (assigned == false)
// =========================================================================================
func (t *SmartRoovesChaincode) getAvailableApartments(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  queryString := fmt.Sprintf("SELECT valueJson FROM <STATE> WHERE json_extract(valueJson, '$.docType', '$.assigned') = '[\"apartment\",%s]'", "false")

  queryResults, err := getQueryResultForQueryString(stub, queryString)
  if err != nil {
    return shim.Error(err.Error())
  }
  return shim.Success(queryResults)
}

// =========================================================================================
// getAvailableTenants return all the tenants that don't have an apartment assigned (apartmentId == "null")
// =========================================================================================
func (t *SmartRoovesChaincode) getAvailableTenants(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  queryString := fmt.Sprintf("SELECT valueJson FROM <STATE> WHERE json_extract(valueJson, '$.docType', '$.apartmentId') = '[\"tenant\",\"%s\"]'", "null")

  queryResults, err := getQueryResultForQueryString(stub, queryString)
  if err != nil {
    return shim.Error(err.Error())
  }
  return shim.Success(queryResults)
}

// =========================================================================================
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// =========================================================================================
func (t *SmartRoovesChaincode) querySmartRooves(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  // "queryString"
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  queryString := args[0]

  queryResults, err := getQueryResultForQueryString(stub, queryString)
  if err != nil {
    return shim.Error(err.Error())
  }
  return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

  fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

  resultsIterator, err := stub.GetQueryResult(queryString)
  if err != nil {
    return nil, err
  }
  defer resultsIterator.Close()

  // buffer is a JSON array containing QueryRecords
  var buffer bytes.Buffer
  buffer.WriteString("[")

  bArrayMemberAlreadyWritten := false
  for resultsIterator.HasNext() {
    queryResponse, err := resultsIterator.Next()
    if err != nil {
      return nil, err
    }
    // Add a comma before array members, suppress it for the first array member
    if bArrayMemberAlreadyWritten == true {
      buffer.WriteString(",")
    }
    buffer.WriteString(string(queryResponse.Value))
    bArrayMemberAlreadyWritten = true
  }
  buffer.WriteString("]")

  fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

  return buffer.Bytes(), nil
}

// ===========================================================================================
// getHistoryForRecord returns the histotical state transitions for a given key of a record
// ===========================================================================================
func (t *SmartRoovesChaincode) getHistoryForRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  recordKey := args[0]

  fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

  resultsIterator, err := stub.GetHistoryForKey(recordKey)
  if err != nil {
    return shim.Error(err.Error())
  }
  defer resultsIterator.Close()

  // buffer is a JSON array containing historic values for the key/value pair
  var buffer bytes.Buffer
  buffer.WriteString("[")

  bArrayMemberAlreadyWritten := false
  for resultsIterator.HasNext() {
    response, err := resultsIterator.Next()
    if err != nil {
      return shim.Error(err.Error())
    }
    // Add a comma before array members, suppress it for the first array member
    if bArrayMemberAlreadyWritten == true {
      buffer.WriteString(",")
    }
    buffer.WriteString("{\"TxId\":")
    buffer.WriteString("\"")
    buffer.WriteString(response.TxId)
    buffer.WriteString("\"")

    buffer.WriteString(", \"Value\":")
    // if it was a delete operation on given key, then we need to set the
    //corresponding value null. Else, we will write the response.Value
    //as-is (as the Value itself a JSON vehiclePart)
    if response.IsDelete {
      buffer.WriteString("null")
    } else {
      buffer.WriteString(string(response.Value))
    }

    buffer.WriteString(", \"Timestamp\":")
    buffer.WriteString("\"")
    buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
    buffer.WriteString("\"")

    buffer.WriteString(", \"IsDelete\":")
    buffer.WriteString("\"")
    buffer.WriteString(strconv.FormatBool(response.IsDelete))
    buffer.WriteString("\"")

    buffer.WriteString("}")
    bArrayMemberAlreadyWritten = true
  }
  buffer.WriteString("]")

  fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

  return shim.Success(buffer.Bytes())
}
