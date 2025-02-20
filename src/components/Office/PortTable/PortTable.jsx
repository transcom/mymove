import React from 'react';
import classnames from 'classnames';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';
import styles from '../ShipmentAddresses/ShipmentAddresses.module.scss';

import { formatPortInfo } from 'utils/formatters';

const PortTable = ({ poeLocation, podLocation }) => {
  return (
    <DataTableWrapper className={classnames('maxw-tablet', 'table--data-point-group', styles.mtoShipmentAddresses)}>
      <DataTable
        columnHeaders={['Port of Embarkation', 'Port of Debarkation']}
        dataRow={[formatPortInfo(poeLocation), formatPortInfo(podLocation)]}
      />
    </DataTableWrapper>
  );
};

export default PortTable;
