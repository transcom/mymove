import React from 'react';
import classnames from 'classnames';
import PropTypes from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';

import styles from './ShipmentWeightDetails.module.scss';

import returnLowestValue from 'utils/returnLowestValue';
import { formatWeight } from 'utils/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const ShipmentWeightDetails = ({ estimatedWeight, actualWeight, shipmentInfo, handleRequestReweighModal }) => {
  const lowestWeight = returnLowestValue(actualWeight, shipmentInfo.reweighWeight);
  const reweighHeader = (
    <div className={styles.shipmentWeight}>
      <span>Shipment weight</span>
      {!shipmentInfo.reweighID && (
        <div className={styles.rightAlignButtonWrapper}>
          <Restricted to={permissionTypes.createReweighRequest}>
            <Button type="button" onClick={() => handleRequestReweighModal(shipmentInfo)} unstyled>
              Request reweigh
            </Button>
          </Restricted>
        </div>
      )}
      {shipmentInfo.reweighID && !shipmentInfo.reweighWeight && <Tag>reweigh requested</Tag>}
      {shipmentInfo.reweighWeight && <Tag>reweighed</Tag>}
    </div>
  );
  return (
    <div className={classnames('maxw-tablet', styles.ShipmentWeightDetails)}>
      {shipmentInfo.shipmentType !== SHIPMENT_OPTIONS.NTSR && (
        <DataTableWrapper className={classnames('maxw-mobile', 'table--data-point-group')}>
          <DataTable
            columnHeaders={['Estimated weight']}
            dataRow={estimatedWeight ? [formatWeight(estimatedWeight)] : ['']}
          />
        </DataTableWrapper>
      )}
      <DataTableWrapper
        className={classnames('table--data-point-group', {
          'maxw-mobile': shipmentInfo.shipmentType !== SHIPMENT_OPTIONS.NTSR,
        })}
      >
        <DataTable columnHeaders={[reweighHeader]} dataRow={lowestWeight ? [formatWeight(lowestWeight)] : ['']} />
      </DataTableWrapper>
    </div>
  );
};

ShipmentWeightDetails.propTypes = {
  estimatedWeight: PropTypes.number,
  actualWeight: PropTypes.number,
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string,
    ifMatchEtag: PropTypes.string,
    reweighID: PropTypes.string,
    reweighWeight: PropTypes.number,
    shipmentType: ShipmentOptionsOneOf.isRequired,
  }).isRequired,
  handleRequestReweighModal: PropTypes.func.isRequired,
};

ShipmentWeightDetails.defaultProps = {
  estimatedWeight: null,
  actualWeight: null,
};

export default ShipmentWeightDetails;
