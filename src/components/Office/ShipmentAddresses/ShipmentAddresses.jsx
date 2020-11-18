import React from 'react';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowRight } from '@fortawesome/free-solid-svg-icons';

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
        Icon={<FontAwesomeIcon icon={faArrowRight} />}
      />
      <DataPoint
        columnHeaders={["Customer's addresses", '']}
        dataRow={[formatAddress(pickupAddress), formatAddress(destinationAddress)]}
        Icon={<FontAwesomeIcon icon={faArrowRight} />}
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
