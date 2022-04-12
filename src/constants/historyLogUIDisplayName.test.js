import {
  getHistoryLogEventNameDisplay,
  detailsPlainTextToDisplay,
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

  describe('displays the correct plain text when each detailsPlainTextToDisplay is called', () => {
    it.each([
      [
        { eventName: 'approveShipment', oldValues: { shipment_type: 'HHG' }, changedValues: { status: 'APPROVED' } },
        'HHG shipment',
      ],
      [
        { eventName: 'requestShipmentDiversion', oldValues: { shipment_type: 'HHG' } },
        `Requested diversion for ${shipmentOptionToDisplay.HHG} shipment`,
      ],
      [
        { eventName: 'requestShipmentDiversion', oldValues: { shipment_type: 'HHG_OUTOF_NTS_DOMESTIC' } },
        `Requested diversion for ${shipmentOptionToDisplay.HHG_OUTOF_NTS_DOMESTIC} shipment`,
      ],
      [
        { eventName: 'requestShipmentDiversion', oldValues: { shipment_type: 'HHG_INTO_NTS_DOMESTIC' } },
        `Requested diversion for ${shipmentOptionToDisplay.HHG_INTO_NTS_DOMESTIC} shipment`,
      ],
      [
        { eventName: 'requestShipmentDiversion', oldValues: { shipment_type: 'PPM' } },
        `Requested diversion for ${shipmentOptionToDisplay.PPM} shipment`,
      ],
      [
        { eventName: 'requestShipmentDiversion', oldValues: { shipment_type: 'HHG_SHORTHAUL_DOMESTIC' } },
        `Requested diversion for ${shipmentOptionToDisplay.HHG_SHORTHAUL_DOMESTIC} shipment`,
      ],
      [
        {
          eventName: 'updateMTOServiceItemStatus',
          context: { name: 'Domestic origin price', shipment_type: 'HHG_INTO_NTS_DOMESTIC' },
        },
        'NTS shipment, Domestic origin price',
      ],
      [
        { eventName: 'setFinancialReviewFlag', changedValues: { financial_review_flag: 'true' } },
        'Move flagged for financial review',
      ],
      [
        { eventName: 'setFinancialReviewFlag', changedValues: { financial_review_flag: 'false' } },
        'Move unflagged for financial review',
      ],
      [
        { eventName: 'requestShipmentCancellation', oldValues: { shipment_type: 'HHG' } },
        `Requested cancellation for ${shipmentOptionToDisplay.HHG} shipment`,
      ],
      [
        { eventName: 'requestShipmentCancellation', oldValues: { shipment_type: 'HHG_OUTOF_NTS_DOMESTIC' } },
        `Requested cancellation for ${shipmentOptionToDisplay.HHG_OUTOF_NTS_DOMESTIC} shipment`,
      ],
      [
        { eventName: 'requestShipmentCancellation', oldValues: { shipment_type: 'HHG_INTO_NTS_DOMESTIC' } },
        `Requested cancellation for ${shipmentOptionToDisplay.HHG_INTO_NTS_DOMESTIC} shipment`,
      ],
      [
        { eventName: 'requestShipmentCancellation', oldValues: { shipment_type: 'PPM' } },
        `Requested cancellation for ${shipmentOptionToDisplay.PPM} shipment`,
      ],
      [
        { eventName: 'requestShipmentCancellation', oldValues: { shipment_type: 'HHG_SHORTHAUL_DOMESTIC' } },
        `Requested cancellation for ${shipmentOptionToDisplay.HHG_SHORTHAUL_DOMESTIC} shipment`,
      ],
      [
        { eventName: 'updateMoveTaskOrderStatus', changedValues: { status: 'APPROVED' } },
        'Created Move Task Order (MTO)',
      ],
      [
        { eventName: 'updateMoveTaskOrderStatus', changedValues: { status: 'REJECTED' } },
        'Rejected Move Task Order (MTO)',
      ],
      [{ eventName: 'acknowledgeExcessWeightRisk' }, 'Dismissed excess weight alert'],
    ])('for history record %s it returns %s', (historyRecord, text) => {
      const displayText = detailsPlainTextToDisplay(historyRecord);
      expect(displayText).toEqual(text);
    });
  });
});
