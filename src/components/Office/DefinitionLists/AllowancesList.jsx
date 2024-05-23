import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatWeight } from 'utils/formatters';
import { ORDERS_BRANCH_OPTIONS } from 'constants/orders';

const AllowancesList = ({ info, showVisualCues }) => {
  const visualCuesStyle = classNames(descriptionListStyles.row, {
    [`${descriptionListStyles.rowWithVisualCue}`]: showVisualCues,
  });

  return (
    <div className={styles.OfficeDefinitionLists}>
      <dl className={descriptionListStyles.descriptionList}>
        <div className={descriptionListStyles.row}>
          <dt>Branch</dt>
          <dd data-testid="branch">{info.branch ? ORDERS_BRANCH_OPTIONS[info.branch] : ''}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Weight allowance</dt>
          <dd data-testid="weightAllowance">{formatWeight(info.authorizedWeight)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Storage in transit (SIT)</dt>
          <dd data-testid="storageInTransit">{info.storageInTransit} days</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Dependents</dt>
          <dd data-testid="dependents">{info.dependents ? 'Authorized' : 'Unauthorized'}</dd>
        </div>
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
    weightAllowance: PropTypes.number,
    authorizedWeight: PropTypes.number,
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
