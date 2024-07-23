'use strict';

module.exports.info = 'Testing asset operations';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class AssetWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 2001;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        // You can add any initialization logic here if required
    }

    async submitTransaction() {
        this.txIndex++;
        const assetID = `asset${this.workerIndex}${this.txIndex}`;

        // Buy Asset
        const buyAssetArgs = {
            contractId: 'asset-transfer',
            contractFunction: 'BuyAsset',
            invokerIdentity: 'User1',
            contractArguments: [assetID, '200', 'Ventela', `AssetName${this.txIndex}`, '2000', 'LocationAaaaa', 'In Plant', '2024-05-05'],
            readOnly: false
        };
        await this.sutAdapter.sendRequests(buyAssetArgs);
    }

    async cleanupWorkloadModule() {
    }
}

function createWorkloadModule() {
    return new AssetWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
