import React from 'react';
import { string, shape } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../../../types/address';
import styles from '../ShipmentCard.module.scss';

import { formatCustomerDestination, getShipmentTypeLabel } from 'utils/shipmentDisplay';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import { formatCustomerDate } from 'utils/formatters';

const NTSRShipmentCard = ({
  destinationLocation,
  destinationZIP,
  receivingAgent,
  remarks,
  requestedDeliveryDate,
  shipmentId,
  shipmentType,
}) => {
  return (
    <div className={styles.ShipmentCard} data-testid="ntsr-summary">
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
            <dt>Requested delivery date</dt>
            <dd>{formatCustomerDate(requestedDeliveryDate)}</dd>
          </div>

          {destinationLocation || destinationZIP ? (
            <div className={styles.row}>
              <dt>Destination</dt>
              <dd>{formatCustomerDestination(destinationLocation, destinationZIP)}</dd>
            </div>
          ) : undefined}

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
        </dl>
      </ShipmentContainer>
    </div>
  );
};

NTSRShipmentCard.propTypes = {
  shipmentType: string.isRequired,
  shipmentId: string.isRequired,
  destinationLocation: AddressShape,
  requestedDeliveryDate: string.isRequired,
  destinationZIP: string.isRequired,
  receivingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  remarks: string,
};

NTSRShipmentCard.defaultProps = {
  destinationLocation: null,
  receivingAgent: null,
  remarks: '',
};

export default NTSRShipmentCard;
