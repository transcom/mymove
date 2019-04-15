import { isNull } from 'lodash';
import {
  TransportShipment,
  DeliverShipment,
  CompletePmSurvey,
  IndexServiceAgents,
  GetAllShipmentDocuments,
} from './api.js';

import * as ReduxHelpers from 'shared/ReduxHelpers';
import { getEntitlements } from 'shared/entitlements.js';
import { selectShipment } from 'shared/Entities/modules/shipments';

// SINGLE RESOURCE ACTION TYPES
const transportShipmentType = 'TRANSPORT_SHIPMENT';
const deliverShipmentType = 'DELIVER_SHIPMENT';
const loadShipmentDocumentsType = 'LOAD_SHIPMENT_DOCUMENTS';
const completePmSurveyType = 'COMPLETE_PM_SURVEY';

const indexServiceAgentsType = 'INDEX_SERVICE_AGENTS';

// MULTIPLE-RESOURCE ACTION TYPES
const loadTspDependenciesType = 'LOAD_TSP_DEPENDENCIES';

// SINGLE RESOURCE ACTION TYPES
const TRANSPORT_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(transportShipmentType);
const DELIVER_SHIPMENT = ReduxHelpers.generateAsyncActionTypes(deliverShipmentType);
const COMPLETE_PM_SURVEY = ReduxHelpers.generateAsyncActionTypes(completePmSurveyType);
const LOAD_SHIPMENT_DOCUMENTS = ReduxHelpers.generateAsyncActionTypes(loadShipmentDocumentsType);

const INDEX_SERVICE_AGENTS = ReduxHelpers.generateAsyncActionTypes(indexServiceAgentsType);

// MULTIPLE-RESOURCE ACTION TYPES

const LOAD_TSP_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(loadTspDependenciesType);

// SINGLE-RESOURCE ACTION CREATORS

export const transportShipment = ReduxHelpers.generateAsyncActionCreator(transportShipmentType, TransportShipment);

export const deliverShipment = ReduxHelpers.generateAsyncActionCreator(deliverShipmentType, DeliverShipment);

export const completePmSurvey = ReduxHelpers.generateAsyncActionCreator(completePmSurveyType, CompletePmSurvey);

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
    // SINGLE-RESOURCE ACTION TYPES

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
