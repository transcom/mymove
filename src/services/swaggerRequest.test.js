import { makeSwaggerRequest, normalizeResponse } from './swaggerRequest';

function mockGetClient(operationMock) {
  return Promise.resolve({
    apis: {
      shipments: {
        getShipment: operationMock,
      },
    },
    spec: {
      paths: {
        'shipments/{shipmentID}': {
          get: {
            consumes: ['application/json'],
            description: 'Returns a Shipment tied to the move ID',
            operationId: 'getShipment',
            parameters: [
              {
                description: 'UUID of the Shipment being updated',
                format: 'uuid',
                in: 'path',
                name: 'shipmentId',
                required: true,
                type: 'string',
              },
            ],
            produces: ['application/json'],
            responses: {
              200: {
                description: 'Returns Shipment for hhg move',
                schema: {
                  type: 'object',
                  properties: {},
                  $$ref: '#/definitions/Shipment',
                },
              },
              400: {
                description: 'Bad request',
              },
            },
            summary: 'Returns a Shipment for the given move',
            tags: ['shipments'],
          },
        },
      },
    },
  });
}

describe('normalizeResponse', () => {
  it('normalizes data with a valid schemaKey', () => {
    const rawData = {
      id: 'abcd-1234',
      type: 'test shipment',
    };

    const expectedData = {
      shipments: {
        'abcd-1234': {
          id: 'abcd-1234',
          type: 'test shipment',
        },
      },
    };

    expect(normalizeResponse(rawData, 'shipment')).toEqual(expectedData);
  });
});

describe('makeSwaggerRequest', () => {
  it('makes a successful request', async () => {
    const mockResponse = {
      ok: true,
      status: 200,
      body: {
        id: 'abcd-1234',
        type: 'test shipment',
      },
    };

    const expectedData = {
      shipments: {
        'abcd-1234': {
          id: 'abcd-1234',
          type: 'test shipment',
        },
      },
    };

    const opMock = jest.fn(() => Promise.resolve(mockResponse));
    const mockClient = await mockGetClient(opMock);
    const request = await makeSwaggerRequest(mockClient, 'shipments.getShipment', { shipmentID: 'abcd-1234' });

    expect(request).toEqual(expectedData);
  });

  it('makes a failed request', async () => {
    const mockResponse = {
      ok: false,
      status: 400,
      body: {
        error: {
          message: 'Test error message',
        },
      },
    };

    const errorResponse = () => Promise.reject(mockResponse);
    const opMock = jest.fn(() => errorResponse());

    const mockClient = await mockGetClient(opMock);

    await makeSwaggerRequest(mockClient, 'shipments.getShipment', { shipmentID: 'abcd-1234' }).catch((error) => {
      expect(error).toEqual(mockResponse);
    });
  });

  it('errors if the request is unknown', async () => {
    const mockResponse = {
      ok: true,
      status: 200,
      body: {
        id: 'abcd-1234',
        type: 'test shipment',
      },
    };

    const opMock = jest.fn(() => Promise.resolve(mockResponse));
    const mockClient = await mockGetClient(opMock);

    await makeSwaggerRequest(mockClient, 'unknown', { shipmentID: 'abcd-1234' }).catch((error) => {
      expect(error).toEqual(new Error(`Operation 'unknown' does not exist!`));
    });
  });

  it('returns the raw response body if normalize is false', async () => {
    const mockResponse = {
      ok: true,
      status: 200,
      body: {
        id: 'abcd-1234',
        type: 'test shipment',
      },
    };

    const opMock = jest.fn(() => Promise.resolve(mockResponse));
    const mockClient = await mockGetClient(opMock);
    const request = await makeSwaggerRequest(
      mockClient,
      'shipments.getShipment',
      { shipmentID: 'abcd-1234' },
      { normalize: false },
    );

    expect(request).toEqual(mockResponse.body);
  });
});
