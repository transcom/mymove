import React from 'react';
import { Field } from 'formik';

import styles from './AllowancesDetailForm.module.scss';

import { TextMaskedInput } from 'components/form/fields/TextInput';
import { DropdownInput } from 'components/form/fields';
import { DropdownArrayOf } from 'types/form';
import { EntitlementShape } from 'types/order';
import { formatWeight, formatDaysInTransit } from 'shared/formatters';

const AllowancesDetailForm = ({ entitlements, rankOptions, branchOptions }) => {
  return (
    <div className={styles.AllowancesDetailForm}>
      <DropdownInput name="agency" label="Branch" options={branchOptions} showDropdownPlaceholderText={false} />
      <DropdownInput name="grade" label="Rank" options={rankOptions} showDropdownPlaceholderText={false} />
      <TextMaskedInput
        defaultValue="0"
        name="authorizedWeight"
        label="Authorized weight"
        id="authorizedWeightInput"
        mask="NUM lbs" // Nested masking imaskjs
        lazy={false} // immediate masking evaluation
        blocks={{
          // our custom masking key
          NUM: {
            mask: Number,
            thousandsSeparator: ',',
            scale: 0, // whole numbers
            signed: false, // positive numbers
          },
        }}
      />
      <dl>
        <dt>Weight allowance</dt>
        <dd data-testid="weightAllowance">{formatWeight(entitlements.totalWeight)}</dd>
        <dt>Pro-gear</dt>
        <dd data-testid="proGearWeight">{formatWeight(entitlements.proGearWeight)}</dd>
        <dt>Spouse pro-gear</dt>
        <dd data-testid="spouseProGearWeight">{formatWeight(entitlements.proGearWeightSpouse)}</dd>
        <dt>Storage in-transit</dt>
        <dd data-testid="storageInTransit">{formatDaysInTransit(entitlements.storageInTransit)}</dd>
      </dl>
      <div className={styles.DependentsAuthorized}>
        <Field type="checkbox" name="dependentsAuthorized" />
        <label htmlFor="dependentsAuthorized"> Dependents Authorized</label>
      </div>
    </div>
  );
};

AllowancesDetailForm.propTypes = {
  entitlements: EntitlementShape.isRequired,
  rankOptions: DropdownArrayOf.isRequired,
  branchOptions: DropdownArrayOf.isRequired,
};

export default AllowancesDetailForm;
