import React from 'react';
import { string, shape } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../../../types/address';
import styles from '../ShipmentCard.module.scss';

import { getShipmentTypeLabel } from 'utils/shipmentDisplay';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import { formatCustomerDate } from 'utils/formatters';

const NTSShipmentCard = ({
  pickupLocation,
  releasingAgent,
  remarks,
  requestedPickupDate,
  shipmentId,
  shipmentType,
}) => {
  return (
    <div className={styles.ShipmentCard} data-testid="nts-summary">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h3>{getShipmentTypeLabel(shipmentType)}</h3>
            <p>#{shipmentId.substring(0, 8).toUpperCase()}</p>
          </div>
          <Button className={styles.editBtn} data-testid="edit-shipment-btn" unstyled disabled>
            Edit
          </Button>
        </div>

        <dl className={styles.shipmentCardSubsection}>
          <div className={styles.row}>
            <dt>Requested pickup date</dt>
            <dd>{formatCustomerDate(requestedPickupDate)}</dd>
          </div>
          {pickupLocation && (
            <div className={styles.row}>
              <dt>Pickup location</dt>
              <dd>
                {pickupLocation.street_address_1} {pickupLocation.street_address_2}
                <br />
                {pickupLocation.city}, {pickupLocation.state} {pickupLocation.postal_code}
              </dd>
            </div>
          )}
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
          {remarks && (
            <div className={`${styles.row} ${styles.remarksRow}`}>
              <dt>Remarks</dt>
              <dd className={styles.remarksCell}>{remarks}</dd>
            </div>
          )}
        </dl>
      </ShipmentContainer>
    </div>
  );
};

NTSShipmentCard.propTypes = {
  shipmentType: string.isRequired,
  shipmentId: string.isRequired,
  requestedPickupDate: string.isRequired,
  pickupLocation: AddressShape.isRequired,
  releasingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  remarks: string,
};

NTSShipmentCard.defaultProps = {
  releasingAgent: null,
  remarks: '',
};

export default NTSShipmentCard;
