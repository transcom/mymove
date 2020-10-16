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
  shipmentType,
}) => {
  const editPath = `/moves/${moveId}/mto-shipments/${shipmentId}/edit-shipment?shipmentNumber=${shipmentNumber}`;
  const isNTS = shipmentType ? shipmentType === SHIPMENT_OPTIONS.NTS : false;
  const isNTSR = shipmentType ? shipmentType === SHIPMENT_OPTIONS.NTSR : false;
  const isHHG = shipmentType ? shipmentType === SHIPMENT_OPTIONS.HHG : false;
  return (
    <div className={styles.ShipmentCard} data-testid="hhg-summary">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h3>
              {shipmentType} {shipmentNumber}
            </h3>
            <p>#{shipmentId.substring(0, 8).toUpperCase()}</p>
          </div>
          <Button
            className={styles.editBtn}
            data-testid="edit-shipment-btn"
            onClick={() => onEditClick(editPath)}
            unstyled
            disabled={isNTS || isNTSR}
          >
            Edit
          </Button>
        </div>

        <dl className={styles.shipmentCardSubsection}>
          {isHHG ||
            (isNTS && (
              <div className={styles.row}>
                <dt>Requested pickup date</dt>
                <dd>{formatCustomerDate(requestedPickupDate)}</dd>
              </div>
            ))}
          {isHHG ||
            (isNTS && (
              <div className={styles.row}>
                <dt>Pickup location</dt>
                <dd>
                  {pickupLocation.street_address_1} {pickupLocation.street_address_2}
                  <br />
                  {pickupLocation.city}, {pickupLocation.state} {pickupLocation.postal_code}
                </dd>
              </div>
            ))}
          {isHHG ||
            (isNTS && releasingAgent && (
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
            ))}
          {isHHG ||
            (isNTSR && (
              <div className={styles.row}>
                <dt>Requested delivery date</dt>
                <dd>{formatCustomerDate(requestedDeliveryDate)}</dd>
              </div>
            ))}
          {isHHG ||
            (isNTSR && (
              <div className={styles.row}>
                <dt>Destination</dt>
                <dd>{formatCustomerDestination(destinationLocation, destinationZIP)}</dd>
              </div>
            ))}
          {isHHG ||
            (isNTSR && receivingAgent && (
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
            ))}
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
  shipmentType: string.isRequired,
  shipmentId: string.isRequired,
  requestedPickupDate: string,
  pickupLocation: AddressShape,
  destinationLocation: AddressShape,
  releasingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  requestedDeliveryDate: string,
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
  requestedDeliveryDate: '',
  requestedPickupDate: '',
  pickupLocation: {},
};

export default HHGShipmentCard;
