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
    it('makes a successful request', function(done) {
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

      action(dispatch, getState, { client, schema });

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

      setTimeout(function() {
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

      done();
    });

    it('makes a failed request', function(done) {
      const dispatch = jest.fn();
      const getState = jest.fn();

      const response = {
        ok: true,
        status: 401,
        body: {
          shipment: {
            id: 'abcd-1234',
          },
        },
      };

      let rejectCallback;
      const opMock = jest.fn(function() {
        return new Promise(function(resolve, reject) {
          rejectCallback = reject;
        });
      });

      const client = mockClient(opMock);

      const action = request.swaggerRequest(
        'shipments.getShipment',
        { shipmentID: 'abcd-1234' },
        { label: 'testRequest' },
      );

      action(dispatch, getState, { client, schema });

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

      expect(function() {
        rejectCallback(response);
      }).toThrow();

      setTimeout(function() {
        expect(dispatch).toHaveBeenLastCalledWith(
          expect.objectContaining({
            type: '@@swagger/shipments.getShipment/FAILURE',
            label: 'testRequest',
          }),
        );
      });

      done();
    });
  });
});
