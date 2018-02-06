import Swagger from 'swagger-client';
let client = null;

async function ensureClientIsLoaded() {
  if (!client) {
    client = await Swagger('/api/v1/swagger.yaml');
  }
}

export async function GetSpec() {
  await ensureClientIsLoaded();
  return client.spec;
}
// because these are async functions,
// they return a promise that is resolved with the return value
// if there is an error, the promise is rejected with that error
export async function IssuesIndex() {
  await ensureClientIsLoaded();
  const response = await client.apis.issues.indexIssues();
  if (!response.ok) {
    throw new Error(
      `failed to load issues index due to server error:
      ${response.url}: ${response.statusText}`,
    );
  }
  return response.body;
}

export async function CreateIssue(issueBody) {
  await ensureClientIsLoaded();
  const response = await client.apis.issues.createIssue({
    createIssuePayload: { description: issueBody },
  });
  if (!response.ok)
    new Error(
      `failed to create issue due to server error:
      ${response.url}: ${response.statusText}`,
    );
}
export async function CreateForm1299(formData) {
  await ensureClientIsLoaded();
  const response = await client.apis.form1299s.createForm1299({
    createForm1299Payload: formData,
  });
  if (!response.ok)
    new Error(
      `failed to create issue due to server error:
      ${response.url}: ${response.statusText}`,
    );
}

export async function ShipmentsIndex() {
  await ensureClientIsLoaded();
  let response;
  // TODO (Rebecca): Fill in response from swagger api
  // const response = await client.apis.shipments.indexShipments();
  if (shipmentsStatus === 'awarded') {
    response = availableResponse;
  } else if (shipmentsStatus === 'available') {
    response = awardedResponse;
  } else {
    response = allResponse;
  }

  if (!response.ok) {
    throw new Error(
      `failed to load shipments index due to server error:
      ${response.url}: ${response.statusText}`,
    );
  }
  return response.body;
}

const availableResponse = {
  ok: true,
  status: 200,
  statusText: 'OK',
  data: `[
  {
    "created_at":"2018-01-31T18:13:15.232Z",
    "name":"Bob and Andie.",
    "id":"ab1eace7-ec68-4794-883d-bc6db16f00fe",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-01-31T18:13:15.232Z"},
    "traffic_distribution_list_id": "3015bd76-d187-4b3d-b19b-a1f781a563c5"
  {
    "created_at":"2018-02-05T13:31:30.396Z",
    "name":"Nino Shipment.",
    "id":"30c5bd76-d917-4b3d-b19b-a1f781a563c5",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-02-05T13:31:30.396Z",
    "traffic_distribution_list_id": "3015bd76-d187-4b3d-g19b-a1f781a563c5"
  },{
    "created_at":"2018-02-05T13:35:11.538Z",
    "name":"REK Shipment.",
    "id":"e8c3a9a9-333f-4e2a-a98d-757dabeba8ce",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-02-05T13:35:11.538Z",
    "traffic_distribution_list_id": "3015bd76-d137-4b3d-b19b-a1f781a563c5"
  }]`,
  body: [
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Bob and Andie.',
      id: 'ab1eace7-ec68-4794-883d-bc6db16f00fe',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '3015bd76-d187-4b3d-b19b-a1f781a563c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Nino Shipment.',
      id: '30c5bd76-d917-4b3d-b19b-a1f781a563c5',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '3015bd76-d187-4b3d-g19b-a1f781a563c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'REK Shipment.',
      id: 'e8c3a9a9-333f-4e2a-a98d-757dabeba8ce',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '3015bd76-d137-4b3d-b19b-a1f781a563c5',
    },
  ],
  url: 'http://localhost:3000/api/v1/available',
};

const awardedResponse = {
  ok: true,
  status: 200,
  statusText: 'OK',
  data: `[
  {
    "created_at":'2018-01-31T18:13:15.232Z',
    "name":"Shirley Shipment",
    "id":"abs1eace7-ec68-4794-883d-bc6db16f00fe",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-01-31T18:13:15.232Z",
    "traffic_distribution_list_id": "3015bd76-d187-4b3d-b19b-a1f781a563c5"
  },
  {
    "created_at":"2018-02-05T13:31:30.396Z",
    "name":"PA Dutch Move",
    "id":"30c5bd76-d917-4b3d-b19b-a1f781a563c5",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-02-05T13:31:30.396Z",
    "traffic_distribution_list_id": "3015bd76-d917-4b3d-b193-a1f781a563c5"
  },
  {
    "created_at":"2018-02-05T13:35:11.538Z",
    "name":"Ft Funston move",
    "id":"e8c3a9a9-333f-4e2a-a98d-757dabeba8ce",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-02-05T13:35:11.538Z",
    "traffic_distribution_list_id": "3015bd76-d917-4b3d-b19b-a1f781a563c5",
 }]`,
  body: [
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Shirley Shipment.',
      id: 'abs1eace7-ec68-4794-883d-bc6db16f00fe',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '30c5bd76-d917-4b3d-b19b-a1f781aa63c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'PA Dutch Move',
      id: '30c5bd76-d917-4b3d-b19b-a1f781a563c5',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '30c5bd76-d917-4bs3d-b19b-a1f781a563c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Ft Funston move',
      id: 'e8c3a9a9-333f-4e2a-a98d-757dabeba8ce',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '30c5sd76-d917-4b3d-b19b-a1f781a563c5',
    },
  ],
  url: 'http://localhost:3000/api/v1/awarded',
};

const allResponse = {
  ok: true,
  status: 200,
  statusText: 'OK',
  data: `[
  {
    "created_at":'2018-01-31T18:13:15.232Z',
    "name":"Shirley Shipment",
    "id":"abs1eace7-ec68-4794-883d-bc6db16f00fe",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-01-31T18:13:15.232Z",
    "traffic_distribution_list_id": "3015bd76-d187-4b3d-b19b-a1f781a563c5"
  },
  {
    "created_at":"2018-02-05T13:31:30.396Z",
    "name":"PA Dutch Move",
    "id":"30c5bd76-d917-4b3d-b19b-a1f781a563c5",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-02-05T13:31:30.396Z",
    "traffic_distribution_list_id": "3015bd76-d917-4b3d-b193-a1f781a563c5"
  },
  {
    "created_at":"2018-02-05T13:35:11.538Z",
    "name":"Ft Funston move",
    "id":"e8c3a9a9-333f-4e2a-a98d-757dabeba8ce",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-02-05T13:35:11.538Z",
    "traffic_distribution_list_id": "3015bd76-d917-4b3d-b19b-a1f781a563c5",
 },
 {
    "created_at":"2018-01-31T18:13:15.232Z",
    "name":"Bob and Andie.",
    "id":"ab1eace7-ec68-4794-883d-bc6db16f00fe",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-01-31T18:13:15.232Z"},
    "traffic_distribution_list_id": "3015bd76-d187-4b3d-b19b-a1f781a563c5"
  {
    "created_at":"2018-02-05T13:31:30.396Z",
    "name":"Nino Shipment.",
    "id":"30c5bd76-d917-4b3d-b9b-a1f781a563c5",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-02-05T13:31:30.396Z",
    "traffic_distribution_list_id": "3015bd76-d187-4b3d-g19b-a1f781a563c5"
  },{
    "created_at":"2018-02-05T13:35:11.538Z",
    "name":"REK Shipment.",
    "id":"e8c3ada9-333f-4e2a-a98d-757dabeba8ce",
    "pickup_date":"2018-12-05T13:35:11.538Z",
    "delivery_date":"2018-12-08T13:35:11.538Z",
    "updated_at":"2018-02-05T13:35:11.538Z",
    "traffic_distribution_list_id": "3015bd76-d137-4b3d-b19b-a1f781a563c5"
  }]`,
  body: [
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Bob and Andie.',
      id: 'ab1eace7-ec68-4794-883d-bc6db16f00fe',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '3015bd76-d187-4b3d-b19b-a1f781a563c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Nino Shipment.',
      id: '30c5bd76-d917-4b3d-b19b-a1f781a563c5',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '3015bd76-d187-4b3d-g19b-a1f781a563c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'REK Shipment.',
      id: 'e8c3a99-333f-4e2a-a98d-757dabeba8ce',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '3015bd76-d137-4b3d-b19b-a1f781a563c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Shirley Shipment.',
      id: 'abs1eace7-ec68-4794-883d-bc6db16f00fe',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '30c5bd76-d917-4b3d-b19b-a1f781aa63c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'PA Dutch Move',
      id: '30c5bd76-d917-4b3d-b19-a1f781a563c5',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '30c5bd76-d917-4bs3d-b19b-a1f781a563c5',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Ft Funston move',
      id: 'e8c3a9a9-333f-4e2a-a98d-757dabeba8ce',
      pickup_date: '2018-12-05T13:35:11.538Z',
      delivery_date: '2018-12-08T13:35:11.538Z',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list_id: '30c5sd76-d917-4b3d-b19b-a1f781a563c5',
    },
  ],
  url: 'http://localhost:3000/api/v1/available',
};
