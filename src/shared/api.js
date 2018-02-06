import Swagger from 'swagger-client';
let client = null;

async function ensureClientIsLoaded() {
  if (!client) {
    client = await Swagger('api/v1/swagger.yaml');
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

export async function AvailableShipmentsIndex() {
  await ensureClientIsLoaded();
  const response = availableResponse;
  // TODO (Rebecca): Fill in response from swagger api
  // const response = await client.apis.shipments.indexIssues();
  if (!response.ok) {
    throw new Error(
      `failed to load shipments index due to server error:
      ${response.url}: ${response.statusText}`,
    );
  }
  return response.body;
}

export async function AwardedShipmentsIndex() {
  await ensureClientIsLoaded();
  const response = awardedResponse;
  // TODO (Rebecca): Fill in response from swagger api
  // const response = await client.apis.shipments.indexIssues();
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
    "updated_at":"2018-01-31T18:13:15.232Z"},
  {
    "created_at":"2018-02-05T13:31:30.396Z",
    "name":"Nino Shipment.",
    "id":"30c5bd76-d917-4b3d-b19b-a1f781a563c5",
    "updated_at":"2018-02-05T13:31:30.396Z"
  },{
    "created_at":"2018-02-05T13:35:11.538Z",
    "name":"REK Shipment.",
    "id":"e8c3a9a9-333f-4e2a-a98d-757dabeba8ce",
    "updated_at":"2018-02-05T13:35:11.538Z"
  }]`,
  body: [
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Bob and Andie.',
      id: 'ab1eace7-ec68-4794-883d-bc6db16f00fe',
      updated_at: '2018-01-31T18:13:15.232Z',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Nino Shipment.',
      id: '30c5bd76-d917-4b3d-b19b-a1f781a563c5',
      updated_at: '2018-01-31T18:13:15.232Z',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'REK Shipment.',
      id: 'e8c3a9a9-333f-4e2a-a98d-757dabeba8ce',
      updated_at: '2018-01-31T18:13:15.232Z',
    },
  ],
  url: 'http://localhost:3000/api/v1/shipments',
};

const awardedResponse = {
  ok: true,
  status: 200,
  statusText: 'OK',
  data: `[
  {
    "created_at":"2018-01-31T18:13:15.232Z",
    "name":"Shirley Shipment",
    "id":"abs1eace7-ec68-4794-883d-bc6db16f00fe",
    "updated_at":"2018-01-31T18:13:15.232Z",
    "traffic_distribution_list": "Johns Movers"
  },
  {
    "created_at":"2018-02-05T13:31:30.396Z",
    "name":"PA Dutch Move",
    "id":"30c5bd76-d917-4b3d-b19b-a1f781a563c5",
    "updated_at":"2018-02-05T13:31:30.396Z",
    "traffic_distribution_list": "Moovers and Shakers"
  },
  {
    "created_at":"2018-02-05T13:35:11.538Z",
    "name":"Ft Funston move",
    "id":"e8c3a9a9-333f-4e2a-a98d-757dabeba8ce",
    "updated_at":"2018-02-05T13:35:11.538Z",
   "traffic_distribution_list": "Shleppers"
 }]`,
  body: [
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Shirley Shipment.',
      id: 'abs1eace7-ec68-4794-883d-bc6db16f00fe',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list: 'Johns Movers',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'PA Dutch Move',
      id: '30c5bd76-d917-4b3d-b19b-a1f781a563c5',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list: 'Moovers and Shakers',
    },
    {
      created_at: '2018-01-31T18:13:15.232Z',
      name: 'Ft Funston move',
      id: 'e8c3a9a9-333f-4e2a-a98d-757dabeba8ce',
      updated_at: '2018-01-31T18:13:15.232Z',
      traffic_distribution_list: 'Shleppers',
    },
  ],
  url: 'http://localhost:3000/api/v1/shipments',
};
