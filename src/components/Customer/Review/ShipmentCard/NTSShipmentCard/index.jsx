import React from 'react';
import { string, shape, number, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../../../types/address';
import styles from '../ShipmentCard.module.scss';
import PickupDisplay from '../PickupDisplay';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';

const NTSShipmentCard = ({
  moveId,
  onEditClick,
  pickupLocation,
  releasingAgent,
  remarks,
  requestedPickupDate,
  shipmentId,
  shipmentType,
  shipmentNumber,
}) => {
  const editPath = `/moves/${moveId}/mto-shipments/${shipmentId}/edit-shipment`;
  return (
    <div className={styles.ShipmentCard} data-testid="nts-summary">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h3>{getShipmentTypeLabel(shipmentType)}</h3>
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
          <PickupDisplay
            shipmentId={shipmentId}
            shipmentType={shipmentType}
            shipmentNumber={shipmentNumber}
            requestedPickupDate={requestedPickupDate}
            pickupLocation={pickupLocation}
            releasingAgent={releasingAgent}
          />
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
  moveId: string.isRequired,
  onEditClick: func.isRequired,
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
  shipmentNumber: number,
};

NTSShipmentCard.defaultProps = {
  releasingAgent: null,
  remarks: '',
  shipmentNumber: 0,
};

export default NTSShipmentCard;
