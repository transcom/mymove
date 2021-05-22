import React from 'react';
import * as PropTypes from 'prop-types';
import { Checkbox } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import { EditButton } from 'components/form/IconButtons';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import styles from 'components/Office/ShipmentDisplay/ShipmentDisplay.module.scss';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatDate } from 'shared/dates';
import { AddressShape } from 'types/address';
import { formatAddress } from 'utils/shipmentDisplay';

const ShipmentDisplay = ({ shipmentType, displayInfo, onChange, shipmentId, isSubmitted, showIcon, editURL }) => {
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
              aria-labelledby={`shipment-display-label-${shipmentId}`}
            />
          )}

          {showIcon && !isSubmitted && <FontAwesomeIcon icon={['far', 'check-circle']} className={styles.approved} />}
          <h3>
            <label id={`shipment-display-label-${shipmentId}`}>{displayInfo.heading}</label>
          </h3>

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
          <EditButton
            onClick={() => {
              window.location.href = editURL;
            }}
            className={styles.editButton}
            data-testid={editURL}
            label="Edit shipment"
            secondary
          />
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
