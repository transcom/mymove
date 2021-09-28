import React, { useState } from 'react';
import classnames from 'classnames';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';

import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';
import ReviewSITExtensionsModal from '../ReviewSITExtensionModal/ReviewSITExtensionModal';
import { SITExtensionShape } from '../../../types/sitExtensions';

import styles from './ShipmentSITExtensions.module.scss';

import { sitExtensionReasons, SIT_EXTENSION_STATUS } from 'constants/sitExtensions';
import { formatDateFromIso } from 'shared/formatters';

const ShipmentSITExtensions = (props) => {
  const { sitExtensions, handleReviewSITExtension } = props;
  const [isReviewSITExtensionModalVisible, setisReviewSITExtensionModalVisible] = useState(false);
  const reviewSITExtensionSubmit = (sitExtensionID, formValues) => {
    setisReviewSITExtensionModalVisible(false);
    handleReviewSITExtension(sitExtensionID, formValues);
  };

  const pendingSITExtension = sitExtensions.find((se) => se.status === SIT_EXTENSION_STATUS.PENDING);
  const showModal = isReviewSITExtensionModalVisible && pendingSITExtension !== undefined;

  const mappedSITExtensionList = sitExtensions.map((sitExt) => {
    return (
      <dl key={sitExt.id}>
        {sitExt.approvedDays > 0 && (
          <div>
            <dt>{sitExt.approvedDays} days added</dt>
            <dd>on {formatDateFromIso(sitExt.decisionDate, 'DD MMM YYYY')}</dd>
          </div>
        )}
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
    <DataTableWrapper className={classnames('maxw-tablet', styles.mtoShipmentSITExtensions)} testID="sitExtensions">
      <div className={styles.title}>
        <p>SIT (STORAGE IN TRANSIT){pendingSITExtension && <Tag>Extension requested</Tag>}</p>
        {pendingSITExtension && (
          <p>
            <Button type="button" onClick={() => setisReviewSITExtensionModalVisible(true)} unstyled>
              View request
            </Button>
          </p>
        )}
      </div>

      <DataTable columnHeaders={['SIT extensions']} dataRow={[mappedSITExtensionList]} />
      {showModal && (
        <ReviewSITExtensionsModal
          onClose={() => setisReviewSITExtensionModalVisible(false)}
          onSubmit={reviewSITExtensionSubmit}
          sitExtension={pendingSITExtension}
        />
      )}
    </DataTableWrapper>
  );
};

ShipmentSITExtensions.propTypes = {
  sitExtensions: PropTypes.arrayOf(SITExtensionShape).isRequired,
  handleReviewSITExtension: PropTypes.func.isRequired,
};

ShipmentSITExtensions.defaultProps = {};

export default ShipmentSITExtensions;
