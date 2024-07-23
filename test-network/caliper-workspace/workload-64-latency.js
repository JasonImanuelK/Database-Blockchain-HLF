'use strict';

module.exports.info = 'Testing asset operations';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class AssetWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 1;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        // You can add any initialization logic here if required
    }

    async submitTransaction() {
        const assetID = `asset${this.workerIndex}${this.txIndex}`;
        const startTime = Date.now();

        // Buy Asset
        const buyAssetArgs = {
            contractId: 'asset-transfer',
            contractFunction: 'BuyAsset',
            invokerIdentity: 'User1',
            contractArguments: [assetID, '200', 'Ventela', `AssetName${this.txIndex}`, '2000', 'LocationAaaaa', 'In Plant', '2024-05-05'],
            readOnly: false
        };
        this.txIndex++;
        await this.sutAdapter.sendRequests(buyAssetArgs);
        const endTime = Date.now();
        const latency = endTime - startTime;
        console.log(`Transaction ${this.txIndex} latency: ${latency} ms`);
    }

    async cleanupWorkloadModule() {
    }
}

function createWorkloadModule() {
    return new AssetWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
