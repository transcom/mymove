import { isNull, get } from 'lodash';
import {
  LoadShipment,
  PatchShipment,
  AcceptShipment,
  TransportShipment,
  DeliverShipment,
  CompletePmSurvey,
  CreateServiceAgent,
  IndexServiceAgents,
  UpdateServiceAgent,
  GetAllShipmentDocuments,
} from './api.js';

import * as ReduxHelpers from 'shared/ReduxHelpers';
import { getEntitlements } from 'shared/entitlements.js';

// SINGLE RESOURCE ACTION TYPES
const loadShipmentType = 'LOAD_SHIPMENT';
const patchShipmentType = 'PATCH_SHIPMENT';
const acceptShipmentType = 'ACCEPT_SHIPMENT';
const transportShipmentType = 'TRANSPORT_SHIPMENT';
const deliverShipmentType = 'DELIVER_SHIPMENT';
const loadShipmentDocumentsType = 'LOAD_SHIPMENT_DOCUMENTS';
const completePmSurveyType = 'COMPLETE_PM_SURVEY';

const indexServiceAgentsType = 'INDEX_SERVICE_AGENTS';
const createServiceAgentsType = 'CREATE_SERVICE_AGENTS';
const updateServiceAgentsType = 'UPDATE_SERVICE_AGENTS';

// MULTIPLE-RESOURCE ACTION TYPES
const loadTspDependenciesType = 'LOAD_TSP_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES
const LOAD_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(loadShipmentType);
const PATCH_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(patchShipmentType);
const ACCEPT_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(acceptShipmentType);
const TRANSPORT_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(transportShipmentType);
const DELIVER_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(deliverShipmentType);
const COMPLETE_PM_SURVEY = ReduxHelpers.generateAsyncActionTypes(completePmSurveyType);
const LOAD_SHIPMENT_DOCUMENTS = ReduxHelpers.generateAsyncActionTypes(loadShipmentDocumentsType);

const INDEX_SERVICE_AGENTS = ReduxHelpers.generateAsyncActionTypes(indexServiceAgentsType);

const CREATE_SERVICE_AGENTS = ReduxHelpers.generateAsyncActionTypes(createServiceAgentsType);

const UPDATE_SERVICE_AGENTS = ReduxHelpers.generateAsyncActionTypes(updateServiceAgentsType);

// MULTIPLE-RESOURCE ACTION TYPES

const LOAD_TSP_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(loadTspDependenciesType);

// SINGLE-RESOURCE ACTION CREATORS

export const loadShipment = ReduxHelpers.generateAsyncActionCreator(loadShipmentType, LoadShipment);

export const patchShipment = ReduxHelpers.generateAsyncActionCreator(patchShipmentType, PatchShipment);

export const acceptShipment = ReduxHelpers.generateAsyncActionCreator(acceptShipmentType, AcceptShipment);

export const transportShipment = ReduxHelpers.generateAsyncActionCreator(transportShipmentType, TransportShipment);

export const deliverShipment = ReduxHelpers.generateAsyncActionCreator(deliverShipmentType, DeliverShipment);

export const completePmSurvey = ReduxHelpers.generateAsyncActionCreator(completePmSurveyType, CompletePmSurvey);

export const getAllShipmentDocuments = ReduxHelpers.generateAsyncActionCreator(
  loadShipmentDocumentsType,
  GetAllShipmentDocuments,
);

export const indexServiceAgents = ReduxHelpers.generateAsyncActionCreator(indexServiceAgentsType, IndexServiceAgents);

export const createServiceAgent = ReduxHelpers.generateAsyncActionCreator(createServiceAgentsType, CreateServiceAgent);

export const updateServiceAgent = ReduxHelpers.generateAsyncActionCreator(updateServiceAgentsType, UpdateServiceAgent);

// MULTIPLE-RESOURCE ACTION CREATORS
//
// These action types typically dispatch to other actions above to
// perform their work and exist to encapsulate when multiple requests
// need to be made in response to a user action.

export function handleServiceAgents(shipmentId, serviceAgents) {
  return async function(dispatch, getState) {
    for (const serviceAgent in serviceAgents) {
      // eslint-disable-next-line security/detect-object-injection
      dispatch(createOrUpdateServiceAgent(shipmentId, serviceAgents[serviceAgent]));
    }
  };
}

export function createOrUpdateServiceAgent(shipmentId, serviceAgent) {
  return async function(dispatch, getState) {
    if (serviceAgent.id) {
      return dispatch(updateServiceAgent(serviceAgent));
    } else if (!serviceAgent.company || !serviceAgent.email || !serviceAgent.phone_number) {
      // Don't send the service agent if it's not got enough details
      // Currently, it should only be the destination agent that gets skipped
      return;
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
      await Promise.all([dispatch(loadShipment(shipmentId)), dispatch(indexServiceAgents(shipmentId))]);
      return dispatch(actions.success());
    } catch (ex) {
      dispatch(actions.error(ex));
      throw new Error(ex);
    }
  };
}

// Selectors
export function loadEntitlements(state) {
  const hasDependents = get(state, 'tsp.shipment.move.has_dependents', null);
  const spouseHasProGear = get(state, 'tsp.shipment.move.spouse_has_pro_gear', null);
  const rank = get(state, 'tsp.shipment.service_member.rank', null);
  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents, spouseHasProGear);
}
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
  storageInTransitIsCreating: false,
  storageInTransitHasCreatedSuccess: false,
  storageInTransitHasCreatedError: null,
  storageInTransitsAreLoading: false,
  storageInTransitsHasLoadSuccess: false,
  storageInTransitsHasLoadError: null,
  shipment: {},
  loadTspDependenciesHasSuccess: false,
  loadTspDependenciesHasError: null,
  flashMessage: false,
  serviceAgents: [],
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

    // PM SURVEY ACTION
    case COMPLETE_PM_SURVEY.start:
      return Object.assign({}, state, {
        pmSurveyIsCompleting: true,
        pmSurveyHasCompletionSuccess: false,
      });
    case COMPLETE_PM_SURVEY.success:
      return Object.assign({}, state, {
        pmSurveyIsCompleting: false,
        pmSurveyHasCompletionSuccess: true,
        pmSurveyHasCompletionError: false,
        shipment: action.payload,
      });
    case COMPLETE_PM_SURVEY.failure:
      return Object.assign({}, state, {
        pmSurveyIsCompleting: false,
        pmSurveyHasCompletionSuccess: false,
        pmSurveyHasCompletionError: null,
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
        console.log('WARNING: An updated Agent did not exist before updating: ', updatedAgent.id);
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
