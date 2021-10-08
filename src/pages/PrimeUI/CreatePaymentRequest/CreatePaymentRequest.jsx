import React from 'react';
import { useParams } from 'react-router-dom';

import LoadingPlaceholder from '../../../shared/LoadingPlaceholder';
import SomethingWentWrong from '../../../shared/SomethingWentWrong';

import { usePrimeSimulatorGetMove } from 'hooks/queries';

function parseServiceItem(serviceItem) {
  let str = '';

  str = `ID: ${serviceItem.id}\n`;
  str += `Service Code: ${serviceItem.reServiceCode}\n`;
  str += `Service Name: ${serviceItem.reServiceName}\n`;
  str += `eTag: ${serviceItem.eTag}\n`;

  return str;
}

function parseAddress(address) {
  let str = '';

  str += `\n\t${address.streetAddress1}\n\t${address.streetAddress2}\n\t${address.streetAddress2}\n`;
  str += `\t${address.city}, ${address.state}, ${address.postalCode}\n`;

  return str;
}

function parseMoveForPaymentRequest(moveTaskOrders) {
  /*
  ================================================
  Move Code:
  Move ID:
  ------------------------------------------------
  Move Service Items:
  ------------------------------------------------
  Shipments (#)
  ************************************************
  ID:
  Shipment ID:
  Shipment Type:
  Shipment ETag:
  Requested Pickup Date:
  Actual Pickup Date:
  Estimated Weight:
  Actual Weight:
  Reweigh Weight:
  Pickup Address:
  Destination Address:
  ------------------------------------------------
  Shipment Service Items:
  Service Item ID:
  Service Item Code:
  Service Item Description:
  ************************************************
  */

  let moveDetails = '';
  let move = moveTaskOrders.moveTaskOrders;
  const firstkey = Object.keys(move)[0];
  move = move[firstkey];

  const equalDivider = '=';
  moveDetails += equalDivider.repeat(50);
  moveDetails += '\n';
  moveDetails += `Move Code: ${move.moveCode}\n`;
  moveDetails += `Move ID: ${move.id}\n`;

  const lineDivider = '-';
  moveDetails += lineDivider.repeat(50);
  moveDetails += '\n';

  moveDetails += 'Move Service Items:\n';
  moveDetails += lineDivider.repeat(50);
  moveDetails += '\n';

  const MoveServiceCodes = ['MS', 'CS'];
  // Get MS and CS service items
  for (let i = 0; i < move.mtoServiceItems.length; i += 1) {
    if (MoveServiceCodes.includes(move.mtoServiceItems[i].reServiceCode)) {
      moveDetails += parseServiceItem(move.mtoServiceItems[i]);
      moveDetails += lineDivider.repeat(50);
      moveDetails += '\n';
    }
  }

  moveDetails += lineDivider.repeat(50);
  moveDetails += '\n';

  moveDetails += `Shipments (${move.mtoShipments.length}):\n`;
  const starDivider = '*';
  moveDetails += starDivider.repeat(50);
  moveDetails += '\n';

  for (let i = 0; i < move.mtoShipments.length; i += 1) {
    const shipment = move.mtoShipments[i];
    moveDetails += `Shipment ID: ${shipment.id}\n`;
    moveDetails += `Shipment Type: ${shipment.shipmentType}\n`;
    moveDetails += `Shipment eTag: ' ${shipment.eTag}\n`;
    moveDetails += `Requested Pickup Date: ${shipment.requestedPickupDate}\n`;
    moveDetails += `Actual Pickup Date: ${shipment.actualPickupDate}\n`;
    moveDetails += `Estimated Weight: ${shipment.primeEstimatedWeight}\n`;
    moveDetails += `Actual Weight: ${shipment.primeActualWeight}\n`;
    moveDetails += `Pickup Address: ${parseAddress(shipment.pickupAddress)}\n`;
    moveDetails += `Destination Address: ${parseAddress(shipment.destinationAddress)}\n`;
    moveDetails += lineDivider.repeat(50);
    moveDetails += '\n';

    moveDetails += 'Shipment Service Items:\n';
    moveDetails += lineDivider.repeat(50);
    moveDetails += '\n';
    for (let j = 0; j < move.mtoServiceItems.length; j += 1) {
      if (move.mtoServiceItems[j].mtoShipmentID === shipment.id) {
        moveDetails += parseServiceItem(move.mtoServiceItems[j]);
        moveDetails += lineDivider.repeat(50);
        moveDetails += '\n';
      }
    }
    moveDetails += starDivider.repeat(50);
    moveDetails += '\n';
  }

  return moveDetails;
}

const CreatePaymentRequest = () => {
  const { moveCode } = useParams();

  const { data, isLoading, isError } = usePrimeSimulatorGetMove(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return <div style={{ 'white-space': 'pre' }}>{parseMoveForPaymentRequest(data)}</div>;
  // return <div style={{ 'white-space': 'pre' }}>{JSON.stringify(data, null, '\t')}</div>;
};

export default CreatePaymentRequest;
