'use strict';
var helper = require('./helper.js');
var logger = helper.getLogger('Query');

/**
 * Query the chaincode with target function and args
 * @param {*} channelName 
 * @param {*} chaincodeName 
 * @param {*Query function name} fcn 
 * @param {*Query function args} args 
 * @param {*} org_name 
 */
var queryChaincode = function (targets, channelName, chaincodeName, fcn, args) {
	logger.info('\n\n============ Query chaincode on organizations ============\n');
	helper.setupChaincodeDeploy();

	var client = null;
	var channel = null;

	return helper.getClient().then(_client => {
		client = _client;
		return client.getUserContext();
	}).then((user) => {
		if(user !== null){
			logger.info('Current query user is: ' + user.getName());
		}
		channel = client.getChannel(channelName);
		// send query
		var request = {
			targets: targets,
			chaincodeId: chaincodeName,
			fcn: fcn,
			args: args
		};
		return channel.queryByChaincode(request, true);
	}, (err) => {
		throw new Error('Failed to create client ' + err);
	}).then((response_payloads) => {
		if (response_payloads) {
			for (let i = 0; i < response_payloads.length; i++) {
				//check the response is correct or not
			}
			return { success: true, message: response_payloads[0].toString('utf8') };
		} else {
			logger.error('response_payloads is null');
			throw new Error('Failed to get response on query');
		}
	}, (err) => {
		logger.error('Failed to send query due to error: ' + err.stack ? err.stack : err);
		throw new Error('Failed, got error on query');
	});
};


exports.queryChaincode = queryChaincode;
