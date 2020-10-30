import React from 'react';
import * as PropTypes from 'prop-types';
import { Checkbox } from '@trussworks/react-uswds';

import ShipmentContainer from '../ShipmentContainer';
import { AddressShape } from '../../../types/address';

import styles from './ShipmentDisplay.module.scss';

import { ReactComponent as ChevronDown } from 'shared/icon/chevron-down.svg';
import { formatAddress } from 'utils/shipmentDisplay';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatDate } from 'shared/dates';
import { ReactComponent as CheckmarkIcon } from 'shared/icon/checkbox--unchecked.svg';

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
          {!isSubmitted && <CheckmarkIcon />}
          <h3>{displayInfo.heading}</h3>
          <ChevronDown />
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
