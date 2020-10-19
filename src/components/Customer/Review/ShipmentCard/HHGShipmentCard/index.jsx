import React from 'react';
import { string, shape, number, func } from 'prop-types';

import { AddressShape } from '../../../../../types/address';
import styles from '../ShipmentCard.module.scss';
import PickupDisplay from '../PickupDisplay';

import { formatCustomerDestination } from 'utils/shipmentDisplay';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import { formatCustomerDate } from 'utils/formatters';

const HHGShipmentCard = ({
  destinationLocation,
  destinationZIP,
  moveId,
  onEditClick,
  pickupLocation,
  receivingAgent,
  releasingAgent,
  remarks,
  requestedDeliveryDate,
  requestedPickupDate,
  shipmentId,
  shipmentNumber,
  shipmentType,
}) => {
  const editPath = `/moves/${moveId}/mto-shipments/${shipmentId}/edit-shipment?shipmentNumber=${shipmentNumber}`;
  return (
    <div className={styles.ShipmentCard} data-testid="hhg-summary">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <PickupDisplay
          shipmentId={shipmentId}
          shipmentType={shipmentType}
          shipmentNumber={shipmentNumber}
          requestedPickupDate={requestedPickupDate}
          pickupLocation={pickupLocation}
          releasingAgent={releasingAgent}
          onEditClick={onEditClick}
          editPath={editPath}
          ableToEdit
        />

        <div className={styles.row}>
          <dt>Requested delivery date</dt>
          <dd>{formatCustomerDate(requestedDeliveryDate)}</dd>
        </div>
        <div className={styles.row}>
          <dt>Destination</dt>
          <dd>{formatCustomerDestination(destinationLocation, destinationZIP)}</dd>
        </div>

        {receivingAgent && (
          <div className={styles.row}>
            <dt>Receiving agent</dt>
            <dd>
              {(receivingAgent.firstName || receivingAgent.lastName) && (
                <>
                  {receivingAgent.firstName} {receivingAgent.lastName} <br />
                </>
              )}
              {receivingAgent.phone && (
                <>
                  {receivingAgent.phone} <br />
                </>
              )}
              {receivingAgent.email}
            </dd>
          </div>
        )}
        {remarks && (
          <div className={`${styles.row} ${styles.remarksRow}`}>
            <dt>Remarks</dt>
            <dd className={styles.remarksCell}>{remarks}</dd>
          </div>
        )}
      </ShipmentContainer>
    </div>
  );
};

HHGShipmentCard.propTypes = {
  moveId: string.isRequired,
  shipmentNumber: number.isRequired,
  shipmentType: string.isRequired,
  shipmentId: string.isRequired,
  requestedPickupDate: string.isRequired,
  pickupLocation: AddressShape.isRequired,
  destinationLocation: AddressShape,
  releasingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  requestedDeliveryDate: string.isRequired,
  destinationZIP: string.isRequired,
  onEditClick: func.isRequired,
  receivingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  remarks: string,
};

HHGShipmentCard.defaultProps = {
  destinationLocation: null,
  releasingAgent: null,
  receivingAgent: null,
  remarks: '',
};

export default HHGShipmentCard;
