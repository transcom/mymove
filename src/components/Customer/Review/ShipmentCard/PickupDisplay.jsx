import React from 'react';
import { string, shape, func, bool, number } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../../types/address';

import styles from './ShipmentCard.module.scss';

import { getShipmentTypeLabel } from 'utils/shipmentDisplay';
import { formatCustomerDate } from 'utils/formatters';

const PickupDisplay = ({
  pickupLocation,
  releasingAgent,
  onEditClick,
  editPath,
  ableToEdit,
  requestedPickupDate,
  shipmentId,
  shipmentType,
  shipmentNumber,
}) => {
  return (
    <div>
      <div className={styles.ShipmentCardHeader}>
        <div>
          <h3>
            {getShipmentTypeLabel(shipmentType)} {shipmentNumber || ''}
          </h3>
          <p>#{shipmentId.substring(0, 8).toUpperCase()}</p>
        </div>
        <Button
          className={styles.editBtn}
          data-testid="edit-shipment-btn"
          onClick={() => onEditClick(editPath)}
          unstyled
          disabled={!ableToEdit}
        >
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
      </dl>
    </div>
  );
};

PickupDisplay.propTypes = {
  shipmentType: string.isRequired,
  shipmentNumber: number,
  shipmentId: string.isRequired,
  requestedPickupDate: string.isRequired,
  pickupLocation: AddressShape.isRequired,
  releasingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  onEditClick: func,
  editPath: string,
  ableToEdit: bool,
};

PickupDisplay.defaultProps = {
  releasingAgent: null,
  onEditClick: null,
  editPath: '',
  ableToEdit: false,
  shipmentNumber: 0,
};

export default PickupDisplay;
