import React from 'react';
import { string, shape, number } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../../../types/address';
import styles from '../ShipmentCard.module.scss';

import hhgShipmentCardStyles from './HHGShipmentCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatCustomerDate } from 'shared/utils';

const HHGShipmentCard = ({
  shipmentNumber,
  shipmentId,
  requestedPickupDate,
  pickupLocation,
  releasingAgent,
  requestedDeliveryDate,
  destinationZIP,
  receivingAgent,
  remarks,
}) => {
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={SHIPMENT_OPTIONS.HHG}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h3>HHG {shipmentNumber}</h3>
            <p>{shipmentId.substring(0, 10)}</p>
          </div>
          <Button className={styles.editBtn} onClick={() => {}} unstyled>
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
                {releasingAgent.name && (
                  <>
                    {releasingAgent.name} <br />
                  </>
                )}
                {releasingAgent.telephone && (
                  <>
                    {releasingAgent.telephone} <br />
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
            <dd>{destinationZIP}</dd>
          </div>
          {receivingAgent && (
            <div className={styles.row}>
              <dt>Receiving agent</dt>
              <dd>
                {receivingAgent.name && (
                  <>
                    {receivingAgent.name} <br />
                  </>
                )}
                {receivingAgent.telephone && (
                  <>
                    {receivingAgent.telephone} <br />
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
  shipmentNumber: number.isRequired,
  shipmentId: string.isRequired,
  requestedPickupDate: string.isRequired,
  pickupLocation: AddressShape.isRequired,
  releasingAgent: shape({
    name: string,
    telephone: string,
    email: string,
  }),
  requestedDeliveryDate: string.isRequired,
  destinationZIP: string.isRequired,
  receivingAgent: shape({
    name: string,
    telephone: string,
    email: string,
  }),
  remarks: string,
};

HHGShipmentCard.defaultProps = {
  releasingAgent: null,
  receivingAgent: null,
  remarks: '',
};

export default HHGShipmentCard;
