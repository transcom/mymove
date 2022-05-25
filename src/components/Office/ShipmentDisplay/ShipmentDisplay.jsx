import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import { useHistory } from 'react-router-dom';
import { Checkbox, Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import { EditButton } from 'components/form/IconButtons';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import ShipmentInfoListSelector from 'components/Office/DefinitionLists/ShipmentInfoListSelector';
import styles from 'components/Office/ShipmentDisplay/ShipmentDisplay.module.scss';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape } from 'types/address';
import { shipmentStatuses } from 'constants/shipments';
import { ShipmentStatusesOneOf } from 'types/shipment';
import { OrdersLOAShape } from 'types/order';
import { AgentShape } from 'types/agent';
import { retrieveSAC, retrieveTAC } from 'utils/shipmentDisplay';

const ShipmentDisplay = ({
  shipmentType,
  displayInfo,
  onChange,
  shipmentId,
  isSubmitted,
  allowApproval,
  editURL,
  ordersLOA,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  neverShow,
}) => {
  const history = useHistory();
  const containerClasses = classnames(styles.container, { [styles.noIcon]: !allowApproval });
  const [isExpanded, setIsExpanded] = useState(false);
  const tac = retrieveTAC(displayInfo.tacType, ordersLOA);
  const sac = retrieveSAC(displayInfo.sacType, ordersLOA);

  const disableApproval = errorIfMissing.some((requiredInfo) => !displayInfo[requiredInfo]);

  const handleExpandClick = () => {
    setIsExpanded((prev) => !prev);
  };
  const expandableIconClasses = classnames({
    'chevron-up': isExpanded,
    'chevron-down': !isExpanded,
  });

  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={containerClasses} shipmentType={shipmentType}>
        <div className={styles.heading}>
          {allowApproval && isSubmitted && !displayInfo.usesExternalVendor && (
            <Checkbox
              id={`shipment-display-checkbox-${shipmentId}`}
              data-testid="shipment-display-checkbox"
              onChange={onChange}
              name="shipments"
              label=""
              value={shipmentId}
              aria-labelledby={`shipment-display-label-${shipmentId}`}
              disabled={disableApproval}
            />
          )}

          {allowApproval && !isSubmitted && (
            <FontAwesomeIcon icon={['far', 'check-circle']} className={styles.approved} />
          )}
          <div className={styles.headingTagWrapper}>
            <h3>
              <label id={`shipment-display-label-${shipmentId}`}>{displayInfo.heading}</label>
            </h3>
            {displayInfo.isDiversion && <Tag>diversion</Tag>}
            {displayInfo.shipmentStatus === shipmentStatuses.CANCELED && <Tag className="usa-tag--red">cancelled</Tag>}
            {displayInfo.shipmentStatus === shipmentStatuses.DIVERSION_REQUESTED && <Tag>diversion requested</Tag>}
            {displayInfo.shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED && (
              <Tag>cancellation requested</Tag>
            )}
            {displayInfo.usesExternalVendor && <Tag>external vendor</Tag>}
          </div>

          <FontAwesomeIcon className={styles.icon} icon={expandableIconClasses} onClick={handleExpandClick} />
        </div>
        <ShipmentInfoListSelector
          className={styles.shipmentDisplayInfo}
          shipment={{ ...displayInfo, tac, sac }}
          shipmentType={shipmentType}
          isExpanded={isExpanded}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          showWhenCollapsed={showWhenCollapsed}
          neverShow={neverShow}
        />
        {editURL && (
          <EditButton
            onClick={() => {
              history.push(editURL);
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
    SHIPMENT_OPTIONS.PPM,
  ]),
  displayInfo: PropTypes.oneOfType([
    PropTypes.shape({
      heading: PropTypes.string.isRequired,
      isDiversion: PropTypes.bool,
      shipmentStatus: ShipmentStatusesOneOf,
      requestedPickupDate: PropTypes.string,
      pickupAddress: AddressShape,
      secondaryPickupAddress: AddressShape,
      destinationAddress: AddressShape,
      destinationType: PropTypes.string,
      displayDestinationType: PropTypes.bool,
      secondaryDeliveryAddress: AddressShape,
      counselorRemarks: PropTypes.string,
      shipmentId: PropTypes.string,
      shipmentType: PropTypes.string,
      usesExternalVendor: PropTypes.bool,
      customerRemarks: PropTypes.string,
      serviceOrderNumber: PropTypes.string,
      requestedDeliveryDate: PropTypes.string,
      agents: PropTypes.arrayOf(AgentShape),
      primeActualWeight: PropTypes.number,
      storageFacility: PropTypes.shape({
        address: AddressShape.isRequired,
        facilityName: PropTypes.string,
        lotNumber: PropTypes.string,
      }),
      tacType: PropTypes.string,
      sacType: PropTypes.string,
      ntsRecordedWeight: PropTypes.number,
    }),
    PropTypes.shape({
      heading: PropTypes.string.isRequired,
      shipmentType: PropTypes.string,
      hasRequestedAdvance: PropTypes.bool,
      advanceAmountRequested: PropTypes.number,
      destinationPostalCode: PropTypes.string,
      estimatedIncentive: PropTypes.number,
      estimatedWeight: PropTypes.string,
      expectedDepartureDate: PropTypes.string,
      pickupPostalCode: PropTypes.string,
      proGearWeight: PropTypes.string,
      secondaryDestinationPostalCode: PropTypes.string,
      secondaryPickupPostalCode: PropTypes.string,
      spouseProGearWeight: PropTypes.string,
      customerRemarks: PropTypes.string,
      tacType: PropTypes.string,
      sacType: PropTypes.string,
    }),
  ]).isRequired,
  allowApproval: PropTypes.bool,
  editURL: PropTypes.string,
  ordersLOA: OrdersLOAShape,
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  neverShow: PropTypes.arrayOf(PropTypes.string),
};

ShipmentDisplay.defaultProps = {
  onChange: () => {},
  shipmentType: SHIPMENT_OPTIONS.HHG,
  allowApproval: true,
  editURL: '',
  ordersLOA: {
    tac: '',
    sac: '',
    ntsTac: '',
    ntsSac: '',
  },
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  neverShow: [],
};

export default ShipmentDisplay;
