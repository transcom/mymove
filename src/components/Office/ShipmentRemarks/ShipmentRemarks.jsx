import React from 'react';
import classnames from 'classnames';
import * as PropTypes from 'prop-types';

import DataPointGroup from '../../DataPointGroup/index';
import DataPoint from '../../DataPoint/index';
import styles from '../ShipmentDetails/ShipmentDetails.module.scss';

const ShipmentRemarks = ({ title, remarks }) => {
  return (
    <DataPointGroup className={classnames('maxw-tablet', styles.ShipmentRemarks)}>
      <DataPoint columnHeaders={[title]} dataRow={[remarks]} />
    </DataPointGroup>
  );
};

ShipmentRemarks.propTypes = {
  title: PropTypes.string.isRequired,
  remarks: PropTypes.string.isRequired,
};

export default ShipmentRemarks;
