import React from 'react';
import classnames from 'classnames';
import { PropTypes } from 'prop-types';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';
import { SITExtensionShape } from '../../../types/sitExtensions';

import styles from './ShipmentSITExtensions.module.scss';

import { sitExtensionReasons } from 'constants/sitExtensions';
import { formatDateFromIso } from 'shared/formatters';

const ShipmentSITExtensions = (props) => {
  const { sitExtensions } = props;

  const mappedSITExtensionList = sitExtensions.map((sitExt) => {
    return (
      <dl key={sitExt.id}>
        <div>
          <dt>{sitExt.approvedDays} days added</dt>
          <dd>on {formatDateFromIso(sitExt.decisionDate, 'DD MMM YYYY')}</dd>
        </div>
        <div>
          <dt>Reason:</dt>
          <dd>{sitExtensionReasons[sitExt.requestReason]}</dd>
        </div>
        {sitExt.contractorRemarks && (
          <div>
            <dt>Contractor remarks:</dt>
            <dd>{sitExt.contractorRemarks}</dd>
          </div>
        )}
        {sitExt.officeRemarks && (
          <div>
            <dt>Office remarks:</dt>
            <dd>{sitExt.officeRemarks}</dd>
          </div>
        )}
      </dl>
    );
  });

  return (
    <DataTableWrapper
      className={classnames('maxw-tablet', styles.mtoShipmentSITExtensions)}
      data-testid="sitExtensions"
    >
      <p>SIT (STORAGE IN TRANSIT)</p>
      <DataTable columnHeaders={['SIT extensions']} dataRow={[mappedSITExtensionList]} />
    </DataTableWrapper>
  );
};

ShipmentSITExtensions.propTypes = {
  sitExtensions: PropTypes.arrayOf(SITExtensionShape).isRequired,
};

ShipmentSITExtensions.defaultProps = {};

export default ShipmentSITExtensions;
