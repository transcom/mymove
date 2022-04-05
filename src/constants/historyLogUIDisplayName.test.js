import { getHistoryLogEventNameDisplay } from './historyLogUIDisplayName';

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
});
