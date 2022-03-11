// TODO: it is possible to have a clash with the operationId and the yaml file?

const historyLogEventNameDisplay = new Map([
  //operationId, UI Display
  ['counselingUpdateOrder', 'Updated orders'], //ghc.yaml
  ['updateOrder', 'Updated orders'], //ghc.yaml
  ['updateAllowance', 'Updated allowances'], //ghc.yaml
  ['counselingUpdateAllowance', 'Updated allowances'], //ghc.yaml
  ['updateMoveTaskOrder', 'Updated move'], //ghc.yaml
  ['updateMTOShipment', 'Updated shipment'], //ghc.yaml
  ['approveShipment', 'Approved shipment'], //ghc.yaml
  ['requestShipmentDiversion', 'Requested diversion'], //ghc.yaml
  ['updateMTOServiceItem', 'Updated service item'], //ghc.yaml
  ['updateMTOServiceItemStatus', ''], //ghc.yaml
  ['updateMoveTaskOrderStatus', ''], //ghc.yaml
  ['setFinancialReviewFlag', 'Flagged move'], //ghc.yaml
  ['createOrders', 'Submitted orders'], //internal.yaml
  ['updateOrders', 'Updated orders'], //internal.yaml
  ['uploadAmendedOrders', 'Updated orders'], //internal.yaml
  ['submitMoveForApproval', 'Submitted move'], //internal.yaml
  ['submitAmendedOrders', 'Updated orders'], //internal.yaml
  ['createMTOShipment', 'Submitted shipments'], //internal.yaml DUPLICATE key
  ['updateMTOShipment', 'Updated shipments'], //internal.yaml
  ['createMTOShipment', 'Requested  shipment approval/Shipment request'], //prime.yaml DUPLICATE key
  ['updateMTOShipment', 'Updated shipment'], //prime.yaml
  ['updateMTOShipmentAddress', 'Updated shipment'], //prime.yaml
  ['createMTOServiceItem', 'Requested service item'], //prime.yaml
]);

/*
{
  "action": "INSERT",
  "actionTstampClk": "2022-03-01T16:01:12.217Z",
  "actionTstampStm": "2022-03-01T16:01:12.215Z",
  "actionTstampTx": "2022-03-01T16:01:12.188Z",
  "changedValues": [
  {
    "columnName": "financial_review_flag",
    "columnValue": "false"
  },
  {
    "columnName": "id",
    "columnValue": "4a0d1cb3-a1f7-40f4-911f-79a22c0c4f56"
  },
  {
    "columnName": "show",
    "columnValue": "true"
  },
  {
    "columnName": "reference_id",
    "columnValue": "3775-6135"
  },
  {
    "columnName": "contractor_id",
    "columnValue": "5db13bb4-6d29-4bdb-bc81-262f4513ecf6"
  },
  {
    "columnName": "locator",
    "columnValue": "YHHFVB"
  },
  {
    "columnName": "created_at",
    "columnValue": "2022-03-01T16:01:12.214551"
  },
  {
    "columnName": "updated_at",
    "columnValue": "2022-03-01T16:01:12.214551"
  },
  {
    "columnName": "status",
    "columnValue": "DRAFT"
  },
  {
    "columnName": "selected_move_type"
  },
  {
    "columnName": "orders_id",
    "columnValue": "e65e8bc8-38be-49bb-b378-cf495aa3861d"
  }
],
  "eventName": "createOrders",
  "id": "fcff54a3-8c70-41d6-b0e1-db394d2f585f",
  "objectId": "4a0d1cb3-a1f7-40f4-911f-79a22c0c4f56",
  "relId": 25807,
  "sessionUserId": "5c3021d9-22bd-46b2-8055-27d0ba45d8f4",
  "tableName": "moves",
  "transactionId": 154032
},
 */

function getHistoryLogEventNameDisplay({ eventName /*operationId*/, changedValues }) {
  switch (eventName) {
    case 'updateMTOServiceItemStatus': {
      // find 'columnName' with 'columnValue'
      const status = changedValues.find((changedValue) => changedValue['columnName'] === 'status');
      switch (status['columnValue']) {
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
      const status = changedValues.find((changedValue) => changedValue['columnName'] === 'status');
      switch (status['columnValue']) {
        case 'APPROVED':
          return 'Move approved';
        case 'REJECTED':
          return 'Move rejected';
        default:
          return '';
      }
    }
    default:
      return historyLogEventNameDisplay.get(eventName);
  }
}
