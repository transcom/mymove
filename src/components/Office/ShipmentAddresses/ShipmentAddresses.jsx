import React from 'react';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { AddressShape } from '../../../types/address';
import { formatAddress } from '../../../utils/shipmentDisplay';
import DataPointGroup from '../../DataPointGroup/index';
import DataPoint from '../../DataPoint/index';

import styles from './ShipmentAddresses.module.scss';

const ShipmentAddresses = ({ pickupAddress, destinationAddress, originDutyStation, destinationDutyStation }) => {
  return (
    <DataPointGroup className={classnames('maxw-tablet', styles.mtoShipmentAddresses)}>
      <DataPoint
        columnHeaders={['Authorized addresses', '']}
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
};

ShipmentAddresses.defaultProps = {
  pickupAddress: {},
  destinationAddress: {},
  originDutyStation: {},
  destinationDutyStation: {},
};

export default ShipmentAddresses;
