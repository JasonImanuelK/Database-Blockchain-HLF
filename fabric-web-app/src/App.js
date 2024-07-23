import React, { useState } from 'react';
import axios from 'axios';

function App() {
  const [assetId, setAssetId] = useState('');
  const [price, setPrice] = useState('');
  const [buyer, setBuyer] = useState('');
  const [assetName, setAssetName] = useState('');
  const [quantity, setQuantity] = useState('');
  const [location, setLocation] = useState('');
  const [status, setStatus] = useState('');
  const [deliveryDate, setDeliveryDate] = useState('');
  const [transactionResponse, setTransactionResponse] = useState('');
  const [queryResponse, setQueryResponse] = useState('');

  const handleInvoke = async (event) => {
    event.preventDefault();
    try {
      const response = await axios.post('http://localhost:8080/invoke', {
        channelid: 'mychannel',
        chaincodeid: 'asset-transfer',
        function: 'BuyAsset',
        args: [assetId, price, buyer, assetName, quantity, location, status, deliveryDate]
      });

      setTransactionResponse(response.data);
    } catch (error) {
      console.error('Error invoking API:', error);
      setTransactionResponse('Error: ' + error.message);
    }
  };

  const handleQuery = async () => {
    try {
      const response = await axios.get('http://localhost:8080/query', {
        params: {
          channelid: 'mychannel',
          chaincodeid: 'asset-transfer',
          function: 'ReadAsset',
          args: [assetId]
        }
      });

      setQueryResponse(response.data);
    } catch (error) {
      console.error('Error querying API:', error);
      setQueryResponse('Error: ' + error.message);
    }
  };

  return (
    <div className="App">
      <h1>Hyperledger Fabric Web App</h1>

      <form onSubmit={handleInvoke}>
        <label>
          Asset ID:
          <input type="text" value={assetId} onChange={(e) => setAssetId(e.target.value)} />
        </label>
        <br />
        <label>
          Price:
          <input type="text" value={price} onChange={(e) => setPrice(e.target.value)} />
        </label>
        <br />
        <label>
          Buyer:
          <input type="text" value={buyer} onChange={(e) => setBuyer(e.target.value)} />
        </label>
        <br />
        <label>
          Asset Name:
          <input type="text" value={assetName} onChange={(e) => setAssetName(e.target.value)} />
        </label>
        <br />
        <label>
          Quantity:
          <input type="text" value={quantity} onChange={(e) => setQuantity(e.target.value)} />
        </label>
        <br />
        <label>
          Location:
          <input type="text" value={location} onChange={(e) => setLocation(e.target.value)} />
        </label>
        <br />
        <label>
          Status:
          <input type="text" value={status} onChange={(e) => setStatus(e.target.value)} />
        </label>
        <br />
        <label>
          Delivery Date:
          <input type="date" value={deliveryDate} onChange={(e) => setDeliveryDate(e.target.value)} />
        </label>
        <br />
        <button type="submit">Invoke Transaction</button>
      </form>

      <div>
        <button onClick={handleQuery}>Query Chaincode</button>
        <p>Query Response: {queryResponse}</p>
      </div>

      <div>
        <h3>Transaction Response:</h3>
        <p>{transactionResponse}</p>
      </div>
    </div>
  );
}

export default App;

