import React from 'react';
import PropTypes from 'prop-types';

import DataPointGroup from '../../DataPointGroup/index';
import DataPoint from '../../DataPoint/index';

import { formatWeight } from 'shared/formatters';

const ShipmentWeightDetails = ({ estimatedWeight, actualWeight }) => {
  const headers = ['Estimated weight', 'Actual weight'];
  const row = [estimatedWeight ? formatWeight(estimatedWeight) : '', actualWeight ? formatWeight(actualWeight) : ''];
  return (
    <DataPointGroup className="maxw-tablet">
      <DataPoint columnHeaders={headers} dataRow={row} />
    </DataPointGroup>
  );
};

ShipmentWeightDetails.propTypes = {
  estimatedWeight: PropTypes.number,
  actualWeight: PropTypes.number,
};

ShipmentWeightDetails.defaultProps = {
  estimatedWeight: null,
  actualWeight: null,
};

export default ShipmentWeightDetails;
