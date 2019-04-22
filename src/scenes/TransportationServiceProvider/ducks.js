import { isNull } from 'lodash';
import { IndexServiceAgents, GetAllShipmentDocuments } from './api.js';

import * as ReduxHelpers from 'shared/ReduxHelpers';
import { getEntitlements } from 'shared/entitlements.js';
import { selectShipment } from 'shared/Entities/modules/shipments';

// SINGLE RESOURCE ACTION TYPES
const loadShipmentDocumentsType = 'LOAD_SHIPMENT_DOCUMENTS';

const indexServiceAgentsType = 'INDEX_SERVICE_AGENTS';

// MULTIPLE-RESOURCE ACTION TYPES
const loadTspDependenciesType = 'LOAD_TSP_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES
const LOAD_SHIPMENT_DOCUMENTS = ReduxHelpers.generateAsyncActionTypes(loadShipmentDocumentsType);

const INDEX_SERVICE_AGENTS = ReduxHelpers.generateAsyncActionTypes(indexServiceAgentsType);

// MULTIPLE-RESOURCE ACTION TYPES

const LOAD_TSP_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(loadTspDependenciesType);

// SINGLE-RESOURCE ACTION CREATORS

export const getAllShipmentDocuments = ReduxHelpers.generateAsyncActionCreator(
  loadShipmentDocumentsType,
  GetAllShipmentDocuments,
);

export const indexServiceAgents = ReduxHelpers.generateAsyncActionCreator(indexServiceAgentsType, IndexServiceAgents);

export function loadEntitlements(state, shipmentId) {
  const shipment = selectShipment(state, shipmentId);
  const move = shipment.move || {};
  const serviceMember = shipment.service_member || {};
  const hasDependents = move.has_dependents;
  const spouseHasProGear = move.spouse_has_progear;
  const rank = serviceMember.rank;

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
  shipmentIsRejecting: false,
  shipmentHasRejectError: null,
  shipmentHasRejectSuccess: false,
  shipmentIsDelivering: false,
  shipmentHasDeliverError: null,
  shipmentHasDeliverSuccess: false,
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
};

export function tspReducer(state = initialState, action) {
  switch (action.type) {
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
