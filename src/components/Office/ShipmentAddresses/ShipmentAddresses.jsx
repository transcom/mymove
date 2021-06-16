import React from 'react';
import classnames from 'classnames';
import { PropTypes } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../types/address';
import { formatAddress } from '../../../utils/shipmentDisplay';
import DataPointGroup from '../../DataPointGroup/index';
import DataPoint from '../../DataPoint/index';

import styles from './ShipmentAddresses.module.scss';

import { shipmentStatuses } from 'constants/shipments';

const ShipmentAddresses = ({
  pickupAddress,
  destinationAddress,
  originDutyStation,
  destinationDutyStation,
  handleDivertShipment,
  shipmentInfo,
}) => {
  return (
    <DataPointGroup className={classnames('maxw-tablet', styles.mtoShipmentAddresses)}>
      <DataPoint
        columnHeaders={[
          'Authorized addresses',
          <div className={styles.rightAlignButtonWrapper}>
            {shipmentInfo.shipmentStatus !== shipmentStatuses.CANCELED && (
              <Button
                type="button"
                onClick={() => handleDivertShipment(shipmentInfo.shipmentID, shipmentInfo.ifMatchEtag)}
                unstyled
              >
                Request diversion
              </Button>
            )}
          </div>,
        ]}
        dataRow={[formatAddress(originDutyStation), formatAddress(destinationDutyStation)]}
        icon={<FontAwesomeIcon icon="arrow-right" />}
      />
      <DataPoint
        columnHeaders={["Customer's addresses", '']}
        dataRow={[formatAddress(pickupAddress), formatAddress(destinationAddress)]}
        icon={<FontAwesomeIcon icon="arrow-right" />}
      />
    </DataPointGroup>
  );
};

ShipmentAddresses.propTypes = {
  pickupAddress: AddressShape,
  destinationAddress: AddressShape,
  originDutyStation: AddressShape,
  destinationDutyStation: AddressShape,
  handleDivertShipment: PropTypes.func.isRequired,
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    ifMatchEtag: PropTypes.string.isRequired,
    shipmentStatus: PropTypes.string.isRequired,
  }).isRequired,
};

ShipmentAddresses.defaultProps = {
  pickupAddress: {},
  destinationAddress: {},
  originDutyStation: {},
  destinationDutyStation: {},
};

export default ShipmentAddresses;
