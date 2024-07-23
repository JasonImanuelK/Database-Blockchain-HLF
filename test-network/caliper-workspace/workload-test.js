'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        // You can add any initialization logic here if required
    }

    async submitTransaction() {
        this.txIndex++;
        let txType = this.txIndex % 3;
        let myArgs;

        switch (txType) {
            case 0:
                // Create transaction
                myArgs = {
                    contractId: 'mycontract',
                    contractFunction: 'createAsset',
                    invokerIdentity: 'User1',
                    contractArguments: [`asset${this.txIndex}`, 'blue', '20', 'Tom', '1300'],
                    readOnly: false
                };
                break;
            case 1:
                // Query transaction
                myArgs = {
                    contractId: 'mycontract',
                    contractFunction: 'readAsset',
                    invokerIdentity: 'User1',
                    contractArguments: [`asset${this.txIndex % 100}`],  // Assuming you have at least 100 assets
                    readOnly: true
                };
                break;
            case 2:
                // Update transaction
                myArgs = {
                    contractId: 'mycontract',
                    contractFunction: 'updateAsset',
                    invokerIdentity: 'User1',
                    contractArguments: [`asset${this.txIndex % 100}`, 'red', '30', 'Jerry', '1500'],
                    readOnly: false
                };
                break;
            default:
                throw new Error('Unexpected txType');
        }

        await this.sutAdapter.sendRequests(myArgs);
    }

    async cleanupWorkloadModule() {
        // Cleanup logic if needed
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
