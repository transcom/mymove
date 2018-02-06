import configureStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import {
  ShipmentsReducer,
  createShowShipmentsRequest,
  createShowShipmentsSuccess,
  createShowShipmentsFailure,
} from './ducks';

// SHIPMENTS TEST

describe(' Shipments Reducer', () => {
  it('Should handle SHOW_SHIPMENTS', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = ShipmentsReducer(initialState, {
      type: 'SHOW_SHIPMENTS',
    });

    expect(newState).toEqual({ shipments: null, hasError: false });
  });

  it('Should handle SHOW_SHIPMENTS_SUCCESS', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = ShipmentsReducer(initialState, {
      type: 'SHOW_SHIPMENTS_SUCCESS',
      shipments: ['Sally Shipment'],
    });

    expect(newState).toEqual({
      shipments: ['Sally Shipment'],
      hasError: false,
    });
  });

  it('Should handle SHOW_SHIPMENTS_FAILURE', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = ShipmentsReducer(initialState, {
      type: 'SHOW_SHIPMENTS_FAILURE',
      error: 'Boring',
    });

    expect(newState).toEqual({ shipments: null, hasError: true });
  });
});

describe(' Shipments Actions', () => {
  const initialState = { shipments: null, hasError: false };
  const mockStore = configureStore();
  let store;

  beforeEach(() => {
    store = mockStore(initialState);
  });

  it('Should check action on dispatching ', () => {
    let action;
    store.dispatch(createShowShipmentsRequest());
    store.dispatch(
      createShowShipmentsSuccess([
        {
          id: '11',
          name: 'Sally Shipment',
          pickup_date: new Date(2018, 11, 17).toString(),
          delivery_date: new Date(2018, 11, 19).toString(),
        },
      ]),
    );
    store.dispatch(createShowShipmentsFailure('Tests r not fun.'));
    action = store.getActions();
    // Add expect about what the contents will be.
    expect(action[0].type).toBe('SHOW_SHIPMENTS');
    expect(action[1].type).toBe('SHOW_SHIPMENTS_SUCCESS');
    expect(action[1].shipments).toEqual([
      {
        id: '11',
        name: 'Sally Shipment',
        pickup_date: new Date(2018, 11, 17).toString(),
        delivery_date: new Date(2018, 11, 19).toString(),
      },
    ]);
    expect(action[2].type).toBe('SHOW_SHIPMENTS_FAILURE');
    expect(action[2].error).toEqual('Tests r not fun.');
  });
});

// TODO: Figure out how to mock the Swagger API call
// describe('available shipments async action creators', () => {
//   const middlewares = [ thunk ]
//   const initialState = { shipments: null, hasError: false };
//   const mockStore = configureStore(middlewares)

//   afterEach(() => {
//     fetchMock.reset()
//     fetchMock.restore()
//   })

//   it('creates SHOW_SHIPMENTS_SUCCESS when submitted shipments have been loaded', () => {
//     fetchMock
//       .getOnce('/submitted', { shipments: { shipments: [{'id': 11, 'name': 'Sally Shipment', pickup_date: new Date(2018, 11, 17), delivery_date: new Date(2018, 11, 19)}] }, headers: { 'content-type': 'application/json' } })

//     const expectedActions = [
//       { type: SHOW_SHIPMENTS },
//       { type: SHOW_SHIPMENTS_SUCCESS, shipments: { shipments: [{'id': 11, 'name':'Sally Shipment', pickup_date: new Date(2018, 11, 17), delivery_date: new Date(2018, 11, 19) }] } }
//     ]

//     const store = mockStore(initialState)

//     return store.dispatch(loadShipments()).then(() => {
//       // return of async actions
//       expect(store.getActions()).toEqual(expectedActions)
//     })
//   })
// })
