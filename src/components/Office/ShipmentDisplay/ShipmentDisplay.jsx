import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import { useHistory } from 'react-router-dom';
import { Checkbox, Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import { EditButton } from 'components/form/IconButtons';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import ShipmentInfoList from 'components/Office/DefinitionLists/ShipmentInfoList';
import NTSRShipmentInfoList from 'components/Office/DefinitionLists/NTSRShipmentInfoList';
import styles from 'components/Office/ShipmentDisplay/ShipmentDisplay.module.scss';
import { LOA_TYPE, SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape } from 'types/address';
import { shipmentStatuses } from 'constants/shipments';
import { ShipmentStatusesOneOf } from 'types/shipment';
import { AgentShape } from 'types/agent';

const ShipmentDisplay = ({
  shipmentType,
  displayInfo,
  onChange,
  shipmentId,
  isSubmitted,
  showIcon,
  editURL,
  ordersLOA,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
}) => {
  const history = useHistory();
  const containerClasses = classnames(styles.container, { [styles.noIcon]: !showIcon });
  const [isExpanded, setIsExpanded] = useState(false);
  let infoList;
  let tac;
  switch (displayInfo.tacType) {
    case LOA_TYPE.HHG:
      tac = ordersLOA.tac;
      break;
    case LOA_TYPE.NTS:
      tac = ordersLOA.ntsTAC;
      break;
    default:
      tac = ordersLOA.tac;
  }

  let sac;
  switch (displayInfo.sacType) {
    case LOA_TYPE.HHG:
      sac = ordersLOA.sac;
      break;
    case LOA_TYPE.NTS:
      sac = ordersLOA.ntsSAC;
      break;
    default:
      sac = ordersLOA.sac;
  }

  const setDisplayInfo = () => {
    switch (shipmentType) {
      case SHIPMENT_OPTIONS.HHG:
        infoList = (
          <ShipmentInfoList className={styles.shipmentDisplayInfo} shipment={displayInfo} shipmentType={shipmentType} />
        );
        break;
      case SHIPMENT_OPTIONS.NTSR:
        infoList = (
          <NTSRShipmentInfoList
            className={styles.shipmentDisplayInfo}
            shipment={{ ...displayInfo, tac, sac }}
            isExpanded={isExpanded}
            warnIfMissing={warnIfMissing}
            errorIfMissing={errorIfMissing}
            showWhenCollapsed={showWhenCollapsed}
          />
        );
        break;
      default:
        infoList = (
          <ShipmentInfoList
            className={styles.shipmentDisplayInfo}
            shipment={displayInfo}
            shipmentType={shipmentType}
            isExpanded={isExpanded}
          />
        );
    }
  };
  setDisplayInfo();

  const handleExpandClick = () => {
    setIsExpanded((prev) => !prev);
    setDisplayInfo();
  };
  const expandableIconClasses = classnames({
    'chevron-up': isExpanded,
    'chevron-down': !isExpanded,
  });

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
            {displayInfo.shipmentStatus === shipmentStatuses.DIVERSION_REQUESTED && <Tag>diversion requested</Tag>}
            {displayInfo.shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED && (
              <Tag>cancellation requested</Tag>
            )}
          </div>

          <FontAwesomeIcon className={styles.icon} icon={expandableIconClasses} onClick={handleExpandClick} />
        </div>
        {infoList}
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
  ]),
  displayInfo: PropTypes.shape({
    heading: PropTypes.string.isRequired,
    isDiversion: PropTypes.bool,
    shipmentStatus: ShipmentStatusesOneOf,
    requestedPickupDate: PropTypes.string,
    pickupAddress: AddressShape,
    secondaryPickupAddress: AddressShape,
    destinationAddress: AddressShape,
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
  }).isRequired,
  showIcon: PropTypes.bool,
  editURL: PropTypes.string,
  ordersLOA: PropTypes.shape({
    tac: PropTypes.string,
    sac: PropTypes.string,
    ntsTAC: PropTypes.string,
    ntsSAC: PropTypes.string,
  }),
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
};

ShipmentDisplay.defaultProps = {
  onChange: () => {},
  shipmentType: SHIPMENT_OPTIONS.HHG,
  showIcon: true,
  editURL: '',
  ordersLOA: {
    tac: '',
    sac: '',
    ntsTAC: '',
    ntsSAC: '',
  },
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
};

export default ShipmentDisplay;
