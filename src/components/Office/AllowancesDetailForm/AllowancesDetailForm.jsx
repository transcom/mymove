import React from 'react';
import PropTypes from 'prop-types';

import styles from './AllowancesDetailForm.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownInput, CheckboxField } from 'components/form/fields';
import { DropdownArrayOf } from 'types/form';
import { EntitlementShape } from 'types/order';
import { formatWeight } from 'utils/formatters';
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
      >
        <Hint data-testid="proGearWeightHint">
          <p>Max. 2,000 lbs</p>
        </Hint>
      </MaskedTextField>

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
      >
        <Hint data-testid="proGearWeightSpouseHint">
          <p>Max. 500 lbs</p>
        </Hint>
      </MaskedTextField>

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
      <MaskedTextField
        data-testid="sitInput"
        defaultValue="0"
        name="storageInTransit"
        label="Storage in transit (days)"
        id="sitInput"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
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
        <CheckboxField
          id="dependentsAuthorizedInput"
          data-testid="dependentsAuthorizedInput"
          name="dependentsAuthorized"
          label="Dependents authorized"
        />
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
