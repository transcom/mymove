import React from 'react';
import * as PropTypes from 'prop-types';

import DataPointGroup from '../../DataPointGroup/index';
import DataPoint from '../../DataPoint/index';

const ShipmentRemarks = ({ title, remarks }) => {
  return (
    <DataPointGroup className="maxw-tablet">
      <DataPoint columnHeaders={[title]} dataRow={[remarks]} />
    </DataPointGroup>
  );
};

ShipmentRemarks.propTypes = {
  title: PropTypes.string.isRequired,
  remarks: PropTypes.string.isRequired,
};

export default ShipmentRemarks;
