import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatWeight } from 'utils/formatters';
import friendlyBranchRank from 'utils/branchRankFormatters';

const AllowancesList = ({ info, showVisualCues }) => {
  const visualCuesStyle = classNames(descriptionListStyles.row, {
    [`${descriptionListStyles.rowWithVisualCue}`]: showVisualCues,
  });

  return (
    <div className={styles.OfficeDefinitionLists}>
      <dl className={descriptionListStyles.descriptionList}>
        <div className={descriptionListStyles.row}>
          <dt>Branch, rank</dt>
          <dd data-testid="branchRank">{friendlyBranchRank(info.branch, info.rank)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Weight allowance</dt>
          <dd data-testid="weightAllowance">{formatWeight(info.weightAllowance)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Authorized weight</dt>
          <dd data-testid="authorizedWeight">{formatWeight(info.authorizedWeight)}</dd>
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
          <dt>RME</dt>
          <dd data-testid="rme">{formatWeight(info.requiredMedicalEquipmentWeight)}</dd>
        </div>
        <div className={visualCuesStyle}>
          <dt>OCIE</dt>
          <dd data-testid="ocie">
            {info.organizationalClothingAndIndividualEquipment ? 'Authorized' : 'Unauthorized'}
          </dd>
        </div>
      </dl>
    </div>
  );
};

AllowancesList.propTypes = {
  info: PropTypes.shape({
    branch: PropTypes.string,
    rank: PropTypes.string,
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
