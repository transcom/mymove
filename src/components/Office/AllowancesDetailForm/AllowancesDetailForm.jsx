import React from 'react';

import styles from './AllowancesDetailForm.module.scss';

import { EntitlementShape } from 'types/moveOrder';
import { formatWeight, formatDaysInTransit } from 'shared/formatters';

const AllowancesDetailForm = ({ entitlements }) => {
  return (
    <div className={styles.AllowancesDetailForm}>
      <dl>
        <dt>Weight allowance</dt>
        <dd>{formatWeight(entitlements.totalWeight)}</dd>
        <dt>Pro-gear</dt>
        <dd>{formatWeight(entitlements.proGearWeight)}</dd>
        <dt>Spouse pro-gear</dt>
        <dd>{formatWeight(entitlements.proGearWeightSpouse)}</dd>
        <dt>Storage in-transit</dt>
        <dd>{formatDaysInTransit(entitlements.storageInTransit)}</dd>
      </dl>
    </div>
  );
};

AllowancesDetailForm.propTypes = {
  entitlements: EntitlementShape.isRequired,
};

export default AllowancesDetailForm;
