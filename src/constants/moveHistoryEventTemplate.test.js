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
} = require('./moveHistoryEventTemplate');

describe('moveHistoryEventTemplate', () => {
  describe('identify Acknowledge excess weight risk event', () => {
    const item = {
      action: 'update',
      eventName: 'acknowledgeExcessWeightRisk',
      tableName: 'moves',
    };

    describe('identify Acknowledge excess weight risk event', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event Acknowledge excess weight risk', () => {
        expect(result).toEqual(acknowledgeExcessWeightRiskEvent);
        expect(result.getDetailsPlainText(item)).toEqual('Dismissed excess weight alert');
      });
    });
  });

  describe('identify Approved shipment event', () => {
    const item = {
      changedValues: { status: 'APPROVED' },
      eventName: 'approveShipment',
      oldValues: { shipment_type: 'HHG' },
      tableName: 'mto_shipments',
    };

    describe('identify Approved shipment event when approveShipment API is used', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event Approved shipment', () => {
        expect(result).toEqual(approveShipmentEvent);
        expect(result.getDetailsPlainText(item)).toEqual('HHG shipment');
      });
    });
  });

  describe('identify Approved shipment diversion event', () => {
    const item = {
      changedValues: { status: 'APPROVED' },
      eventName: 'approveShipmentDiversion',
      oldValues: { shipment_type: 'HHG' },
    };

    describe('identify Approved shipment diversion event when approveShipmentDiversion API is used', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be string Approved shipment', () => {
        expect(result).toEqual(approveShipmentDiversionEvent);
        expect(result.getDetailsPlainText(item)).toEqual('HHG shipment');
      });
    });
  });

  describe('identify Submitted orders event', () => {
    const item = {
      action: 'INSERT',
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'createOrders',
      tableName: 'orders',
    };

    describe('display Submitted orders when createOrders API is used', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event Submitted orders', () => {
        expect(result).toEqual(createOrdersEvent);
      });
    });
  });

  describe('identify Request shipment cancellation event', () => {
    const item = {
      action: 'UPDATE',
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'requestShipmentCancellation',
      oldValues: { shipment_type: 'PPM' },
      tableName: '',
    };

    describe('identify Request shipment cancellation event when requestShipmentCancellation API is used', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event Request shipment cancellation', () => {
        expect(result).toEqual(requestShipmentCancellationEvent);
        expect(result.getDetailsPlainText(item)).toEqual('Requested cancellation for PPM shipment');
      });
    });
  });

  describe('identify Set financial review flag event for flagged move', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'setFinancialReviewFlag',
      changedValues: { financial_review_flag: 'true' },
      tableName: 'moves',
    };

    describe('identify Set financial review flag event when move is flagged', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event Set financial review flag', () => {
        expect(result).toEqual(setFinancialReviewFlagEvent);
        expect(result.getDetailsPlainText(item)).toEqual('Move flagged for financial review');
      });
    });
  });

  describe('identify Set financial review flag event for unflagged move', () => {
    const item = {
      action: 'UPDATE',
      eventName: 'setFinancialReviewFlag',
      changedValues: { financial_review_flag: 'false' },
      tableName: 'moves',
    };

    describe('identify Set financial review flag event when move is unflagged', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event Set financial review flag', () => {
        expect(result).toEqual(setFinancialReviewFlagEvent);
        expect(result.getDetailsPlainText(item)).toEqual('Move unflagged for financial review');
      });
    });
  });

  describe('identify Approved service item', () => {
    const item = {
      changedValues: { status: 'APPROVED' },
      context: { name: 'Domestic origin price', shipment_type: 'HHG_INTO_NTS_DOMESTIC' },
      eventName: 'updateMTOServiceItemStatus',
      tableName: 'mto_service_items',
    };

    describe('identify Approved service item when updateMTOServiceItemStatus API is used', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event Approved service item', () => {
        expect(result).toEqual(updateServiceItemStatusEvent);
        expect(result.getDetailsPlainText(item)).toEqual('NTS shipment, Domestic origin price');
      });
    });
  });

  describe('identify Move approved event', () => {
    const item = {
      action: 'UPDATE',
      changedValues: {
        status: 'APPROVED',
      },
      eventName: 'updateMoveTaskOrderStatus',
      tableName: 'moves',
    };

    describe('identify Move approved event when updateMoveTaskOrderStatus API is used', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event update move task order status', () => {
        expect(result).toEqual(updateMoveTaskOrderStatusEvent);
        expect(result.getDetailsPlainText(item)).toEqual('Created Move Task Order (MTO)');
      });
    });
  });

  describe('identify Move rejected event', () => {
    const item = {
      action: 'UPDATE',
      changedValues: { status: 'REJECTED' },
      eventName: 'updateMoveTaskOrderStatus',
      tableName: 'moves',
    };

    describe('identify Move rejected event when updateMoveTaskOrderStatus API is used', () => {
      const result = getMoveHistoryEventTemplate(item);
      it('should be event update move task order status', () => {
        expect(result).toEqual(updateMoveTaskOrderStatusEvent);
        expect(result.getDetailsPlainText(item)).toEqual('Rejected Move Task Order (MTO)');
      });
    });
  });
});
