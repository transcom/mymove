import React from 'react';
import { string, shape, number } from 'prop-types';

import { AddressShape } from '../../../../../types/address';
import styles from '../ShipmentCard.module.scss';
import PickupDisplay from '../PickupDisplay';

import ShipmentContainer from 'components/Office/ShipmentContainer';

const NTSShipmentCard = ({
  pickupLocation,
  releasingAgent,
  remarks,
  requestedPickupDate,
  shipmentId,
  shipmentType,
  shipmentNumber,
}) => {
  return (
    <div className={styles.ShipmentCard} data-testid="nts-summary">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <PickupDisplay
          shipmentId={shipmentId}
          shipmentType={shipmentType}
          shipmentNumber={shipmentNumber}
          requestedPickupDate={requestedPickupDate}
          pickupLocation={pickupLocation}
          releasingAgent={releasingAgent}
          onEditClick={() => {}}
          ableToEdit={false}
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
  shipmentNumber: number,
};

NTSShipmentCard.defaultProps = {
  releasingAgent: null,
  remarks: '',
  shipmentNumber: 0,
};

export default NTSShipmentCard;
