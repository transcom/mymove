import React from 'react';
import * as PropTypes from 'prop-types';
import { Checkbox, Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import { EditButton } from 'components/form/IconButtons';
import ShipmentContainer from 'components/Office/ShipmentContainer';
import ShipmentInfoList from 'components/Office/DefinitionLists/ShipmentInfoList';
import styles from 'components/Office/ShipmentDisplay/ShipmentDisplay.module.scss';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape } from 'types/address';
import { shipmentStatuses } from 'constants/shipments';

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
          <div className={styles.headingTagWrapper}>
            <h3>
              <label id={`shipment-display-label-${shipmentId}`}>{displayInfo.heading}</label>
            </h3>
            {displayInfo.isDiversion && <Tag>diversion</Tag>}
            {displayInfo.shipmentStatus === shipmentStatuses.CANCELED && <Tag className="usa-tag--red">cancelled</Tag>}
          </div>

          <FontAwesomeIcon icon="chevron-down" />
        </div>
        <ShipmentInfoList className={styles.shipmentDisplayInfo} shipment={displayInfo} />
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
    isDiversion: PropTypes.bool,
    shipmentStatus: PropTypes.string,
    requestedPickupDate: PropTypes.string.isRequired,
    pickupAddress: AddressShape.isRequired,
    secondaryPickupAddress: AddressShape,
    destinationAddress: AddressShape,
    secondaryDeliveryAddress: AddressShape,
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
