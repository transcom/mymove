import {
  getHistoryLogEventNameDisplay,
  eventNamePlainTextToDisplay,
  shipmentOptionToDisplay,
} from './historyLogUIDisplayName';

describe('historyLogUIDisplay', () => {
  describe('display Submitted orders', () => {
    const item = {
      changedValues: {
        status: 'DRAFT',
      },
      eventName: 'createOrders',
    };

    // ['createOrders', 'Submitted orders'], //internal.yaml
    describe('display Submitted orders when createOrders API is used', () => {
      const result = getHistoryLogEventNameDisplay(item);
      it('should be string Submitted orders', () => {
        expect(result).toEqual('Submitted orders');
      });
    });
  });

  describe('display Approved service item', () => {
    const item = {
      changedValues: {
        status: 'APPROVED',
      },
      eventName: 'updateMTOServiceItemStatus',
    };

    // ['createOrders', 'Submitted orders'], //internal.yaml
    describe('display Approved service item when updateMTOServiceItemStatus API is used', () => {
      const result = getHistoryLogEventNameDisplay(item);
      it('should be string Approved service item', () => {
        expect(result).toEqual('Approved service item');
      });
    });
  });

  describe('display Approved shipment', () => {
    const item = {
      changedValues: {
        status: 'APPROVED',
      },
      eventName: 'approveShipmentDiversion',
    };

    // ['createOrders', 'Submitted orders'], //internal.yaml
    describe('display Approved shipment when approveShipmentDiversion API is used', () => {
      const result = getHistoryLogEventNameDisplay(item);
      it('should be string Approved shipment', () => {
        expect(result).toEqual('Approved shipment');
      });
    });
  });

  describe('display Approved shipment', () => {
    const item = {
      changedValues: {
        status: 'APPROVED',
      },
      eventName: 'approveShipment',
    };

    // ['createOrders', 'Submitted orders'], //internal.yaml
    describe('display Approved shipment when approveShipment API is used', () => {
      const result = getHistoryLogEventNameDisplay(item);
      it('should be string Approved shipment', () => {
        expect(result).toEqual('Approved shipment');
      });
    });
  });

  describe('display Move rejected', () => {
    const item = {
      changedValues: {
        status: 'REJECTED',
      },
      eventName: 'updateMoveTaskOrderStatus',
    };

    // ['createOrders', 'Submitted orders'], //internal.yaml
    describe('display Move rejected when updateMoveTaskOrderStatus API is used', () => {
      const result = getHistoryLogEventNameDisplay(item);
      it('should be string Approved service item', () => {
        expect(result).toEqual('Move rejected');
      });
    });
  });

  describe('displays the correct plain text when each eventNamePlainTextToDisplay is called', () => {
    it.each([
      ['approveShipment', 'HHG shipment', { shipment_type: 'HHG' }, { status: 'APPROVED' }],
      [
        'requestShipmentDiversion',
        `Requested diversion for ${shipmentOptionToDisplay.HHG} shipment`,
        { shipment_type: 'HHG' },
        {},
      ],
      [
        'requestShipmentDiversion',
        `Requested diversion for ${shipmentOptionToDisplay.HHG_OUTOF_NTS_DOMESTIC} shipment`,
        { shipment_type: 'HHG_OUTOF_NTS_DOMESTIC' },
        {},
      ],
      [
        'requestShipmentDiversion',
        `Requested diversion for ${shipmentOptionToDisplay.HHG_INTO_NTS_DOMESTIC} shipment`,
        { shipment_type: 'HHG_INTO_NTS_DOMESTIC' },
        {},
      ],
      [
        'requestShipmentDiversion',
        `Requested diversion for ${shipmentOptionToDisplay.PPM} shipment`,
        { shipment_type: 'PPM' },
        {},
      ],
      [
        'requestShipmentDiversion',
        `Requested diversion for ${shipmentOptionToDisplay.HHG_SHORTHAUL_DOMESTIC} shipment`,
        { shipment_type: 'HHG_SHORTHAUL_DOMESTIC' },
        {},
      ],
      ['updateMTOServiceItemStatus', 'Service item status', {}, {}],
      ['setFinancialReviewFlag', 'Move flagged for financial review', {}, { financial_review_flag: 'true' }],
      ['setFinancialReviewFlag', 'Move unflagged for financial review', {}, { financial_review_flag: 'false' }],
      [
        'requestShipmentCancellation',
        `Requested cancellation for ${shipmentOptionToDisplay.HHG} shipment`,
        { shipment_type: 'HHG' },
        {},
      ],
      [
        'requestShipmentCancellation',
        `Requested cancellation for ${shipmentOptionToDisplay.HHG_OUTOF_NTS_DOMESTIC} shipment`,
        { shipment_type: 'HHG_OUTOF_NTS_DOMESTIC' },
        {},
      ],
      [
        'requestShipmentCancellation',
        `Requested cancellation for ${shipmentOptionToDisplay.HHG_INTO_NTS_DOMESTIC} shipment`,
        { shipment_type: 'HHG_INTO_NTS_DOMESTIC' },
        {},
      ],
      [
        'requestShipmentCancellation',
        `Requested cancellation for ${shipmentOptionToDisplay.PPM} shipment`,
        { shipment_type: 'PPM' },
        {},
      ],
      [
        'requestShipmentCancellation',
        `Requested cancellation for ${shipmentOptionToDisplay.HHG_SHORTHAUL_DOMESTIC} shipment`,
        { shipment_type: 'HHG_SHORTHAUL_DOMESTIC' },
        {},
      ],
      ['updateMoveTaskOrderStatus', 'Created Move Task Order (MTO)', {}, { status: 'APPROVED' }],
      ['updateMoveTaskOrderStatus', 'Rejected Move Task Order (MTO)', {}, { status: 'REJECTED' }],
    ])('for event name %s it returns %s', (eventName, text, oldValues, changedValues) => {
      const displayText = eventNamePlainTextToDisplay[eventName](changedValues, oldValues);
      expect(displayText).toEqual(text);
    });
  });
});
