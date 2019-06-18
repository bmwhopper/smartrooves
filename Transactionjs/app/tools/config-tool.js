'use strict'

var fsx = require('fs-extra');
var Client = require('fabric-client');
var base_config_path = "../"

var ConfigTool = class {

    constructor() {
        this._client = null;
    }

    cleanUpConfigCache(orgName) {
        let client = Client.loadFromConfig(base_config_path + orgName + '.yaml');
        let client_config = client.getClientConfig();

        let store_path = client_config.credentialStore.path;
        fsx.removeSync(store_path);

        let crypto_path = client_config.credentialStore.cryptoStore.path;
        fsx.removeSync(crypto_path);
    }

    initClient() {
        if (this._client != null) {
            return Promise.resolve(this._client);
        }
        var client = Client.loadFromConfig(base_config_path + 'network.yaml');
        return client.initCredentialStores().then((nothing) => {
            this._client = client;
            return Promise.resolve(client);
        });
    }
}

module.exports = ConfigTool;
