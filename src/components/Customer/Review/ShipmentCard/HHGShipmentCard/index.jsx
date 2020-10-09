import React from 'react';
import { string, shape, number, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../../../types/address';
import styles from '../ShipmentCard.module.scss';

import hhgShipmentCardStyles from './HHGShipmentCard.module.scss';

import { formatCustomerDestination } from 'utils/shipmentDisplay';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';
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
}) => {
  const editPath = `/moves/${moveId}/mto-shipments/${shipmentId}/edit-shipment?shipmentNumber=${shipmentNumber}`;
  return (
    <div className={styles.ShipmentCard} data-testid="hhg-summary">
      <ShipmentContainer className={styles.container} shipmentType={SHIPMENT_OPTIONS.HHG}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h3>HHG {shipmentNumber}</h3>
            <p>#{shipmentId.substring(0, 8).toUpperCase()}</p>
          </div>
          <Button
            className={styles.editBtn}
            data-testid="edit-shipment-btn"
            onClick={() => onEditClick(editPath)}
            unstyled
          >
            Edit
          </Button>
        </div>

        <dl className={styles.shipmentCardSubsection}>
          <div className={styles.row}>
            <dt>Requested pickup date</dt>
            <dd>{formatCustomerDate(requestedPickupDate)}</dd>
          </div>
          <div className={styles.row}>
            <dt>Pickup location</dt>
            <dd>
              {pickupLocation.street_address_1} {pickupLocation.street_address_2}
              <br />
              {pickupLocation.city}, {pickupLocation.state} {pickupLocation.postal_code}
            </dd>
          </div>
          {releasingAgent && (
            <div className={styles.row}>
              <dt>Releasing agent</dt>
              <dd>
                {(releasingAgent.firstName || releasingAgent.lastName) && (
                  <>
                    {releasingAgent.firstName} {releasingAgent.lastName} <br />
                  </>
                )}
                {releasingAgent.phone && (
                  <>
                    {releasingAgent.phone} <br />
                  </>
                )}
                {releasingAgent.email}
              </dd>
            </div>
          )}
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
            <div className={`${styles.row} ${hhgShipmentCardStyles.remarksRow}`}>
              <dt>Remarks</dt>
              <dd className={hhgShipmentCardStyles.remarksCell}>{remarks}</dd>
            </div>
          )}
        </dl>
      </ShipmentContainer>
    </div>
  );
};

HHGShipmentCard.propTypes = {
  moveId: string.isRequired,
  shipmentNumber: number.isRequired,
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
