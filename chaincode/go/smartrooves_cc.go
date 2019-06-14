package main

import (
  //"bytes"
  "encoding/json"
  "fmt"
  "strconv"
  "strings"
  //"time"

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
  if function == "initApartment" { //create a new initApartment
    return t.initApartment(stub, args)
  // } else if function == "transferVehiclePart" { //change owner of a specific vehicle part
    // return t.transferVehiclePart(stub, args)
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