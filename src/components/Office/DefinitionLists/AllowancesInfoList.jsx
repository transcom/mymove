import React from 'react';
import * as PropTypes from 'prop-types';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatWeight, formatDaysInTransit } from 'shared/formatters';

const AllowancesInfoList = ({ info }) => {
  const titleCase = (input) => {
    if (input && input.length > 0) {
      const friendlyInput = input.toLowerCase().replace('_', ' ').split(' ');
      return friendlyInput
        .map((word) => {
          return word.replace(word[0], word[0].toUpperCase());
        })
        .join(' ');
    }
    return input;
  };
  const friendlyRankDisplay = (rank) => {
    if (rank) {
      const friendlyRank = rank.split('_');
      return `${friendlyRank[0]}-${friendlyRank[1]} ${titleCase(friendlyRank.slice(2).join(' '))}`;
    }
    return rank;
  };

  return (
    <div className={styles.OfficeDefinitionLists}>
      <dl className={descriptionListStyles.descriptionList}>
        <div className={descriptionListStyles.row}>
          <dt>Branch, rank</dt>
          <dd data-testid="branchRank">{`${titleCase(info.branch)}, ${friendlyRankDisplay(info.rank)}`}</dd>
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
          <dt>Storage in transit</dt>
          <dd data-testid="storageInTransit">
            {info.storageInTransit ? formatDaysInTransit(info.storageInTransit) : ''}
          </dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Dependents</dt>
          <dd data-testid="dependents">{info.dependents ? 'Authorized' : 'Unauthorized'}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Pro-gear</dt>
          <dd data-testid="progear">{formatWeight(info.progear)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>Spouse pro-gear</dt>
          <dd data-testid="spouseProgear">{formatWeight(info.spouseProgear)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>RME</dt>
          <dd data-testid="rme">{formatWeight(info.requiredMedicalEquipmentWeight)}</dd>
        </div>
        <div className={descriptionListStyles.row}>
          <dt>OCIE</dt>
          <dd data-testid="ocie">
            {info.organizationalClothingAndIndividualEquipment ? 'Authorized' : 'Unauthorized'}
          </dd>
        </div>
      </dl>
    </div>
  );
};

AllowancesInfoList.propTypes = {
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
};

export default AllowancesInfoList;
