import React from 'react';
import { string, shape } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../../../types/address';
import styles from '../ShipmentCard.module.scss';
import DeliveryDisplay from '../DeliveryDisplay';

import { getShipmentTypeLabel } from 'utils/shipmentDisplay';
import ShipmentContainer from 'components/Office/ShipmentContainer';

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

        <DeliveryDisplay
          shipmentId={shipmentId}
          shipmentType={shipmentType}
          requestedDeliveryDate={requestedDeliveryDate}
          destinationLocation={destinationLocation}
          destinationZIP={destinationZIP}
          receivingAgent={receivingAgent}
        />
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
