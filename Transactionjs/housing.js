var install = require('./app/install-chaincode.js');
var instantiate = require('./app/instantiate-chaincode.js');
var invoke = require('./app/invoke-transaction.js');
var query = require('./app/query.js');
var enroll = require('./app/enroll-user.js');
var Client = require('fabric-client');

//Abstract commandline arguments  
var args = process.argv.splice(2);
// Param declaration
var CHANNEL_NAME = 'default';
var CHAINCODE_ID = 'smartrooves_cc';
//var CHAINCODE_PATH = 'github.com';
var CHAINCODE_VERSION = 'v12';
ng; 

var USER_NAME;
var USER_SECRET;
var base_config_path = "../"
var targets;

//Reload args
param_check(args);

var installChaincodeRequest = {
    chaincodePath: CHAINCODE_PATH,
    chaincodeId: CHAINCODE_ID,
    chaincodeVersion: CHAINCODE_VERSION,
    chaincodeType: CHAINCODE_TYPE
};

var instantiateChaincodeRequest = {
    chanName: CHANNEL_NAME,
    chaincodeId: CHAINCODE_ID,
    chaincodeVersion: CHAINCODE_VERSION,
    fcn: 'init',
    args: ["a"]
};

var targetPart = "ser".concat(getRandomInt(0,100000).toString());
var targetCar = "mer".concat(getRandomInt(0,1000000).toString());

var createVehiclePartRequest = {
    chanName: CHANNEL_NAME,
    chaincodeId: CHAINCODE_ID,
    fcn: 'initVehiclePart',
    args: [targetPart, 'tasa', '1502688979', 'airbag 2020', 'aaimler ag / mercedes', 'false', '0'],
};

var createCarRequest = {
    chanName: CHANNEL_NAME,
    chaincodeId: CHAINCODE_ID, 
    fcn: 'initVehicle',
    args: [targetCar, 'mercedes', 'c class', '1502688979', targetPart, 'mercedes', 'false', '0'],
};

var queryVehiclePartRequest = {
    chanName: CHANNEL_NAME,
    chaincodeId: CHAINCODE_ID,
    fcn: 'queryVehiclePartByOwner',
    args: ['aaimler ag / mercedes'],
};

         


// STEP 1
// Install chaincode
try {
    install.installChaincode(targets, installChaincodeRequest.chaincodeId, installChaincodeRequest.chaincodePath,installChaincodeRequest.chaincodeVersion, installChaincodeRequest.chaincodeType).then((result) => {
                console.log(result)
                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 1/6 : Successfully installed chaincode' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                sleep(2000);
                // STEP 2
                // Instantiate chaincode
                return instantiate.instantiateChaincode(targets, instantiateChaincodeRequest.chanName, instantiateChaincodeRequest.chaincodeId, 
                    instantiateChaincodeRequest.chaincodeVersion,installChaincodeRequest.chaincodeType, instantiateChaincodeRequest.fcn,instantiateChaincodeRequest.args);
        }).then((result2) => {
                console.log(result2)
                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 2/6 : Successfully instantiated chaincode on the channel' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                sleep(2000);

        	    // STEP 3
        	    // Enroll a user
        	    return enroll.enrollUser(USER_NAME,USER_SECRET);
    	}).then((result3) => {
		    if(result3 !== 'Do not need enroll'){
			    console.log(result3)

        		console.log(
            		'\n\n*******************************************************************************' +
            		'\n*******************************************************************************' +
            		'\n*                                          ' +
            		'\n* STEP 3/6 : Successfully enrolled a user' +
            		'\n*                                          ' +
            		'\n*******************************************************************************' +
            		'\n*******************************************************************************\n');

        		sleep(2000);
		    }else{
			    console.log(
            		'\n\n*******************************************************************************' +
            		'\n*******************************************************************************' +
            		'\n*                                          ' +
            		'\n* STEP 3/6 : Enroll a user, but no account is specified, so skipped it' +
            		'\n*                                          ' +
            		'\n*******************************************************************************' +
            		'\n*******************************************************************************\n');
		    }

                // STEP 4
                // invoke chaincode to create a vehicle part
                return invoke.invokeChaincode(targets, createVehiclePartRequest.chanName, createVehiclePartRequest.chaincodeId,
                    createVehiclePartRequest.fcn, createVehiclePartRequest.args);
        }).then((result3) => {
                console.log(result3)

                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 4/6 : Successfully committed vehicle part to ledger' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                sleep(2000);

                // STEP 5
                // invoke chaincode to create a vehicle
                return invoke.invokeChaincode(targets, createCarRequest.chanName, createCarRequest.chaincodeId,
                    createCarRequest.fcn, createCarRequest.args);
        }).then((result4) => {
                console.log(result4)

                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 5/6 : Successfully committed vehicle to ledger' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                sleep(2000);

                // STEP 6
                // invoke chaincode to query vehicle by rich query with query param "Owner"
                return invoke.invokeChaincode(targets, queryVehiclePartRequest.chanName, queryVehiclePartRequest.chaincodeId,
                    queryVehiclePartRequest.fcn, queryVehiclePartRequest.args);
        }).then((result5) => {
                console.log(result5)

                console.log(
                    '\n\n*******************************************************************************' +
                    '\n*******************************************************************************' +
                    '\n*                                          ' +
                    '\n* STEP 6/6 : Successfully queried vehicle part from ledger' +
                    '\n*                                          ' +
                    '\n*******************************************************************************' +
                    '\n*******************************************************************************\n');

                console.log("All Steps Completed Sucessfully");
                process.exit();
        });
} catch (e) {
    console.log(
        '\n\n*******************************************************************************' +
        '\n*******************************************************************************' +
        '\n*                                          ' +
        '\n* Error!!!!!' +
        '\n*                                          ' +
        '\n*******************************************************************************' +
        '\n*******************************************************************************\n');
    console.log(e);
    return;
}


function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function param_check(args) {
    var targetPeerName="";

    if (args.length > 0) {
        var parmUser = args.indexOf('-u');
        var parmChannel = args.indexOf('-c');
        var parmLang = args.indexOf('-l');
        var parmChaincode = args.indexOf('-n');
        var parmPeer = args.indexOf('-p');
        if(parmUser !== -1){
            USER_NAME = args[parmUser + 1];
            USER_SECRET = args[parmUser + 2];
            if(USER_NAME === undefined || USER_SECRET === undefined){
                console.log('Please input the username and password.');
                process.exit();
            }
        }
        if(parmChannel !== -1){
            CHANNEL_NAME = args[parmChannel + 1];
        }
        if(parmLang !== -1){
            CHAINCODE_TYPE = args[parmLang + 1];
        }	
        if(parmChaincode !== -1){
            CHAINCODE_ID = args[parmChaincode +1];
        }     
        if (parmPeer !== -1) {
            targetPeerName = args[parmPeer +1];  
        } 	
    }

    if (targetPeerName.length == 0) {
        var client = Client.loadFromConfig(base_config_path + 'network.yaml');
        var client_org = client.getClientConfig().organization;
        var targets_1 = client.getPeersForOrg(client_org);
        targets =  targets_1.splice(0, 1);
    } else {
        targets = [targetPeerName];
    }
    //console.log("json_targets=" + JSON.stringify(targets));
}

function getRandomInt(min, max) {
	min = Math.ceil(min);
	max = Math.floor(max);
	return Math.floor(Math.random() * (max - min)) + min; //The maximum is exclusive and the minimum is inclusive
}
