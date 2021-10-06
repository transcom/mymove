import React from 'react';
import classnames from 'classnames';
import * as PropTypes from 'prop-types';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';
import styles from '../ShipmentDetails/ShipmentDetails.module.scss';

const ShipmentRemarks = ({ title, remarks }) => {
  return (
    <DataTableWrapper className={classnames('maxw-tablet', 'table--data-point-group', styles.ShipmentRemarks)}>
      <DataTable columnHeaders={[title]} dataRow={[remarks]} />
    </DataTableWrapper>
  );
};

ShipmentRemarks.propTypes = {
  title: PropTypes.string.isRequired,
  remarks: PropTypes.string.isRequired,
};

export default ShipmentRemarks;
