import {
  LoadShipment,
  PatchShipment,
  AcceptShipment,
  RejectShipment,
  TransportShipment,
  DeliverShipment,
  CreateServiceAgent,
  IndexServiceAgents,
  UpdateServiceAgent,
  GenerateGBL,
  GetAllShipmentDocuments,
} from './api.js';

import * as ReduxHelpers from 'shared/ReduxHelpers';

// SINGLE RESOURCE ACTION TYPES
const loadShipmentType = 'LOAD_SHIPMENT';
const patchShipmentType = 'PATCH_SHIPMENT';
const acceptShipmentType = 'ACCEPT_SHIPMENT';
const generateGBLType = 'GENERATE_GBL';
const rejectShipmentType = 'REJECT_SHIPMENT';
const transportShipmentType = 'TRANSPORT_SHIPMENT';
const deliverShipmentType = 'TRANSPORT_SHIPMENT';
const loadShipmentDocumentsType = 'LOAD_SHIPMENT_DOCUMENTS';

const indexServiceAgentsType = 'INDEX_SERVICE_AGENTS';
const createServiceAgentsType = 'CREATE_SERVICE_AGENTS';
const updateServiceAgentsType = 'UPDATE_SERVICE_AGENTS';

// MULTIPLE-RESOURCE ACTION TYPES
const loadTspDependenciesType = 'LOAD_TSP_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES
const LOAD_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(loadShipmentType);
const PATCH_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(patchShipmentType);
const ACCEPT_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(
  acceptShipmentType,
);
const REJECT_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(
  rejectShipmentType,
);
const TRANSPORT_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(
  transportShipmentType,
);
const DELIVER_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(
  deliverShipmentType,
);
const LOAD_SHIPMENT_DOCUMENTS = ReduxHelpers.generateAsyncActionTypes(
  loadShipmentDocumentsType,
);

const GENERATE_GBL = ReduxHelpers.generateAsyncActionTypes(generateGBLType);

const INDEX_SERVICE_AGENTS = ReduxHelpers.generateAsyncActionTypes(
  indexServiceAgentsType,
);

const CREATE_SERVICE_AGENTS = ReduxHelpers.generateAsyncActionTypes(
  createServiceAgentsType,
);

const UPDATE_SERVICE_AGENTS = ReduxHelpers.generateAsyncActionTypes(
  updateServiceAgentsType,
);

// MULTIPLE-RESOURCE ACTION TYPES

const LOAD_TSP_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(
  loadTspDependenciesType,
);

// SINGLE-RESOURCE ACTION CREATORS

export const loadShipment = ReduxHelpers.generateAsyncActionCreator(
  loadShipmentType,
  LoadShipment,
);

export const patchShipment = ReduxHelpers.generateAsyncActionCreator(
  patchShipmentType,
  PatchShipment,
);

export const acceptShipment = ReduxHelpers.generateAsyncActionCreator(
  acceptShipmentType,
  AcceptShipment,
);

export const generateGBL = ReduxHelpers.generateAsyncActionCreator(
  generateGBLType,
  GenerateGBL,
);

export const rejectShipment = ReduxHelpers.generateAsyncActionCreator(
  rejectShipmentType,
  RejectShipment,
);

export const transportShipment = ReduxHelpers.generateAsyncActionCreator(
  transportShipmentType,
  TransportShipment,
);

export const deliverShipment = ReduxHelpers.generateAsyncActionCreator(
  deliverShipmentType,
  DeliverShipment,
);

export const getAllShipmentDocuments = ReduxHelpers.generateAsyncActionCreator(
  loadShipmentDocumentsType,
  GetAllShipmentDocuments,
);

export const indexServiceAgents = ReduxHelpers.generateAsyncActionCreator(
  indexServiceAgentsType,
  IndexServiceAgents,
);

export const createServiceAgent = ReduxHelpers.generateAsyncActionCreator(
  createServiceAgentsType,
  CreateServiceAgent,
);

export const updateServiceAgent = ReduxHelpers.generateAsyncActionCreator(
  updateServiceAgentsType,
  UpdateServiceAgent,
);

// MULTIPLE-RESOURCE ACTION CREATORS
//
// These action types typically dispatch to other actions above to
// perform their work and exist to encapsulate when multiple requests
// need to be made in response to a user action.

export function createOrUpdateServiceAgent(shipmentId, serviceAgent) {
  return async function(dispatch, getState) {
    if (serviceAgent.id) {
      return dispatch(updateServiceAgent(serviceAgent));
    } else {
      return dispatch(createServiceAgent(shipmentId, serviceAgent));
    }
  };
}

export function loadShipmentDependencies(shipmentId) {
  const actions = ReduxHelpers.generateAsyncActions(loadTspDependenciesType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      await Promise.all([
        dispatch(loadShipment(shipmentId)),
        dispatch(indexServiceAgents(shipmentId)),
      ]);
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

// Selectors

// Reducer
const initialState = {
  shipmentIsLoading: false,
  shipmentHasLoadSuccess: false,
  shipmentHasLoadError: null,
  shipmentPatchSuccess: false,
  shipmentPatchError: null,
  shipmentIsAccepting: false,
  shipmentHasAcceptError: null,
  shipmentHasAcceptSuccess: false,
  shipmentIsRejecting: false,
  shipmentHasRejectError: null,
  shipmentHasRejectSuccess: false,
  shipmentIsSendingTransport: false,
  shipmentHasTransportError: null,
  shipmentHasTransportSuccess: false,
  shipmentIsDelivering: false,
  shipmentHasDeliverError: null,
  shipmentHasDeliverSuccess: false,
  serviceAgentsAreLoading: false,
  serviceAgentsHasLoadSucces: false,
  serviceAgentsHasLoadError: null,
  serviceAgentIsCreating: false,
  serviceAgentHasCreatedSucces: false,
  serviceAgentHasCreatedError: null,
  serviceAgentIsUpdating: false,
  serviceAgentHasUpdatedSucces: false,
  serviceAgentHasUpdatedError: null,
  loadTspDependenciesHasSuccess: false,
  loadTspDependenciesHasError: null,
  flashMessage: false,
  serviceAgents: [],
  generateGBLSuccess: false,
  generateGBLError: null,
};

export function tspReducer(state = initialState, action) {
  switch (action.type) {
    // SINGLE-RESOURCE ACTION TYPES

    // SHIPMENTS
    case LOAD_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsLoading: true,
        shipmentHasLoadSuccess: false,
      });
    case LOAD_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsLoading: false,
        shipmentHasLoadSuccess: true,
        shipmentHasLoadError: false,
        shipment: action.payload,
      });
    case LOAD_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsLoading: false,
        shipmentHasLoadSuccess: false,
        shipmentHasLoadError: null,
        shipment: null,
        error: action.error.message,
      });
    case PATCH_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentPatchSuccess: false,
      });
    case PATCH_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentPatchSuccess: true,
        shipmentPatchError: false,
        shipment: action.payload,
      });
    case PATCH_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentPatchSuccess: false,
        shipmentPatchError: null,
        error: action.error.message,
      });
    case ACCEPT_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsAccepting: true,
        shipmentHasAcceptSuccess: false,
      });
    case ACCEPT_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsAccepting: false,
        shipmentHasAcceptSuccess: true,
        shipmentHasAcceptError: false,
        shipment: action.payload,
      });
    case ACCEPT_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsAccepting: false,
        shipmentHasAcceptSuccess: false,
        shipmentHasAcceptError: null,
        error: action.error.message,
      });
    case REJECT_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsRejecting: true,
        shipmentHasRejectSuccess: false,
      });
    case REJECT_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsRejecting: false,
        shipmentHasRejectSuccess: true,
        shipmentHasRejectError: false,
        shipment: action.payload,
      });
    case REJECT_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsRejecting: false,
        shipmentHasRejectSuccess: false,
        shipmentHasRejectError: null,
        error: action.error.message,
      });
    case TRANSPORT_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsSendingTransport: true,
        shipmentHasTransportSuccess: false,
      });
    case TRANSPORT_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsSendingTransport: false,
        shipmentHasTransportSuccess: true,
        shipmentHasTransportError: false,
        shipment: action.payload,
      });
    case TRANSPORT_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsSendingTransport: false,
        shipmentHasTransportSuccess: false,
        shipmentHasTransportError: null,
        error: action.error.message,
      });

    case DELIVER_SHIPMENT.start:
      return Object.assign({}, state, {
        shipmentIsDelivering: true,
        shipmentHasDeliverSuccess: false,
      });
    case DELIVER_SHIPMENT.success:
      return Object.assign({}, state, {
        shipmentIsDelivering: false,
        shipmentHasDeliverSuccess: true,
        shipmentHasDeliverError: false,
        shipment: action.payload,
      });
    case DELIVER_SHIPMENT.failure:
      return Object.assign({}, state, {
        shipmentIsDelivering: false,
        shipmentHasDeliverSuccess: false,
        shipmentHasDeliverError: null,
        error: action.error.message,
      });

    // LOAD SHIPMENT DOCUMENTS
    case LOAD_SHIPMENT_DOCUMENTS.start:
      return Object.assign({}, state, {
        loadingShipmentDocuments: true,
        loadShipmentDocumentsSuccess: false,
      });
    case LOAD_SHIPMENT_DOCUMENTS.success:
      return Object.assign({}, state, {
        loadingShipmentDocuments: false,
        loadShipmentDocumentsSuccess: true,
        loadingShipmentDocumentsError: false,
        shipmentDocuments: action.payload,
      });
    case LOAD_SHIPMENT_DOCUMENTS.failure:
      return Object.assign({}, state, {
        loadingShipmentDocuments: false,
        loadShipmentDocumentsSuccess: false,
        loadingShipmentDocumentsError: true,
        error: action.error.message,
      });

    // SERVICE AGENTS
    case INDEX_SERVICE_AGENTS.start:
      return Object.assign({}, state, {
        serviceAgentsAreLoading: true,
        serviceAgentsHasLoadSucces: false,
      });
    case INDEX_SERVICE_AGENTS.success:
      return Object.assign({}, state, {
        serviceAgentsAreLoading: false,
        serviceAgentsHasLoadSucces: true,
        serviceAgentsHasLoadError: false,
        serviceAgents: action.payload,
      });
    case INDEX_SERVICE_AGENTS.failure:
      return Object.assign({}, state, {
        serviceAgentsAreLoading: false,
        serviceAgentsHasLoadSucces: false,
        serviceAgentsHasLoadError: null,
        serviceAgents: [],
        error: action.error.message,
      });

    case CREATE_SERVICE_AGENTS.start:
      return Object.assign({}, state, {
        serviceAgentIsCreating: true,
        serviceAgentHasCreatedSucces: false,
      });
    case CREATE_SERVICE_AGENTS.success:
      const serviceAgents = state.serviceAgents;
      serviceAgents.push(action.payload);
      return Object.assign({}, state, {
        serviceAgentIsCreating: false,
        serviceAgentHasCreatedSucces: true,
        serviceAgentHasCreatedError: false,
        serviceAgents,
      });
    case CREATE_SERVICE_AGENTS.failure:
      return Object.assign({}, state, {
        serviceAgentIsCreating: false,
        serviceAgentHasCreatedSucces: false,
        serviceAgentHasCreatedError: null,
        serviceAgents: [],
        error: action.error.message,
      });

    case UPDATE_SERVICE_AGENTS.start:
      return Object.assign({}, state, {
        serviceAgentIsUpdating: true,
        serviceAgentHasUpdatedSucces: false,
      });
    case UPDATE_SERVICE_AGENTS.success:
      const updatedAgent = action.payload;
      const updatedAgents = [];
      let extant = false;
      state.serviceAgents.forEach(agent => {
        if (agent.id === updatedAgent.id) {
          extant = true;
          updatedAgents.push(updatedAgent);
        } else {
          updatedAgents.push(agent);
        }
      });
      if (!extant) {
        console.log(
          'WARNING: An updated Agent did not exist before updating: ',
          updatedAgent.id,
        );
        updatedAgents.push(updatedAgent);
      }

      return Object.assign({}, state, {
        serviceAgentIsUpdating: false,
        serviceAgentHasUpdatedSucces: true,
        serviceAgentHasUpdatedError: false,
        serviceAgents: updatedAgents,
      });
    case UPDATE_SERVICE_AGENTS.failure:
      return Object.assign({}, state, {
        serviceAgentIsUpdating: false,
        serviceAgentHasUpdatedSucces: false,
        serviceAgentHasUpdatedError: null,
        serviceAgents: [],
        error: action.error.message,
      });
    // Gov bill of lading
    case GENERATE_GBL.start:
      return Object.assign({}, state, {
        generateGBLSuccess: false,
        generateGBLError: null,
      });
    case GENERATE_GBL.success:
      return Object.assign({}, state, {
        generateGBLSuccess: true,
        generateGBLError: false,
      });
    case GENERATE_GBL.failure:
      return Object.assign({}, state, {
        generateGBLSuccess: false,
        generateGBLError: true,
        error: action.error,
      });

    // MULTIPLE-RESOURCE ACTION TYPES
    //
    // These action types typically dispatch to other actions above to
    // perform their work and exist to encapsulate when multiple requests
    // need to be made in response to a user action.

    // ALL TSP DEPENDENCIES
    case LOAD_TSP_DEPENDENCIES.start:
      return Object.assign({}, state, {
        loadTspDependenciesHasSuccess: false,
        loadTspDependenciesHasError: false,
      });
    case LOAD_TSP_DEPENDENCIES.success:
      return Object.assign({}, state, {
        loadTspDependenciesHasSuccess: true,
        loadTspDependenciesHasError: false,
      });
    case LOAD_TSP_DEPENDENCIES.failure:
      return Object.assign({}, state, {
        loadTspDependenciesHasSuccess: false,
        loadTspDependenciesHasError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}
