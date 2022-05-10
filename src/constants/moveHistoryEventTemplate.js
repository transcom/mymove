import moveHistoryOperations from './moveHistoryOperations';
import { shipmentTypes } from './shipments';

import { formatMoveHistoryFullAddress, formatMoveHistoryAgent } from 'utils/formatters';
import { dbActions, dbTables } from 'constants/historyLogUIDisplayName';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

function propertiesMatch(p1, p2) {
  return p1 === '*' || p2 === '*' || p1 === p2;
}

export const detailsTypes = {
  PLAIN_TEXT: 'PLAIN_TEXT',
  LABELED: 'LABELED',
  LABELED_SHIPMENT: 'LABELED_SHIPMENT',
  PAYMENT: 'PAYMENT',
  STATUS: 'STATUS',
};

const buildMoveHistoryEventTemplate = ({
  action = '*',
  eventName = '*',
  tableName = '*',
  detailsType = detailsTypes.PLAIN_TEXT,
  getEventNameDisplay = () => {
    return 'Undefined event type';
  },
  getDetailsPlainText = () => {
    return 'Undefined details';
  },
  getStatusDetails = () => {
    return 'Undefined status';
  },
  getDetailsLabeledDetails = null,
}) => {
  const eventType = {};
  eventType.action = action;
  eventType.eventName = eventName;
  eventType.tableName = tableName;
  eventType.detailsType = detailsType;
  eventType.getEventNameDisplay = getEventNameDisplay;
  eventType.getDetailsPlainText = getDetailsPlainText;
  eventType.getStatusDetails = getStatusDetails;
  eventType.getDetailsLabeledDetails = getDetailsLabeledDetails;

  eventType.matches = (other) => {
    if (eventType === undefined || other === undefined) {
      return false;
    }
    return (
      propertiesMatch(eventType.action, other?.action) &&
      propertiesMatch(eventType.eventName, other?.eventName) &&
      propertiesMatch(eventType.tableName, other?.tableName)
    );
  };

  return eventType;
};

export const acknowledgeExcessWeightRiskEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.acknowledgeExcessWeightRisk,
  tableName: dbTables.moves,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated move',
  getDetailsPlainText: () => 'Dismissed excess weight alert',
});

export const approveShipmentEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.approveShipment,
  tableName: dbTables.mto_shipments,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved shipment',
  getDetailsPlainText: (historyRecord) => {
    return `${shipmentTypes[historyRecord.oldValues?.shipment_type]} shipment`;
  },
});

export const approveShipmentDiversionEvent = buildMoveHistoryEventTemplate({
  action: '*',
  eventName: moveHistoryOperations.approveShipmentDiversion,
  tableName: '*',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved diversion',
  getDetailsPlainText: (historyRecord) => {
    return `${shipmentTypes[historyRecord.oldValues?.shipment_type]} shipment`;
  },
});

export const createBasicServiceItemEvent = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: moveHistoryOperations.updateMoveTaskOrderStatus,
  tableName: dbTables.mto_service_items,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved service item',
  getDetailsPlainText: (historyRecord) => {
    return `${historyRecord.context[0]?.name}`;
  },
});

export const createMTOShipmentEvent = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: moveHistoryOperations.createMTOShipment,
  tableName: dbTables.mto_shipments,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Submitted/Requested shipments',
  getDetailsPlainText: () => '-',
});

export const createMTOShipmentAddressesEvent = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: '',
  tableName: dbTables.addresses,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ changedValues, context }) => {
    const address = formatMoveHistoryFullAddress(changedValues);

    const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;

    let addressLabel = '';
    if (addressType === 'pickupAddress') {
      addressLabel = 'pickup_address';
    } else if (addressType === 'destinationAddress') {
      addressLabel = 'destination_address';
    }

    const newChangedValues = {
      shipment_type: context[0]?.shipment_type,
      ...changedValues,
    };

    newChangedValues[addressLabel] = address;

    return newChangedValues;
  },
});

export const createMTOShipmentAgentEvent = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: moveHistoryOperations.updateMTOShipment,
  tableName: dbTables.mto_agents,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ changedValues, oldValues, context }) => {
    const agent = formatMoveHistoryAgent(changedValues);

    const agentType = changedValues.agent_type ?? oldValues.agent_type;

    let agentLabel = '';
    if (agentType === 'RECEIVING_AGENT') {
      agentLabel = 'receiving_agent';
    } else if (agentType === 'RELEASING_AGENT') {
      agentLabel = 'releasing_agent';
    }

    const newChangedValues = {
      shipment_type: context[0]?.shipment_type,
      ...changedValues,
    };

    newChangedValues[agentLabel] = agent;

    return newChangedValues;
  },
});

export const createOrdersEvent = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: moveHistoryOperations.createOrders,
  tableName: dbTables.orders,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Submitted orders',
  getDetailsPlainText: () => '-',
});

export const createPaymentRequestReweighUpdate = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: moveHistoryOperations.updateReweigh,
  tableName: dbTables.payment_requests,
  detailsType: detailsTypes.STATUS,
  getEventNameDisplay: () => 'Created payment request',
  getStatusDetails: () => {
    return 'Pending';
  },
});

export const createPaymentRequestShipmentUpdate = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: moveHistoryOperations.updateMTOShipment,
  tableName: dbTables.payment_requests,
  detailsType: detailsTypes.STATUS,
  getEventNameDisplay: () => 'Created payment request',
  getStatusDetails: () => {
    return 'Pending';
  },
});

export const createStandardServiceItemEvent = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: moveHistoryOperations.approveShipment,
  tableName: dbTables.mto_service_items,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved service item',
  getDetailsPlainText: (historyRecord) => {
    return `${shipmentTypes[historyRecord.context[0]?.shipment_type]} shipment, ${historyRecord.context[0]?.name}`;
  },
});

export const requestShipmentCancellationEvent = buildMoveHistoryEventTemplate({
  action: '*',
  eventName: moveHistoryOperations.requestShipmentCancellation,
  tableName: '*',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsPlainText: (historyRecord) => {
    return `Requested cancellation for ${shipmentTypes[historyRecord.oldValues?.shipment_type]} shipment`;
  },
});

export const requestShipmentDiversionEvent = buildMoveHistoryEventTemplate({
  action: '*',
  eventName: moveHistoryOperations.requestShipmentDiversion,
  tableName: '*',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Requested diversion',
  getDetailsPlainText: (historyRecord) => {
    return `Requested diversion for ${shipmentTypes[historyRecord.oldValues?.shipment_type]} shipment`;
  },
});

export const requestShipmentReweighEvent = buildMoveHistoryEventTemplate({
  action: dbActions.INSERT,
  eventName: moveHistoryOperations.requestShipmentReweigh,
  tableName: 'reweighs',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsPlainText: (historyRecord) => {
    return `${shipmentTypes[historyRecord.context[0]?.shipment_type]} shipment, reweigh requested`;
  },
});

export const setFinancialReviewFlagEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.setFinancialReviewFlag,
  tableName: dbTables.moves,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => {
    return 'Flagged move';
  },
  getDetailsPlainText: (historyRecord) => {
    return historyRecord.changedValues?.financial_review_flag === 'true'
      ? 'Move flagged for financial review'
      : 'Move unflagged for financial review';
  },
});

export const submitMoveForApprovalEvent = buildMoveHistoryEventTemplate({
  action: '*',
  eventName: moveHistoryOperations.submitMoveForApproval,
  tableName: '*',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Submitted move',
  getDetailsPlainText: () => '-',
});

export const updateAllowanceEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateAllowance,
  tableName: dbTables.entitlements,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated allowances',
});

export const updateBillableWeightEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateBillableWeight,
  tableName: dbTables.entitlements,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated move',
});

export const updateMoveTaskOrderEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateMoveTaskOrder,
  tableName: dbTables.moves,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated move',
});

export const updateMoveTaskOrderStatusEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateMoveTaskOrderStatus,
  tableName: dbTables.moves,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: (historyRecord) => {
    return historyRecord.changedValues?.available_to_prime_at ? 'Approved move' : 'Move status updated';
  },
  getDetailsPlainText: (historyRecord) => {
    return historyRecord.changedValues?.available_to_prime_at ? 'Created Move Task Order (MTO)' : '-';
  },
});

export const updateMTOShipmentEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateMTOShipment,
  tableName: dbTables.mto_shipments,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: (historyRecord) => {
    return {
      shipment_type: historyRecord.oldValues.shipment_type,
      ...historyRecord.changedValues,
    };
  },
});

export const updateMTOShipmentAddressesEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateMTOShipment,
  tableName: dbTables.addresses,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ oldValues, changedValues, context }) => {
    let newChangedValues = {
      street_address_1: oldValues.street_address_1,
      street_address_2: oldValues.street_address_2,
      city: oldValues.city,
      state: oldValues.state,
      postal_code: oldValues.postal_code,
      ...changedValues,
    };

    const address = formatMoveHistoryFullAddress(newChangedValues);

    const addressType = context.filter((contextObject) => contextObject.address_type)[0].address_type;

    let addressLabel = '';
    if (addressType === 'pickupAddress') {
      addressLabel = 'pickup_address';
    } else if (addressType === 'destinationAddress') {
      addressLabel = 'destination_address';
    }

    newChangedValues = {
      shipment_type: context[0]?.shipment_type,
      ...changedValues,
    };

    newChangedValues[addressLabel] = address;

    return newChangedValues;
  },
});

export const updateMTOShipmentAgentEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateMTOShipment,
  tableName: dbTables.mto_agents,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ oldValues, changedValues, context }) => {
    let newChangedValues = {
      email: oldValues.email,
      first_name: oldValues.first_name,
      last_name: oldValues.last_name,
      phone: oldValues.phone,
      ...changedValues,
    };

    const agent = formatMoveHistoryAgent(newChangedValues);

    const agentType = changedValues.agent_type ?? oldValues.agent_type;

    let agentLabel = '';
    if (agentType === 'RECEIVING_AGENT') {
      agentLabel = 'receiving_agent';
    } else if (agentType === 'RELEASING_AGENT') {
      agentLabel = 'releasing_agent';
    }

    newChangedValues = {
      shipment_type: context[0].shipment_type,
      ...changedValues,
    };

    newChangedValues[agentLabel] = agent;

    return newChangedValues;
  },
});

export const updateMTOShipmentDeprecatePaymentRequest = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateMTOShipment,
  tableName: dbTables.payment_requests,
  detailsType: detailsTypes.STATUS,
  getEventNameDisplay: ({ oldValues }) => `Updated payment request ${oldValues?.payment_request_number}`,
  getStatusDetails: ({ changedValues }) => {
    const { status } = changedValues;
    return PAYMENT_REQUEST_STATUS_LABELS[status];
  },
});

export const updatePaymentRequestStatus = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updatePaymentRequestStatus,
  tableName: dbTables.payment_requests,
  detailsType: detailsTypes.PAYMENT,
  getEventNameDisplay: () => 'Submitted payment request',
});

export const updateServiceItemStatusEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateServiceItemStatus,
  tableName: dbTables.mto_service_items,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: (historyRecord) => {
    switch (historyRecord.changedValues?.status) {
      case 'APPROVED':
        return 'Approved service item';
      case 'REJECTED':
        return 'Rejected service item';
      default:
        return '';
    }
  },
  getDetailsPlainText: (historyRecord) => {
    return `${shipmentTypes[historyRecord.context[0]?.shipment_type]} shipment, ${historyRecord.context[0]?.name}`;
  },
});

export const uploadAmendedOrdersEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.uploadAmendedOrders,
  tableName: 'orders',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated orders',
  getDetailsPlainText: () => '-',
});

export const updateOrderEvent = buildMoveHistoryEventTemplate({
  action: dbActions.UPDATE,
  eventName: '*',
  tableName: dbTables.orders,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated orders',
  getDetailsLabeledDetails: (historyRecord) => {
    let newChangedValues;

    if (historyRecord.context) {
      newChangedValues = {
        ...historyRecord.changedValues,
        ...historyRecord.context[0],
      };
    } else {
      newChangedValues = historyRecord.changedValues;
    }

    // merge context with change values for only this event
    return newChangedValues;
  },
});

export const undefinedEvent = buildMoveHistoryEventTemplate({
  action: null,
  eventName: null,
  tableName: null,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: ({ tableName }) => {
    switch (tableName) {
      case dbTables.orders:
        return 'Updated order';
      case dbTables.mto_service_items:
        return 'Updated service item';
      case dbTables.entitlements:
        return 'Updated allowances';
      case dbTables.payment_requests:
        return 'Updated payment request';
      case dbTables.mto_shipments:
      case dbTables.mto_agents:
      case dbTables.addresses:
        return 'Updated shipment';
      case dbTables.moves:
      default:
        return 'Updated move';
    }
  },
  getDetailsPlainText: () => {
    return '-';
  },
});

const allMoveHistoryEventTemplates = [
  acknowledgeExcessWeightRiskEvent,
  approveShipmentEvent,
  approveShipmentDiversionEvent,
  createMTOShipmentEvent,
  createMTOShipmentAddressesEvent,
  createMTOShipmentAgentEvent,
  createOrdersEvent,
  createPaymentRequestReweighUpdate,
  createPaymentRequestShipmentUpdate,
  createBasicServiceItemEvent,
  createStandardServiceItemEvent,
  requestShipmentCancellationEvent,
  requestShipmentDiversionEvent,
  requestShipmentReweighEvent,
  setFinancialReviewFlagEvent,
  submitMoveForApprovalEvent,
  updateAllowanceEvent,
  uploadAmendedOrdersEvent,
  updateBillableWeightEvent,
  updateMoveTaskOrderEvent,
  updateMoveTaskOrderStatusEvent,
  updateMTOShipmentEvent,
  updateMTOShipmentAddressesEvent,
  updateMTOShipmentAgentEvent,
  updateMTOShipmentDeprecatePaymentRequest,
  updateOrderEvent,
  updatePaymentRequestStatus,
  updateServiceItemStatusEvent,
  updateBillableWeightEvent,
  updateAllowanceEvent,
];

const getMoveHistoryEventTemplate = (historyRecord) => {
  return allMoveHistoryEventTemplates.find((eventType) => eventType.matches(historyRecord)) || undefinedEvent;
};

export default getMoveHistoryEventTemplate;
