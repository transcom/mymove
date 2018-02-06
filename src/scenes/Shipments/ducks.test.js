import configureStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import {
  availableShipmentsReducer,
  createShowAvailableShipmentsRequest,
  createShowAvailableShipmentsSuccess,
  createShowAvailableShipmentsFailure,
  awardedShipmentsReducer,
  createShowAwardedShipmentsRequest,
  createShowAwardedShipmentsSuccess,
  createShowAwardedShipmentsFailure,
} from './ducks';

// AVAILABLE SHIPMENTS TEST

describe('Available Shipments Reducer', () => {
  it('Should handle SHOW_AVAILABLE_SHIPMENTS', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = availableShipmentsReducer(initialState, {
      type: 'SHOW_AVAILABLE_SHIPMENTS',
    });

    expect(newState).toEqual({ shipments: null, hasError: false });
  });

  it('Should handle SHOW_AVAILABLE_SHIPMENTS_SUCCESS', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = availableShipmentsReducer(initialState, {
      type: 'SHOW_AVAILABLE_SHIPMENTS_SUCCESS',
      shipments: 'Sally Shipment',
    });

    expect(newState).toEqual({ shipments: 'Sally Shipment', hasError: false });
  });

  it('Should handle SHOW_AVAILABLE_SHIPMENTS_FAILURE', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = availableShipmentsReducer(initialState, {
      type: 'SHOW_AVAILABLE_SHIPMENTS_FAILURE',
      error: 'Boring',
    });

    expect(newState).toEqual({ shipments: null, hasError: true });
  });
});

describe('Available Shipments Actions', () => {
  const initialState = { shipments: null, hasError: false };
  const mockStore = configureStore();
  let store;

  beforeEach(() => {
    store = mockStore(initialState);
  });

  it('Should check action on dispatching ', () => {
    let action;
    store.dispatch(createShowAvailableShipmentsRequest());
    store.dispatch(
      createShowAvailableShipmentsSuccess([
        { id: '11', name: 'Sally Shipment' },
      ]),
    );
    store.dispatch(createShowAvailableShipmentsFailure('Tests r not fun.'));
    action = store.getActions();
    // Add expect about what the contents will be.
    expect(action[0].type).toBe('SHOW_AVAILABLE_SHIPMENTS');
    expect(action[1].type).toBe('SHOW_AVAILABLE_SHIPMENTS_SUCCESS');
    expect(action[1].shipments).toEqual([{ id: '11', name: 'Sally Shipment' }]);
    expect(action[2].type).toBe('SHOW_AVAILABLE_SHIPMENTS_FAILURE');
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

//   it('creates SHOW_AVAILABLE_SHIPMENTS_SUCCESS when submitted shipments have been loaded', () => {
//     fetchMock
//       .getOnce('/submitted', { shipments: { shipments: [{'id': 11, 'name': 'Sally Shipment'}] }, headers: { 'content-type': 'application/json' } })

//     const expectedActions = [
//       { type: SHOW_AVAILABLE_SHIPMENTS },
//       { type: SHOW_AVAILABLE_SHIPMENTS_SUCCESS, shipments: { shipments: [{'id': 11, 'name':'Sally Shipment'}] } }
//     ]

//     const store = mockStore(initialState)

//     return store.dispatch(loadAvailableShipments()).then(() => {
//       // return of async actions
//       expect(store.getActions()).toEqual(expectedActions)
//     })
//   })
// })

// AWARDED SHIPMENTS TESTS
describe('Awarded Shipments Reducer', () => {
  it('Should handle SHOW_AWARDED_SHIPMENTS', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = awardedShipmentsReducer(initialState, {
      type: 'SHOW_AWARDED_SHIPMENTS',
    });

    expect(newState).toEqual({ shipments: null, hasError: false });
  });

  it('Should handle SHOW_AWARDED_SHIPMENTS_SUCCESS', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = awardedShipmentsReducer(initialState, {
      type: 'SHOW_AWARDED_SHIPMENTS_SUCCESS',
      shipments: 'Sally Shipment',
    });

    expect(newState).toEqual({ shipments: 'Sally Shipment', hasError: false });
  });

  it('Should handle SHOW_AWARDED_SHIPMENTS_FAILURE', () => {
    const initialState = { shipments: null, hasError: false };

    const newState = awardedShipmentsReducer(initialState, {
      type: 'SHOW_AWARDED_SHIPMENTS_FAILURE',
      error: 'Boring',
    });

    expect(newState).toEqual({ shipments: null, hasError: true });
  });
});

describe('Awarded Shipments Actions', () => {
  const initialState = { shipments: null, hasError: false };
  const mockStore = configureStore();
  let store;

  beforeEach(() => {
    store = mockStore(initialState);
  });

  it('Should check action on dispatching ', () => {
    let action;
    store.dispatch(createShowAwardedShipmentsRequest());
    store.dispatch(
      createShowAwardedShipmentsSuccess([{ id: '11', name: 'Sally Shipment' }]),
    );
    store.dispatch(createShowAwardedShipmentsFailure('Tests r not fun.'));
    action = store.getActions();
    // Add expect about what the contents will be.
    expect(action[0].type).toBe('SHOW_AWARDED_SHIPMENTS');
    expect(action[1].type).toBe('SHOW_AWARDED_SHIPMENTS_SUCCESS');
    expect(action[1].shipments).toEqual([{ id: '11', name: 'Sally Shipment' }]);
    expect(action[2].type).toBe('SHOW_AWARDED_SHIPMENTS_FAILURE');
    expect(action[2].error).toEqual('Tests r not fun.');
  });
});

// TODO: Figure out how to mock the Swagger API call
// describe('async action creators', () => {
//   const middlewares = [ thunk ]
//   const initialState = { shipments: null, hasError: false };
//   const mockStore = configureStore(middlewares)

//   afterEach(() => {
//     fetchMock.reset()
//     fetchMock.restore()
//   })

//   it('creates SHOW_AWARDED_SHIPMENTS_SUCCESS when submitted shipments have been loaded', () => {
//     fetchMock
//       .getOnce('/submitted', { shipments: { shipments: [{'id': 11, 'name': 'Sally Shipment'}] }, headers: { 'content-type': 'application/json' } })

//     const expectedActions = [
//       { type: SHOW_AWARDED_SHIPMENTS },
//       { type: SHOW_AWARDED_SHIPMENTS_SUCCESS, shipments: { shipments: [{'id': 11, 'name':'Sally Shipment'}] } }
//     ]

//     const store = mockStore(initialState)

//     return store.dispatch(loadAvailableShipments()).then(() => {
//       // return of async actions
//       expect(store.getActions()).toEqual(expectedActions)
//     })
//   })
// })
