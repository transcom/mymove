import React from 'react';
import classNames from 'classnames';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';

import { AddressShape } from '../../../types/address';

import styles from './shipmentHeading.module.scss';

import { shipmentStatuses } from 'constants/shipments';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

function ShipmentHeading({ shipmentInfo, handleShowCancellationModal, isMoveLocked }) {
  const { shipmentStatus } = shipmentInfo;
  // cancelation modal is visible if shipment is not already canceled, AND if shipment cancellation hasn't already been requested
  const showRequestCancellation =
    shipmentStatus !== shipmentStatuses.CANCELED && shipmentStatus !== shipmentStatuses.CANCELLATION_REQUESTED;
  const isCancellationRequested = shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED;

  return (
    <div className={classNames(styles.shipmentHeading, 'shipment-heading')}>
      <div className={styles.shipmentHeadingType}>
        <span className={styles.marketCodeIndicator}>{shipmentInfo.marketCode}</span>
        <h2>{shipmentInfo.shipmentType}</h2>
        <div>
          {shipmentStatus === shipmentStatuses.CANCELED && <Tag className="usa-tag--cancellation">canceled</Tag>}
          {shipmentInfo.isDiversion && <Tag className="usa-tag--diversion">diversion</Tag>}
          {!shipmentInfo.isDiversion && shipmentStatus === shipmentStatuses.DIVERSION_REQUESTED && (
            <Tag className="usa-tag--diversion">diversion requested</Tag>
          )}
        </div>
      </div>
      <div>
        <h4>#{shipmentInfo.shipmentLocator}</h4>
      </div>
      <div className={styles.row}>
        {showRequestCancellation && (
          <Restricted to={permissionTypes.createShipmentCancellation}>
            <Restricted to={permissionTypes.updateMTOPage}>
              <Button
                data-testid="requestCancellationBtn"
                type="button"
                onClick={() => handleShowCancellationModal(shipmentInfo)}
                unstyled
                disabled={isMoveLocked}
              >
                Request Cancellation
              </Button>
            </Restricted>
          </Restricted>
        )}
        {isCancellationRequested && <Tag className="usa-tag--cancellation">Cancellation Requested</Tag>}
      </div>
    </div>
  );
}

ShipmentHeading.propTypes = {
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    shipmentType: PropTypes.string.isRequired,
    isDiversion: PropTypes.bool,
    originCity: PropTypes.string.isRequired,
    originState: PropTypes.string.isRequired,
    originPostalCode: PropTypes.string.isRequired,
    destinationAddress: AddressShape,
    scheduledPickupDate: PropTypes.string.isRequired,
    shipmentStatus: PropTypes.string.isRequired,
    ifMatchEtag: PropTypes.string.isRequired,
    moveTaskOrderID: PropTypes.string.isRequired,
  }).isRequired,
  handleShowCancellationModal: PropTypes.func.isRequired,
};

export default ShipmentHeading;
