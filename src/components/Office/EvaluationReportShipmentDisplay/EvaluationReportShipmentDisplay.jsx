import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import styles from './EvaluationReportShipmentDisplay.module.scss';

import ShipmentInfoListSelector from 'components/Office/DefinitionLists/ShipmentInfoListSelector';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape } from 'types/address';
import { AgentShape } from 'types/agent';
import { OrdersLOAShape } from 'types/order';
import { shipmentStatuses } from 'constants/shipments';
import { ShipmentStatusesOneOf } from 'types/shipment';
import { formatAddress, retrieveSAC, retrieveTAC } from 'utils/shipmentDisplay';

const EvaluationReportShipmentDisplay = ({
  shipmentType,
  displayInfo,
  shipmentId,
  allowApproval,
  ordersLOA,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  neverShow,
}) => {
  const containerClasses = classnames(styles.container, { [styles.noIcon]: !allowApproval });
  const [isExpanded, setIsExpanded] = useState(true);
  const tac = retrieveTAC(displayInfo.tacType, ordersLOA);
  const sac = retrieveSAC(displayInfo.sacType, ordersLOA);

  const destinationAddressString = displayInfo.destinationAddress ? formatAddress(displayInfo.destinationAddress) : '';
  const pickupAddressString = displayInfo.pickupAddress ? formatAddress(displayInfo.pickupAddress) : '';

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
          <div className={styles.headingShipmentID}>Shipment ID: #{shipmentId.slice(0, 5).toUpperCase()}</div>
          <FontAwesomeIcon className={styles.icon} icon={expandableIconClasses} onClick={handleExpandClick} />
        </div>
        {isExpanded && displayInfo.shipmentType === SHIPMENT_OPTIONS.NTS && (
          <div className={styles.ntsHeaderText}>
            <div className={styles.ntsHeaderTextField}>Pickup address</div>
            <div className={styles.ntsHeaderTextField}>
              {displayInfo?.storageFacility ? displayInfo.storageFacility.facilityName : ''}
            </div>
          </div>
        )}
        {isExpanded && displayInfo.shipmentType === SHIPMENT_OPTIONS.NTSR && (
          <div className={styles.ntsHeaderText}>
            <div className={styles.ntsHeaderTextField}>
              {displayInfo?.storageFacility ? displayInfo.storageFacility.facilityName : ''}
            </div>
            <div className={styles.ntsHeaderTextField}>Delivery address</div>
          </div>
        )}
        {isExpanded && displayInfo.pickupAddress && displayInfo.destinationAddress && (
          <div className={styles.shipmentAddresses}>
            <div className={styles.shipmentAddressTextFields}>{pickupAddressString}</div>
            <div className={styles.shipmentAddressArrow}>
              <FontAwesomeIcon icon="arrow-right" />
            </div>
            <div className={styles.shipmentAddressTextFields}>{destinationAddressString}</div>
          </div>
        )}
        <ShipmentInfoListSelector
          className={styles.shipmentDisplayInfo}
          shipment={{ ...displayInfo, tac, sac }}
          shipmentType={shipmentType}
          isExpanded={isExpanded}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          showWhenCollapsed={showWhenCollapsed}
          neverShow={neverShow}
          isForEvaluationReport
        />
      </ShipmentContainer>
    </div>
  );
};

EvaluationReportShipmentDisplay.propTypes = {
  shipmentId: PropTypes.string.isRequired,
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
  ordersLOA: OrdersLOAShape,
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  neverShow: PropTypes.arrayOf(PropTypes.string),
};

EvaluationReportShipmentDisplay.defaultProps = {
  shipmentType: SHIPMENT_OPTIONS.HHG,
  allowApproval: true,
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

export default EvaluationReportShipmentDisplay;
