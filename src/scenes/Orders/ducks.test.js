import {
  CREATE_OR_UPDATE_ORDERS,
  GET_ORDERS,
  SHOW_CURRENT_ORDERS,
  DELETE_UPLOAD,
  ADD_UPLOADS,
  ordersReducer,
} from './ducks';
import loggedInUserPayload, { emptyPayload } from 'shared/User/sampleLoggedInUserPayload';

const expectedOrders = {
  id: '51953e97-25a7-430c-ba6d-3bd980a38b71',
  has_dependents: false,
  status: 'DRAFT',
  issue_date: '2018-05-11',
  new_duty_station: {
    address: {
      city: 'Fort Worth',
      country: 'United States',
      postal_code: '76127',
      state: 'TX',
      street_address_1: 'n/a',
    },
    affiliation: 'NAVY',
    created_at: '2018-05-20T18:36:45.034Z',
    id: '44db8bfb-db7c-4c8d-bc08-5d683c4469ed',
    name: 'NAS Fort Worth',
    updated_at: '2018-05-20T18:36:45.034Z',
  },
  orders_type: 'PERMANENT_CHANGE_OF_STATION',
  report_by_date: '2018-05-29',
  service_member_id: '1694e00e-17ff-43fe-af6d-ab0519a18ff2',
  uploaded_orders: {
    id: '24f18674-eec7-4c1f-b8c0-cb343a8c4f77',
    name: 'uploaded_orders',
    service_member_id: '1694e00e-17ff-43fe-af6d-ab0519a18ff2',
    uploads: [
      {
        bytes: 3932969,
        content_type: 'image/jpeg',
        created_at: '2018-05-25T21:38:06.235Z',
        filename: 'last vacccination.jpg',
        id: 'd56df2e3-1481-4dff-9a02-ef5c6bcae491',
        updated_at: '2018-05-25T21:38:06.235Z',
        url:
          '/storage/documents/24f18674-eec7-4c1f-b8c0-cb343a8c4f77/uploads/d56df2e3-1481-4dff-9a02-ef5c6bcae491?contentType=image%2Fjpeg',
      },
      {
        bytes: 58036,
        content_type: 'image/png',
        created_at: '2018-05-25T21:38:57.655Z',
        filename: 'image (2).png',
        id: 'e2010a83-ac1e-45a2-9eb1-4e144c443c41',
        updated_at: '2018-05-25T21:38:57.655Z',
        url:
          '/storage/documents/24f18674-eec7-4c1f-b8c0-cb343a8c4f77/uploads/e2010a83-ac1e-45a2-9eb1-4e144c443c41?contentType=image%2Fpng',
      },
    ],
  },
};
const ordersWithOneUpload = {
  id: '51953e97-25a7-430c-ba6d-3bd980a38b71',
  has_dependents: false,
  issue_date: '2018-05-11',
  status: 'DRAFT',
  new_duty_station: {
    address: {
      city: 'Fort Worth',
      country: 'United States',
      postal_code: '76127',
      state: 'TX',
      street_address_1: 'n/a',
    },
    affiliation: 'NAVY',
    created_at: '2018-05-20T18:36:45.034Z',
    id: '44db8bfb-db7c-4c8d-bc08-5d683c4469ed',
    name: 'NAS Fort Worth',
    updated_at: '2018-05-20T18:36:45.034Z',
  },
  orders_type: 'PERMANENT_CHANGE_OF_STATION',
  report_by_date: '2018-05-29',
  service_member_id: '1694e00e-17ff-43fe-af6d-ab0519a18ff2',
  uploaded_orders: {
    id: '24f18674-eec7-4c1f-b8c0-cb343a8c4f77',
    name: 'uploaded_orders',
    service_member_id: '1694e00e-17ff-43fe-af6d-ab0519a18ff2',
    uploads: [
      {
        bytes: 3932969,
        content_type: 'image/jpeg',
        created_at: '2018-05-25T21:38:06.235Z',
        filename: 'last vacccination.jpg',
        id: 'd56df2e3-1481-4dff-9a02-ef5c6bcae491',
        updated_at: '2018-05-25T21:38:06.235Z',
        url:
          '/storage/documents/24f18674-eec7-4c1f-b8c0-cb343a8c4f77/uploads/d56df2e3-1481-4dff-9a02-ef5c6bcae491?contentType=image%2Fjpeg',
      },
    ],
  },
};
const ordersPayload = Object.freeze({
  created_at: '2018-05-25T21:36:10.219Z',
  has_dependents: false,
  id: '51953e97-25a7-430c-ba6d-3bd980a38b71',
  issue_date: '2018-05-11',
  status: 'DRAFT',
  moves: [
    {
      created_at: '2018-05-25T21:36:10.235Z',
      id: '593cc830-1a3e-44b3-ba5a-8809f02dfa7d',
      locator: 'WUMGLQ',
      orders_id: '51953e97-25a7-430c-ba6d-3bd980a38b71',
      personally_procured_moves: [
        {
          destination_postal_code: '76127',
          incentive_estimate_min: 1495409,
          incentive_estimate_max: 1652821,
          has_additional_postal_code: false,
          has_requested_advance: false,
          has_sit: false,
          id: 'cd67c9e4-ef59-45e5-94bc-767aaafe559e',
          pickup_postal_code: '80913',
          original_move_date: '2018-06-28',
          size: 'L',
          status: 'DRAFT',
          weight_estimate: 9000,
        },
      ],
      selected_move_type: 'PPM',
      status: 'DRAFT',
    },
  ],
  new_duty_station: {
    address: {
      city: 'Fort Worth',
      country: 'United States',
      postal_code: '76127',
      state: 'TX',
      street_address_1: 'n/a',
    },
    affiliation: 'NAVY',
    created_at: '2018-05-20T18:36:45.034Z',
    id: '44db8bfb-db7c-4c8d-bc08-5d683c4469ed',
    name: 'NAS Fort Worth',
    updated_at: '2018-05-20T18:36:45.034Z',
  },
  orders_type: 'PERMANENT_CHANGE_OF_STATION',
  report_by_date: '2018-05-29',
  service_member_id: '1694e00e-17ff-43fe-af6d-ab0519a18ff2',
  updated_at: '2018-05-25T21:39:02.429Z',
  uploaded_orders: {
    id: '24f18674-eec7-4c1f-b8c0-cb343a8c4f77',
    name: 'uploaded_orders',
    service_member_id: '1694e00e-17ff-43fe-af6d-ab0519a18ff2',
    uploads: [
      {
        bytes: 3932969,
        content_type: 'image/jpeg',
        created_at: '2018-05-25T21:38:06.235Z',
        filename: 'last vacccination.jpg',
        id: 'd56df2e3-1481-4dff-9a02-ef5c6bcae491',
        updated_at: '2018-05-25T21:38:06.235Z',
        url:
          '/storage/documents/24f18674-eec7-4c1f-b8c0-cb343a8c4f77/uploads/d56df2e3-1481-4dff-9a02-ef5c6bcae491?contentType=image%2Fjpeg',
      },
      {
        bytes: 58036,
        content_type: 'image/png',
        created_at: '2018-05-25T21:38:57.655Z',
        filename: 'image (2).png',
        id: 'e2010a83-ac1e-45a2-9eb1-4e144c443c41',
        updated_at: '2018-05-25T21:38:57.655Z',
        url:
          '/storage/documents/24f18674-eec7-4c1f-b8c0-cb343a8c4f77/uploads/e2010a83-ac1e-45a2-9eb1-4e144c443c41?contentType=image%2Fpng',
      },
    ],
  },
});
describe('orders Reducer', () => {
  describe('GET_LOGGED_IN_USER', () => {
    it('Should handle GET_LOGGED_IN_USER.success', () => {
      const initialState = {};
      const newState = ordersReducer(initialState, loggedInUserPayload);

      expect(newState).toEqual({
        currentOrders: { ...expectedOrders },
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    });
    it('Should handle empty payload', () => {
      const initialState = {};
      const newState = ordersReducer(initialState, emptyPayload);

      expect(newState).toEqual({
        currentOrders: null,
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    });
  });

  describe('CREATE_OR_UPDATE_ORDERS', () => {
    it('Should handle CREATE_OR_UPDATE_ORDERS_SUCCESS', () => {
      const initialState = {};

      const newState = ordersReducer(initialState, {
        type: CREATE_OR_UPDATE_ORDERS.success,
        payload: ordersPayload,
      });

      expect(newState).toEqual({
        currentOrders: { ...expectedOrders },
        error: null,
        hasSubmitError: false,
        hasSubmitSuccess: true,
      });
    });

    it('Should handle CREATE_OR_UPDATE_ORDERS_FAILURE', () => {
      const initialState = {};

      const newState = ordersReducer(initialState, {
        type: CREATE_OR_UPDATE_ORDERS.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentOrders: {},
        error: 'No bueno.',
        hasSubmitError: true,
        hasSubmitSuccess: false,
      });
    });
  });

  describe('GET_ORDERS', () => {
    it('Should handle GET_ORDERS_SUCCESS', () => {
      const initialState = {};
      const newState = ordersReducer(initialState, {
        type: GET_ORDERS.success,
        payload: ordersPayload,
      });

      expect(newState).toEqual({
        currentOrders: { ...expectedOrders },
        error: null,
        hasLoadError: false,
        hasLoadSuccess: true,
      });
    });

    it('Should handle GET_ORDERS_FAILURE', () => {
      const initialState = {};

      const newState = ordersReducer(initialState, {
        type: GET_ORDERS.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentOrders: null,
        hasLoadError: true,
        hasLoadSuccess: false,
        error: 'No bueno.',
      });
    });
  });
  describe('SHOW_CURRENT_ORDERS', () => {
    it('Should handle SHOW_CURRENT_ORDERS_SUCCESS', () => {
      const initialState = {};
      const newState = ordersReducer(initialState, {
        type: SHOW_CURRENT_ORDERS.success,
        payload: ordersPayload,
      });

      expect(newState).toEqual({
        currentOrders: { ...expectedOrders },
        showCurrentOrdersError: false,
        showCurrentOrdersSuccess: true,
      });
    });

    it('Should handle SHOW_CURRENT_ORDERS_FAILURE', () => {
      const initialState = {};

      const newState = ordersReducer(initialState, {
        type: SHOW_CURRENT_ORDERS.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentOrders: null,
        showCurrentOrdersError: true,
        error: 'No bueno.',
      });
    });
  });
  describe('DELETE_UPLOAD', () => {
    it('Should handle DELETE_UPLOAD_SUCCESS', () => {
      const initialState = { currentOrders: expectedOrders };
      const newState = ordersReducer(initialState, {
        type: DELETE_UPLOAD.success,
        payload: 'e2010a83-ac1e-45a2-9eb1-4e144c443c41',
      });
      expect(expectedOrders.uploaded_orders.uploads.length).toEqual(2);
      expect(newState).toEqual({
        currentOrders: ordersWithOneUpload,
        hasSubmitError: false,
        hasSubmitSuccess: true,
        error: null,
      });
    });

    it('Should handle DELETE_UPLOAD_FAILURE', () => {
      const initialState = { currentOrders: expectedOrders };

      const newState = ordersReducer(initialState, {
        payload: ordersPayload,
        type: DELETE_UPLOAD.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentOrders: expectedOrders,
        hasSubmitError: true,
        hasSubmitSuccess: false,
        error: 'No bueno.',
      });
    });
  });
  describe('ADD_UPLOADS', () => {
    it('Should handle ADD_UPLOADS_SUCCESS', () => {
      const initialState = { currentOrders: ordersWithOneUpload };
      const newState = ordersReducer(initialState, {
        type: ADD_UPLOADS.success,
        payload: [expectedOrders.uploaded_orders.uploads[1]],
      });

      expect(initialState.currentOrders.uploaded_orders.uploads.length).toEqual(1);
      expect(newState).toEqual({
        currentOrders: expectedOrders,
        hasSubmitError: false,
        hasSubmitSuccess: true,
        error: null,
      });
    });

    it('Should handle ADD_UPLOADS_FAILURE', () => {
      const initialState = { currentOrders: ordersWithOneUpload };

      const newState = ordersReducer(initialState, {
        payload: ordersPayload,
        type: ADD_UPLOADS.failure,
        error: 'No bueno.',
      });

      expect(newState).toEqual({
        currentOrders: { ...initialState.currentOrders },
        hasSubmitError: true,
        hasSubmitSuccess: false,
        error: 'No bueno.',
      });
    });
  });
});
