import React from 'react';
import * as PropTypes from 'prop-types';
import { Checkbox } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheckCircle } from '@fortawesome/free-regular-svg-icons';
import { faChevronDown } from '@fortawesome/free-solid-svg-icons';

import ShipmentContainer from '../ShipmentContainer';
import { AddressShape } from '../../../types/address';

import styles from './ShipmentDisplay.module.scss';

import { formatAddress } from 'utils/shipmentDisplay';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatDate } from 'shared/dates';

const ShipmentDisplay = ({ shipmentType, displayInfo, onChange, shipmentId, isSubmitted }) => {
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <div className={styles.heading}>
          {isSubmitted && (
            <Checkbox
              id={`shipment-display-checkbox-${shipmentId}`}
              data-testid="shipment-display-checkbox"
              onChange={onChange}
              name="shipments"
              label=""
              value={shipmentId}
            />
          )}
          {!isSubmitted && <FontAwesomeIcon icon={faCheckCircle} />}
          <h3>{displayInfo.heading}</h3>
          <FontAwesomeIcon icon={faChevronDown} />
        </div>
        <dl>
          <div className={styles.row}>
            <dt>Requested move date</dt>
            <dd>{formatDate(displayInfo.requestedMoveDate, 'DD MMM YYYY')}</dd>
          </div>
          <div className={styles.row}>
            <dt>Current address</dt>
            <dd>{displayInfo.currentAddress && formatAddress(displayInfo.currentAddress)}</dd>
          </div>
          <div className={styles.row}>
            <dt className={styles.label}>Destination address</dt>
            <dd data-testid="shipmentDestinationAddress">{formatAddress(displayInfo.destinationAddress)}</dd>
          </div>
        </dl>
      </ShipmentContainer>
    </div>
  );
};

ShipmentDisplay.propTypes = {
  onChange: PropTypes.func,
  shipmentId: PropTypes.string.isRequired,
  isSubmitted: PropTypes.bool.isRequired,
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
  ]),
  displayInfo: PropTypes.shape({
    heading: PropTypes.string.isRequired,
    requestedMoveDate: PropTypes.string.isRequired,
    currentAddress: AddressShape.isRequired,
    destinationAddress: AddressShape,
  }).isRequired,
};

ShipmentDisplay.defaultProps = {
  onChange: () => {},
  shipmentType: SHIPMENT_OPTIONS.HHG,
};

export default ShipmentDisplay;
