import { getHistoryLogEventNameDisplay } from './historyLogUIDisplayName';

describe('historyLogUIDisplay', () => {
  describe('display Submitted orders', () => {
    const item = {
      changedValues: [
        {
          columnName: 'status',
          columnValue: 'DRAFT',
        },
      ],
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
      changedValues: [
        {
          columnName: 'status',
          columnValue: 'APPROVED',
        },
      ],
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

  describe('display Move rejected', () => {
    const item = {
      changedValues: [
        {
          columnName: 'status',
          columnValue: 'REJECTED',
        },
      ],
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
