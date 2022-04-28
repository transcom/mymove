import moveHistoryOperations from './moveHistoryOperations';
import { shipmentTypes } from './shipments';

function propertiesMatch(p1, p2) {
  return p1 === '*' || p2 === '*' || p1 === p2;
}

export const detailsTypes = {
  PLAIN_TEXT: 'PLAIN_TEXT',
  LABELED: 'LABELED',
  LABELED_SHIPMENT: 'LABELED_SHIPMENT',
  PAYMENT: 'PAYMENT',
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
  getDetailsLabeledDetails = null,
}) => {
  const eventType = {};
  eventType.action = action;
  eventType.eventName = eventName;
  eventType.tableName = tableName;
  eventType.detailsType = detailsType;
  eventType.getEventNameDisplay = getEventNameDisplay;
  eventType.getDetailsPlainText = getDetailsPlainText;
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
  action: 'UPDATE',
  eventName: moveHistoryOperations.acknowledgeExcessWeightRisk,
  tableName: 'moves',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated move',
  getDetailsPlainText: () => 'Dismissed excess weight alert',
});

export const updateBillableWeightEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: moveHistoryOperations.updateBillableWeight,
  tableName: 'entitlements',
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated move',
});

export const updateAllowanceEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: moveHistoryOperations.updateAllowance,
  tableName: 'entitlements',
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated allowances',
});

export const approveShipmentEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: moveHistoryOperations.approveShipment,
  tableName: 'mto_shipments',
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

export const createMTOShipmentEvent = buildMoveHistoryEventTemplate({
  action: 'INSERT',
  eventName: moveHistoryOperations.createMTOShipment,
  tableName: 'mto_shipments',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Submitted/Requested shipments',
  getDetailsPlainText: () => '-',
});

export const createOrdersEvent = buildMoveHistoryEventTemplate({
  action: 'INSERT',
  eventName: moveHistoryOperations.createOrders,
  tableName: 'orders',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Submitted orders',
  getDetailsPlainText: () => '-',
});

export const createBasicServiceItemEvent = buildMoveHistoryEventTemplate({
  action: 'INSERT',
  eventName: moveHistoryOperations.updateMoveTaskOrderStatus,
  tableName: 'mto_service_items',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved service item',
  getDetailsPlainText: (historyRecord) => {
    return `${historyRecord.context[0]?.name}`;
  },
});

export const createStandardServiceItemEvent = buildMoveHistoryEventTemplate({
  action: 'INSERT',
  eventName: moveHistoryOperations.approveShipment,
  tableName: 'mto_service_items',
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
  action: 'INSERT',
  eventName: moveHistoryOperations.requestShipmentReweigh,
  tableName: 'reweighs',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsPlainText: (historyRecord) => {
    return `${shipmentTypes[historyRecord.context[0]?.shipment_type]} shipment, reweigh requested`;
  },
});

export const setFinancialReviewFlagEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: moveHistoryOperations.setFinancialReviewFlag,
  tableName: 'moves',
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

export const updateMoveTaskOrderEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: moveHistoryOperations.updateMoveTaskOrder,
  tableName: 'moves',
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated move',
});

export const updateMoveTaskOrderStatusEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: moveHistoryOperations.updateMoveTaskOrderStatus,
  tableName: 'moves',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: (historyRecord) => {
    return historyRecord.changedValues?.status === 'APPROVED' ? 'Approved move' : 'Rejected move';
  },
  getDetailsPlainText: (historyRecord) => {
    return historyRecord.changedValues?.status === 'APPROVED'
      ? 'Created Move Task Order (MTO)'
      : 'Rejected Move Task Order (MTO)';
  },
});

export const updateServiceItemStatusEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: moveHistoryOperations.updateServiceItemStatus,
  tableName: 'mto_service_items',
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
  action: 'UPDATE',
  eventName: moveHistoryOperations.uploadAmendedOrders,
  tableName: 'orders',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated orders',
  getDetailsPlainText: () => '-',
});

export const updateOrderEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: '*',
  tableName: 'orders',
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated orders',
  getDetailsLabeledDetails: ({ changedValues, context }) => {
    let newChangedValues;

    if (context) {
      newChangedValues = {
        ...changedValues,
        ...context[0],
      };
    } else {
      newChangedValues = changedValues;
    }

    // merge context with change values for only this event
    return newChangedValues;
  },
});

export const updateMTOReviewdBillableWeightsAtEvent = buildMoveHistoryEventTemplate({
  action: 'UPDATE',
  eventName: moveHistoryOperations.UpdateMTOReviewedBillableWeightsAt,
  tableName: 'moves',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsPlainText: (historyRecord) => historyRecord.changedValues?.billable_weights_reviewed_at,
});

export const undefinedEvent = buildMoveHistoryEventTemplate({
  action: '*',
  eventName: '*',
  tableName: '*',
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: () => {
    return 'Undefined event type';
  },
  getDetailsPlainText: () => {
    return 'Undefined event details';
  },
});

const allMoveHistoryEventTemplates = [
  acknowledgeExcessWeightRiskEvent,
  approveShipmentEvent,
  approveShipmentDiversionEvent,
  createMTOShipmentEvent,
  createOrdersEvent,
  createBasicServiceItemEvent,
  createStandardServiceItemEvent,
  requestShipmentCancellationEvent,
  requestShipmentDiversionEvent,
  requestShipmentReweighEvent,
  setFinancialReviewFlagEvent,
  submitMoveForApprovalEvent,
  updateMoveTaskOrderEvent,
  updateMoveTaskOrderStatusEvent,
  updateOrderEvent,
  updateServiceItemStatusEvent,
  uploadAmendedOrdersEvent,
  updateBillableWeightEvent,
  updateAllowanceEvent,
  updateMTOReviewdBillableWeightsAtEvent,
];

const getMoveHistoryEventTemplate = (historyRecord) => {
  return allMoveHistoryEventTemplates.find((eventType) => eventType.matches(historyRecord)) || undefinedEvent;
};

export default getMoveHistoryEventTemplate;
