import React from 'react';
import * as PropTypes from 'prop-types';
import { Checkbox, Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import ShipmentContainer from '../ShipmentContainer';

import styles from './ShipmentDisplay.module.scss';

import { AddressShape } from 'types/address';
import { formatAddress } from 'utils/shipmentDisplay';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatDate } from 'shared/dates';

const ShipmentDisplay = ({ shipmentType, displayInfo, onChange, shipmentId, isSubmitted, showIcon, editURL }) => {
  const editButtonClasses = classnames(
    'usa-button',
    'usa-button--secondary',
    'usa-button--small',
    'usa-button--icon',
    'margin-left-1',
    styles.editButton,
  );

  const containerClasses = classnames(styles.container, { [styles.noIcon]: !showIcon });

  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={containerClasses} shipmentType={shipmentType}>
        <div className={styles.heading}>
          {showIcon && isSubmitted && (
            <Checkbox
              id={`shipment-display-checkbox-${shipmentId}`}
              data-testid="shipment-display-checkbox"
              onChange={onChange}
              name="shipments"
              label=""
              value={shipmentId}
            />
          )}
          {showIcon && !isSubmitted && <FontAwesomeIcon icon={['far', 'check-circle']} className={styles.approved} />}
          <h3>{displayInfo.heading}</h3>
          <FontAwesomeIcon icon="chevron-down" />
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
          <div className={styles.row}>
            <dt className={styles.label}>Counselor remarks</dt>
            <dd data-testid="counselorRemarks">{displayInfo.counselorRemarks || '—'}</dd>
          </div>
        </dl>
        {editURL && (
          <Button to={editURL} className={editButtonClasses} data-testid="editButton">
            <span className="icon">
              <FontAwesomeIcon icon="pen" alt=" " inverse />
            </span>
            <span>Edit shipment</span>
          </Button>
        )}
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
    SHIPMENT_OPTIONS.NTSR,
  ]),
  displayInfo: PropTypes.shape({
    heading: PropTypes.string.isRequired,
    requestedMoveDate: PropTypes.string.isRequired,
    currentAddress: AddressShape.isRequired,
    destinationAddress: AddressShape,
    counselorRemarks: PropTypes.string,
  }).isRequired,
  showIcon: PropTypes.bool,
  editURL: PropTypes.string,
};

ShipmentDisplay.defaultProps = {
  onChange: () => {},
  shipmentType: SHIPMENT_OPTIONS.HHG,
  showIcon: true,
  editURL: '',
};

export default ShipmentDisplay;
