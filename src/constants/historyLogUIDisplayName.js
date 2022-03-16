const historyLogEventNameDisplay = {
  // operationId, UI Display
  counselingUpdateOrder: 'Updated orders', // ghc.yaml
  updateOrder: 'Updated orders', // ghc.yaml
  updateAllowance: 'Updated allowances', // ghc.yaml
  counselingUpdateAllowance: 'Updated allowances', // ghc.yaml
  updateMoveTaskOrder: 'Updated move', // ghc.yaml
  updateMTOShipment: 'Updated shipment', // ghc.yaml internal.yaml prime.yaml
  approveShipment: 'Approved shipment', // ghc.yaml
  requestShipmentDiversion: 'Requested diversion', // ghc.yaml
  updateMTOServiceItem: 'Updated service item', // ghc.yaml
  updateMTOServiceItemStatus: '', // ghc.yaml
  updateMoveTaskOrderStatus: '', // ghc.yaml
  setFinancialReviewFlag: 'Flagged move', // ghc.yaml
  requestShipmentCancellation: 'Updated shipment', // ghc.yaml
  createOrders: 'Submitted orders', // internal.yaml
  updateOrders: 'Updated orders', // internal.yaml
  uploadAmendedOrders: 'Updated orders', // internal.yaml
  submitMoveForApproval: 'Submitted move', // internal.yaml
  submitAmendedOrders: 'Updated orders', // internal.yaml
  createMTOShipment: 'Submitted/Requested shipments', // internal.yaml prime.yaml
  updateMTOShipmentAddress: 'Updated shipment', // prime.yaml
  createMTOServiceItem: 'Requested service item', // prime.yaml
};

export function getHistoryLogEventNameDisplay({ eventName /* operationId */, changedValues }) {
  switch (eventName) {
    case 'updateMTOServiceItemStatus': {
      // find 'columnName' with 'columnValue'
      const status = changedValues?.find((changedValue) => changedValue?.columnName === 'status');
      switch (status?.columnValue) {
        case 'APPROVED':
          return 'Approved service item';
        case 'REJECTED':
          return 'Rejected service item';
        default:
          return '';
      }
    }
    case 'updateMoveTaskOrderStatus': {
      // find 'columnName' with 'columnValue'
      const status = changedValues?.find((changedValue) => changedValue?.columnName === 'status');
      switch (status?.columnValue) {
        case 'APPROVED':
          return 'Move approved';
        case 'REJECTED':
          return 'Move rejected';
        default:
          return '';
      }
    }
    default:
      return historyLogEventNameDisplay[eventName];
  }
}

export default getHistoryLogEventNameDisplay;
