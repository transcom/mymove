import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';
import { FEATURE_FLAG_KEYS } from '../../../shared/constants';

import styles from './AllowancesDetailForm.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownInput, CheckboxField } from 'components/form/fields';
import { DropdownArrayOf } from 'types/form';
import { EntitlementShape } from 'types/order';
import { formatWeight } from 'utils/formatters';
import Hint from 'components/Hint';

const AllowancesDetailForm = ({ header, entitlements, branchOptions, formIsDisabled }) => {
  const [enableUB, setEnableUB] = useState(false);
  const renderOconusFields = !!(
    entitlements?.accompaniedTour ||
    entitlements?.dependentsTwelveAndOver ||
    entitlements?.dependentsUnderTwelve
  );
  const [isAdminWeightLocationChecked, setIsAdminWeightLocationChecked] = useState(entitlements?.weightRestriction > 0);
  useEffect(() => {
    // Functional component version of "componentDidMount"
    // By leaving the dependency array empty this will only run once
    const checkUBFeatureFlag = async () => {
      const enabled = await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.UNACCOMPANIED_BAGGAGE);
      if (enabled) {
        setEnableUB(true);
      }
    };
    checkUBFeatureFlag();
  }, []);

  useEffect(() => {
    if (!isAdminWeightLocationChecked) {
      // Find the weight restriction input and reset its value to 0
      const weightRestrictionInput = document.getElementById('weightRestrictionId');
      if (weightRestrictionInput) {
        weightRestrictionInput.value = '0';
      }
    }
  }, [isAdminWeightLocationChecked]);

  const handleAdminWeightLocationChange = (e) => {
    setIsAdminWeightLocationChecked(e.target.checked);
    if (!e.target.checked) {
      // Find the weight restriction input and update both DOM and form state
      const weightRestrictionInput = document.querySelector('input[name="weightRestriction"]');
      if (weightRestrictionInput) {
        weightRestrictionInput.value = '0';
        // Create and dispatch both input and change events
        const inputEvent = new Event('input', { bubbles: true });
        const changeEvent = new Event('change', { bubbles: true });
        weightRestrictionInput.dispatchEvent(inputEvent);
        weightRestrictionInput.dispatchEvent(changeEvent);
      }
    }
  };

  return (
    <div className={styles.AllowancesDetailForm}>
      {header && <h3 data-testid="header">{header}</h3>}
      {enableUB && renderOconusFields && (
        <>
          <MaskedTextField
            data-testid="dependentsUnderTwelveInput"
            defaultValue="0"
            name="dependentsUnderTwelve"
            label="Number of dependents under the age of 12"
            id="dependentsUnderTwelveInput"
            mask={Number}
            scale={0}
            signed={false}
            thousandsSeparator=","
            lazy={false}
            isDisabled={formIsDisabled}
          />

          <MaskedTextField
            data-testid="dependentsTwelveAndOverInput"
            defaultValue="0"
            name="dependentsTwelveAndOver"
            label="Number of dependents of the age 12 or over"
            id="dependentsTwelveAndOverInput"
            mask={Number}
            scale={0}
            signed={false}
            thousandsSeparator=","
            lazy={false}
            isDisabled={formIsDisabled}
          />
          <div className={styles.wrappedCheckbox}>
            <CheckboxField
              id="accompaniedTourInput"
              data-testid="accompaniedTourInput"
              name="accompaniedTour"
              label="Accompanied tour"
              isDisabled={formIsDisabled}
            />
          </div>
        </>
      )}

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
        isDisabled={formIsDisabled}
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
        isDisabled={formIsDisabled}
      >
        <Hint data-testid="proGearWeightSpouseHint">
          <p>Max. 500 lbs</p>
        </Hint>
      </MaskedTextField>

      <MaskedTextField
        data-testid="rmeInput"
        defaultValue="0"
        name="requiredMedicalEquipmentWeight"
        label="Required medical equipment estimated weight (lbs)"
        id="rmeInput"
        mask={Number}
        scale={0} // digits after point, 0 for integers
        signed={false} // disallow negative
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
        isDisabled={formIsDisabled}
      />
      <DropdownInput
        data-testid="branchInput"
        name="agency"
        label="Branch"
        options={branchOptions}
        showDropdownPlaceholderText={false}
        isDisabled={formIsDisabled}
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
        isDisabled={formIsDisabled}
      />
      <dl>
        <dt>Standard weight allowance</dt>
        <dd data-testid="weightAllowance">{formatWeight(entitlements.totalWeight)}</dd>
      </dl>
      <div className={styles.wrappedCheckbox}>
        <CheckboxField
          data-testid="ocieInput"
          id="ocieInput"
          name="organizationalClothingAndIndividualEquipment"
          label="OCIE authorized (Army only)"
          isDisabled={formIsDisabled}
        />
      </div>
      <div className={styles.wrappedCheckbox}>
        <CheckboxField
          data-testid="gunSafeInput"
          id="gunSafeInput"
          name="gunSafe"
          label="Gun safe authorized"
          isDisabled={formIsDisabled}
        />
      </div>
      <div className={styles.wrappedCheckbox}>
        <CheckboxField
          data-testid="adminWeightLocation"
          id="adminWeightLocation"
          name="adminRestrictedWeightLocation"
          label="Admin restricted weight location"
          isDisabled={formIsDisabled}
          onChange={handleAdminWeightLocationChange}
          checked={isAdminWeightLocationChecked}
        />
      </div>
      {isAdminWeightLocationChecked && (
        <MaskedTextField
          data-testid="weightRestrictionInput"
          id="weightRestrictionId"
          defaultValue="0"
          name="weightRestriction"
          label="Weight Restriction (lbs)"
          mask={Number}
          scale={0} // digits after point, 0 for integers
          signed={false} // disallow negative
          thousandsSeparator=","
          lazy={false} // immediate masking evaluation
          isDisabled={formIsDisabled}
        />
      )}
      <div className={styles.wrappedCheckbox}>
        <CheckboxField
          id="dependentsAuthorizedInput"
          data-testid="dependentsAuthorizedInput"
          name="dependentsAuthorized"
          label="Dependents authorized"
          isDisabled={formIsDisabled}
        />
      </div>
    </div>
  );
};

AllowancesDetailForm.propTypes = {
  entitlements: EntitlementShape.isRequired,
  branchOptions: DropdownArrayOf.isRequired,
  header: PropTypes.string,
  formIsDisabled: PropTypes.bool,
};

AllowancesDetailForm.defaultProps = {
  header: null,
  formIsDisabled: false,
};

export default AllowancesDetailForm;
