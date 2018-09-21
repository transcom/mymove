import * as request from './request';
import * as schema from 'shared/Entities/schema';

function mockClient(operationMock) {
  return {
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
  };
}

describe('swaggerRequest', function() {
  describe('making a request', function() {
    it('makes a successful request', function() {
      expect.assertions(3);

      const dispatch = jest.fn();
      const getState = jest.fn();

      const response = {
        ok: true,
        status: 200,
        body: {
          shipment: {
            id: 'abcd-1234',
          },
        },
      };

      let resolveCallback;
      const opMock = jest.fn(function() {
        return new Promise(function(resolve, reject) {
          resolveCallback = resolve;
        });
      });
      const client = mockClient(opMock);

      const action = request.swaggerRequest(
        'shipments.getShipment',
        { shipmentID: 'abcd-1234' },
        { label: 'testRequest' },
      );

      const result = action(dispatch, getState, { client, schema });

      expect(dispatch).toHaveBeenLastCalledWith(
        expect.objectContaining({
          type: '@@swagger/shipments.getShipment/START',
          label: 'testRequest',
          request: expect.objectContaining({
            id: expect.any(String),
            operationPath: 'shipments.getShipment',
            params: { shipmentID: 'abcd-1234' },
            start: expect.any(Date),
            isLoading: true,
          }),
        }),
      );

      resolveCallback(response);

      result.then(function(response) {
        expect(dispatch).toHaveBeenLastCalledWith(
          expect.objectContaining({
            type: '@@swagger/shipments.getShipment/SUCCESS',
            label: 'testRequest',
            entities: expect.objectContaining({
              shipments: expect.objectContaining({
                'abcd-1234': expect.objectContaining({
                  id: 'abcd-1234',
                }),
              }),
            }),
          }),
        );
      });

      return expect(result).resolves.toEqual(response);
    });

    it('makes a failed request', function() {
      expect.assertions(3);

      const dispatch = jest.fn();
      const getState = jest.fn();

      const response = {
        ok: true,
        status: 400,
        body: {
          shipment: {
            id: 'abcd-1234',
          },
        },
      };

      let rejectCallback;
      let promise;
      const opMock = jest.fn(function() {
        promise = new Promise(function(resolve, reject) {
          rejectCallback = reject;
        });
        return promise;
      });

      const client = mockClient(opMock);

      const action = request.swaggerRequest(
        'shipments.getShipment',
        { shipmentID: 'abcd-1234' },
        { label: 'testRequest' },
      );

      const result = action(dispatch, getState, { client, schema });

      expect(dispatch).toHaveBeenLastCalledWith(
        expect.objectContaining({
          type: '@@swagger/shipments.getShipment/START',
          label: 'testRequest',
          request: expect.objectContaining({
            id: expect.any(String),
            operationPath: 'shipments.getShipment',
            params: { shipmentID: 'abcd-1234' },
            start: expect.any(Date),
            isLoading: true,
          }),
        }),
      );

      result.catch(function(response) {
        expect(dispatch).toHaveBeenLastCalledWith(
          expect.objectContaining({
            type: '@@swagger/shipments.getShipment/FAILURE',
            label: 'testRequest',
          }),
        );
      });

      rejectCallback(response);

      return expect(result).rejects.toEqual(response);
    });
  });
});
