import React from 'react';
import { Field } from 'formik';

import styles from './AllowancesDetailForm.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField';
import CheckboxField from 'components/form/fields/CheckboxField';
import { DropdownInput } from 'components/form/fields';
import { DropdownArrayOf } from 'types/form';
import { EntitlementShape } from 'types/order';
import { formatWeight, formatDaysInTransit } from 'shared/formatters';

const AllowancesDetailForm = ({ entitlements, rankOptions, branchOptions }) => {
  return (
    <div className={styles.AllowancesDetailForm}>
      <MaskedTextField
        defaultValue="0"
        name="proGearWeight"
        label="Pro-gear (lbs)"
        id="proGearWeightInput"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSearator=","
        lazy={false} // immediate masking evaluation
      />
      <MaskedTextField
        defaultValue="0"
        name="proGearWeightSpouse"
        label="Spouse pro-gear (lbs)"
        id="proGearWeightSpouseInput"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSearator=","
        lazy={false} // immediate masking evaluation
      />
      <MaskedTextField
        defaultValue="0"
        name="requiredMedicalEquipmentWeight"
        label="RME estimated weight (lbs)"
        id="rmeInput"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSearator=","
        lazy={false} // immediate masking evaluation
      />
      <DropdownInput name="agency" label="Branch" options={branchOptions} showDropdownPlaceholderText={false} />
      <DropdownInput name="grade" label="Rank" options={rankOptions} showDropdownPlaceholderText={false} />
      <div className={styles.DependentsAuthorized}>
        <CheckboxField
          id="ocieInput"
          name="organizationalClothingAndIndividualEquipment"
          label="OCIE authorized (Army only)"
        />
      </div>
      {/* TODO - Get a bool value to show or hide this field */}
      {/* <MaskedTextField */}
      {/*  defaultValue="0" */}
      {/*  name="authorizedWeight" */}
      {/*  label="Authorized weight" */}
      {/*  id="authorizedWeightInput" */}
      {/*  mask="NUM lbs" // Nested masking imaskjs */}
      {/*  lazy={false} // immediate masking evaluation */}
      {/*  blocks={{ */}
      {/*    // our custom masking key */}
      {/*    NUM: { */}
      {/*      mask: Number, */}
      {/*      thousandsSeparator: ',', */}
      {/*      scale: 0, // whole numbers */}
      {/*      signed: false, // positive numbers */}
      {/*    }, */}
      {/*  }} */}
      {/* /> */}
      <dl>
        <dt>Authorized weight</dt>
        <dd data-testid="authorizedWeight">{formatWeight(entitlements.authorizedWeight)}</dd>
        <dt>Weight allowance</dt>
        <dd data-testid="weightAllowance">{formatWeight(entitlements.totalWeight)}</dd>
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
