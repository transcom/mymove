import React from 'react';

import { AddressShape } from '../../../types/address';
import { DutyStationShape } from '../../../types/dutyStation';
import formatAddress from '../../../utils/shipmentDisplay';
import DataPointGroup from '../../DataPointGroup/index';
import DataPoint from '../../DataPoint/index';

import { ReactComponent as ArrowRight } from 'shared/icon/arrow-right.svg';

const ShipmentAddresses = ({ pickupAddress, destinationAddress, originDutyStation, destinationDutyStation }) => {
  return (
    <DataPointGroup>
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
  originDutyStation: DutyStationShape,
  destinationDutyStation: DutyStationShape,
};

ShipmentAddresses.defaultProps = {
  pickupAddress: {},
  destinationAddress: {},
  originDutyStation: {},
  destinationDutyStation: {},
};

export default ShipmentAddresses;
