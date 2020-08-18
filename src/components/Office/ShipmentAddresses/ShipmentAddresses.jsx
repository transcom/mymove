import React from 'react';

import { AddressShape } from '../../../types/address';
import formatAddress from '../../../utils/shipmentDisplay';
import DataPointGroup from '../../DataPointGroup/index';
import DataPoint from '../../DataPoint/index';

import { ReactComponent as ArrowRight } from 'shared/icon/arrow-right.svg';

const ShipmentAddresses = ({ pickupAddress, destinationAddress, originDutyStation, destinationDutyStation }) => {
  return (
    <DataPointGroup className="maxw-tablet">
      <DataPoint
        columnHeaders={['Authorized addresses', '']}
        dataRow={[formatAddress(originDutyStation), formatAddress(destinationDutyStation)]}
        Icon={ArrowRight}
      />
      <DataPoint
        columnHeaders={["Customer's addresses", '']}
        dataRow={[formatAddress(pickupAddress), formatAddress(destinationAddress)]}
        Icon={ArrowRight}
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
