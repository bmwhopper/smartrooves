'use strict';
var helper = require('./helper.js');
var logger = helper.getLogger('enroll-user');

/**
 * Enroll a user
 * @param {*Target user's name} name 
 * @param {*Target user's password} password 
 */
var enrollUser = function (name,password) {

	if(name === undefined || password === undefined){
		return "Do not need enroll";
	}

	logger.info('\n\n============ Start enroll a user ============\n');

	var client = null;
	return helper.getClient().then(_client => {
		client = _client;
		return client.setUserContext({username:name, password:password});
	}).then((user) => {
		client.saveUserToStateStore();
		return { status: "Successfully enroll the user: " + name};
	}).catch(err => {
		logger.error('Failed to enroll the user: ' + err.stack ? err.stack : err);
		throw new Error('Failed to enroll the user: ' + err.stack ? err.stack : err);
	});
};

exports.enrollUser = enrollUser;
