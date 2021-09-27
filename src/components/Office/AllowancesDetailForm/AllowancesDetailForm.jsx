import React from 'react';
import PropTypes from 'prop-types';

import styles from './AllowancesDetailForm.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField';
import { DropdownInput, CheckboxField } from 'components/form/fields';
import { DropdownArrayOf } from 'types/form';
import { EntitlementShape } from 'types/order';
import { formatWeight } from 'shared/formatters';
import Hint from 'components/Hint';

const AllowancesDetailForm = ({ header, entitlements, rankOptions, branchOptions, editableAuthorizedWeight }) => {
  return (
    <div className={styles.AllowancesDetailForm}>
      {header && <h3 data-testid="header">{header}</h3>}
      <MaskedTextField
        data-testid="proGearWeightInput"
        defaultValue="0"
        name="proGearWeight"
        label="Pro-gear (lbs)"
        id="proGearWeightInput"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
        formGroupClassName={styles.fieldWithHint}
      />
      <Hint data-testid="proGearWeightHint">
        <p>Max. 2,000 lbs</p>
      </Hint>

      <MaskedTextField
        data-testid="proGearWeightSpouseInput"
        defaultValue="0"
        name="proGearWeightSpouse"
        label="Spouse pro-gear (lbs)"
        id="proGearWeightSpouseInput"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
        formGroupClassName={styles.fieldWithHint}
      />
      <Hint data-testid="proGearWeightSpouseHint">
        <p>Max. 500 lbs</p>
      </Hint>

      <MaskedTextField
        data-testid="rmeInput"
        defaultValue="0"
        name="requiredMedicalEquipmentWeight"
        label="RME estimated weight (lbs)"
        id="rmeInput"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
      />
      <DropdownInput
        data-testid="branchInput"
        name="agency"
        label="Branch"
        options={branchOptions}
        showDropdownPlaceholderText={false}
      />
      <DropdownInput
        data-testid="rankInput"
        name="grade"
        label="Rank"
        options={rankOptions}
        showDropdownPlaceholderText={false}
      />
      <div className={styles.wrappedCheckbox}>
        <CheckboxField
          data-testid="ocieInput"
          id="ocieInput"
          name="organizationalClothingAndIndividualEquipment"
          label="OCIE authorized (Army only)"
        />
      </div>

      {editableAuthorizedWeight && (
        <MaskedTextField
          data-testid="authorizedWeightInput"
          defaultValue="0"
          name="authorizedWeight"
          label="Authorized weight (lbs)"
          id="authorizedWeightInput"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          signed={false} // disallow negative
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
        />
      )}

      <dl>
        {!editableAuthorizedWeight && (
          <>
            <dt>Authorized weight</dt>
            <dd data-testid="authorizedWeight">{formatWeight(entitlements.authorizedWeight)}</dd>
          </>
        )}
        <dt>Weight allowance</dt>
        <dd data-testid="weightAllowance">{formatWeight(entitlements.totalWeight)}</dd>
      </dl>
      <div className={styles.wrappedCheckbox}>
        <CheckboxField id="dependentsAuthorizedInput" name="dependentsAuthorized" label="Dependents authorized" />
      </div>
    </div>
  );
};

AllowancesDetailForm.propTypes = {
  entitlements: EntitlementShape.isRequired,
  rankOptions: DropdownArrayOf.isRequired,
  branchOptions: DropdownArrayOf.isRequired,
  header: PropTypes.string,
  editableAuthorizedWeight: PropTypes.bool,
};

AllowancesDetailForm.defaultProps = {
  header: null,
  editableAuthorizedWeight: false,
};

export default AllowancesDetailForm;
