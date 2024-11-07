import React, { useState, useEffect } from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';
import { FEATURE_FLAG_KEYS, DEFAULT_EMPTY_VALUE } from '../../../shared/constants';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatWeight } from 'utils/formatters';
import { ORDERS_BRANCH_OPTIONS } from 'constants/orders';

const AllowancesList = ({ info, showVisualCues }) => {
  const [enableUB, setEnableUB] = useState(false);

  const visualCuesStyle = classNames(descriptionListStyles.row, {
    [`${descriptionListStyles.rowWithVisualCue}`]: showVisualCues,
  });

  useEffect(() => {
    const checkUBFeatureFlag = async () => {
      const enabled = await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.UNACCOMPANIED_BAGGAGE);
      if (enabled) {
        setEnableUB(true);
      }
    };
    checkUBFeatureFlag();
  }, []);

  return (
    <div className={styles.OfficeDefinitionLists}>
      <dl className={descriptionListStyles.descriptionList}>
        <div className={descriptionListStyles.row}>
          <dt>Branch</dt>
          <dd data-testid="branch">{info.branch ? ORDERS_BRANCH_OPTIONS[info.branch] : ''}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Weight allowance</dt>
          <dd data-testid="weightAllowance">{formatWeight(info.totalWeight)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Storage in transit (SIT)</dt>
          <dd data-testid="storageInTransit">{info.storageInTransit} days</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Dependents</dt>
          <dd data-testid="dependents">{info.dependents ? 'Authorized' : 'Unauthorized'}</dd>
        </div>
        {/* Begin OCONUS fields */}
        {/* As these fields are grouped together and only apply to OCONUS orders
        They will all be NULL for CONUS orders. If one of these fields are present,
        it will be safe to assume it is an OCONUS order. With this, if one field is present
        we show all three. Otherwise, we show none */}
        {/* Wrap in FF */}
        {enableUB && (info?.accompaniedTour || info?.dependentsTwelveAndOver || info?.dependentsUnderTwelve) && (
          <>
            <div className={descriptionListStyles.row}>
              <dt>Accompanied tour</dt>
              <dd data-testid="ordersAccompaniedTour">{info.accompaniedTour ? 'Yes' : 'No'}</dd>
            </div>
            <div className={descriptionListStyles.row}>
              <dt>Dependents under age 12</dt>
              <dd data-testid="ordersDependentsUnderTwelve">
                {info.dependentsUnderTwelve ? info.dependentsUnderTwelve : DEFAULT_EMPTY_VALUE}
              </dd>
            </div>
            <div className={descriptionListStyles.row}>
              <dt>Dependents over age 12</dt>
              <dd data-testid="ordersDependentsTwelveAndOver">
                {info.dependentsTwelveAndOver ? info.dependentsTwelveAndOver : DEFAULT_EMPTY_VALUE}
              </dd>
            </div>
          </>
        )}
        {/* End OCONUS fields */}
        <div className={visualCuesStyle}>
          <dt>Pro-gear</dt>
          <dd data-testid="progear">{formatWeight(info.progear)}</dd>
        </div>
        <div className={visualCuesStyle}>
          <dt>Spouse pro-gear</dt>
          <dd data-testid="spouseProgear">{formatWeight(info.spouseProgear)}</dd>
        </div>
        <div className={visualCuesStyle}>
          <dt>Required medical equipment</dt>
          <dd data-testid="rme">{formatWeight(info.requiredMedicalEquipmentWeight)}</dd>
        </div>
        <div className={visualCuesStyle}>
          <dt>OCIE</dt>
          <dd data-testid="ocie">
            {info.organizationalClothingAndIndividualEquipment ? 'Authorized' : 'Unauthorized'}
          </dd>
        </div>
        <div className={visualCuesStyle}>
          <dt>Gun Safe</dt>
          <dd data-testid="gunSafe"> {info.gunSafe ? 'Authorized' : 'Unauthorized'} </dd>
        </div>
      </dl>
    </div>
  );
};

AllowancesList.propTypes = {
  info: PropTypes.shape({
    branch: PropTypes.string,
    grade: PropTypes.string,
    totalWeight: PropTypes.string,
    progear: PropTypes.number,
    spouseProgear: PropTypes.number,
    storageInTransit: PropTypes.number,
    dependents: PropTypes.bool,
    requiredMedicalEquipmentWeight: PropTypes.number,
    organizationalClothingAndIndividualEquipment: PropTypes.bool,
  }).isRequired,
  showVisualCues: PropTypes.bool,
};

AllowancesList.defaultProps = {
  showVisualCues: false,
};

export default AllowancesList;
