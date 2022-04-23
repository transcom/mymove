const {
  default: getMoveHistoryEventTemplate,
  approveShipmentEvent,
  approveShipmentDiversionEvent,
  createOrdersEvent,
  requestShipmentCancellationEvent,
  setFinancialReviewFlagEvent,
  updateMoveTaskOrderStatusEvent,
  updateServiceItemStatusEvent,
  acknowledgeExcessWeightRiskEvent,
  createStandardServiceItemEvent,
  createBasicServiceItemEvent,
} = require('./moveHistoryEventTemplate');

describe('moveHistoryEventTemplate', () => {
  describe('when given an Acknowledge excess weight risk history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'acknowledgeExcessWeightRisk',
      tableName: 'moves',
    };
    it('correctly matches the Acknowledge excess weight risk event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(acknowledgeExcessWeightRiskEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Dismissed excess weight alert');
    });
  });

  describe('when given an Approved shipment history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: { status: 'APPROVED' },
      eventName: 'approveShipment',
      oldValues: { shipment_type: 'HHG' },
      tableName: 'mto_shipments',
    };
    it('correctly matches the Approved shipment event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(approveShipmentEvent);
      expect(result.getDetailsPlainText(item)).toEqual('HHG shipment');
    });
  });

  describe('when given an Approved shipment diversion history record', () => {
    const item = {
      changedValues: { status: 'APPROVED' },
      eventName: 'approveShipmentDiversion',
      oldValues: { shipment_type: 'HHG' },
    };
    it('correctly matches the Approved shipment event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(approveShipmentDiversionEvent);
      expect(result.getDetailsPlainText(item)).toEqual('HHG shipment');
    });
  });

  describe('when given a Submitted orders history record', () => {
    const item = {
      action: 'INSERT',
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'createOrders',
      tableName: 'orders',
    };
    it('correctly matches the Submitted orders event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createOrdersEvent);
    });
  });

  describe('when given a Create basic service item history record', () => {
    const item = {
      action: 'INSERT',
      context: {
        name: 'Domestic linehaul',
      },
      eventName: 'updateMoveTaskOrderStatus',
      tableName: 'mto_service_items',
    };
    it('correctly matches the Create basic service item event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createBasicServiceItemEvent);
      expect(result.getEventNameDisplay(result)).toEqual('Approved service item');
      expect(result.getDetailsPlainText(item)).toEqual('Domestic linehaul');
    });
  });

  describe('when given a Create standard service item history record', () => {
    const item = {
      action: 'INSERT',
      context: [
        {
          shipment_type: 'HHG',
          name: 'Domestic linehaul',
        },
      ],
      eventName: 'approveShipment',
      tableName: 'mto_service_items',
    };
    it('correctly matches the Create standard service item event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(createStandardServiceItemEvent);
      expect(result.getEventNameDisplay(result)).toEqual('Approved service item');
      expect(result.getDetailsPlainText(item)).toEqual('HHG shipment, Domestic linehaul');
    });
  });

  describe('when given a Request shipment cancellation history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'requestShipmentCancellation',
      oldValues: { shipment_type: 'PPM' },
      tableName: '',
    };
    it('correctly matches the Request shipment cancellation event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(requestShipmentCancellationEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Requested cancellation for PPM shipment');
    });
  });

  describe('when given a Set financial review flag event for flagged move history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'setFinancialReviewFlag',
      changedValues: { financial_review_flag: 'true' },
      tableName: 'moves',
    };
    it('correctly matches the Set financial review flag event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(setFinancialReviewFlagEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Move flagged for financial review');
    });
  });

  describe('when given a Set financial review flag event for unflagged move history record', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'setFinancialReviewFlag',
      changedValues: { financial_review_flag: 'false' },
      tableName: 'moves',
    };
    it('correctly matches the Set financial review flag event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(setFinancialReviewFlagEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Move unflagged for financial review');
    });
  });

  describe('when given an Approved service item history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: { status: 'APPROVED' },
      context: [{ name: 'Domestic origin price', shipment_type: 'HHG_INTO_NTS_DOMESTIC' }],
      eventName: 'updateMTOServiceItemStatus',
      tableName: 'mto_service_items',
    };
    it('correctly matches the Approved service item event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateServiceItemStatusEvent);
      expect(result.getDetailsPlainText(item)).toEqual('NTS shipment, Domestic origin price');
    });
  });

  describe('when given a Move approved history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: {
        status: 'APPROVED',
      },
      eventName: 'updateMoveTaskOrderStatus',
      tableName: 'moves',
    };
    it('correctly matches the Update move task order status event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMoveTaskOrderStatusEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Created Move Task Order (MTO)');
    });
  });

  describe('when given a Move rejected history record', () => {
    const item = {
      action: 'UPDATE',
      changedValues: { status: 'REJECTED' },
      eventName: 'updateMoveTaskOrderStatus',
      tableName: 'moves',
    };
    it('correctly matches the Update move task order status event', () => {
      const result = getMoveHistoryEventTemplate(item);
      expect(result).toEqual(updateMoveTaskOrderStatusEvent);
      expect(result.getDetailsPlainText(item)).toEqual('Rejected Move Task Order (MTO)');
    });
  });
});
