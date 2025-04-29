import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { useFormikContext } from 'formik';
import { Label } from '@trussworks/react-uswds';

import { isBooleanFlagEnabled } from '../../../utils/featureFlags';
import { FEATURE_FLAG_KEYS } from '../../../shared/constants';

import styles from './AllowancesDetailForm.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownInput, CheckboxField } from 'components/form/fields';
import { DropdownArrayOf } from 'types/form';
import { EntitlementShape } from 'types/order';
import { formatWeight } from 'utils/formatters';
import Hint from 'components/Hint';
import ToolTip from 'shared/ToolTip/ToolTip';

const AllowancesDetailForm = ({ header, entitlements, branchOptions, formIsDisabled, civilianTDYUBMove }) => {
  const [enableUB, setEnableUB] = useState(false);
  const renderOconusFields = !!(
    entitlements?.accompaniedTour ||
    entitlements?.dependentsTwelveAndOver ||
    entitlements?.dependentsUnderTwelve
  );
  const { values, setFieldValue } = useFormikContext();
  const [isAdminWeightLocationChecked, setIsAdminWeightLocationChecked] = useState(entitlements?.weightRestriction > 0);
  const [isAdminUBWeightLocationChecked, setIsAdminUBWeightLocationChecked] = useState(
    entitlements?.ubWeightRestriction > 0,
  );
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
      setFieldValue('weightRestriction', `${values.weightRestriction}`);
    }
  }, [setFieldValue, values.weightRestriction, isAdminWeightLocationChecked]);

  useEffect(() => {
    if (!isAdminUBWeightLocationChecked) {
      setFieldValue('ubWeightRestriction', `${values.ubWeightRestriction}`);
    }
  }, [setFieldValue, values.ubWeightRestriction, isAdminUBWeightLocationChecked]);

  const handleAdminWeightLocationChange = (e) => {
    const isChecked = e.target.checked;
    setIsAdminWeightLocationChecked(isChecked);

    if (!isChecked) {
      setFieldValue('weightRestriction', `${values.weightRestriction}`);
    } else if (isChecked && values.weightRestriction) {
      setFieldValue('weightRestriction', `${values.weightRestriction}`);
    } else {
      setFieldValue('weightRestriction', null);
    }
  };

  const handleAdminUBWeightLocationChange = (e) => {
    const isChecked = e.target.checked;
    setIsAdminUBWeightLocationChecked(isChecked);

    if (!isChecked) {
      setFieldValue('ubWeightRestriction', `${values.ubWeightRestriction}`);
    } else if (isChecked && values.ubWeightRestriction) {
      setFieldValue('ubWeightRestriction', `${values.ubWeightRestriction}`);
    } else {
      setFieldValue('ubWeightRestriction', null);
    }
  };

  useEffect(() => {
    if (civilianTDYUBMove) {
      setFieldValue('ubAllowance', `${entitlements.unaccompaniedBaggageAllowance}`);
    }
  }, [setFieldValue, entitlements.unaccompaniedBaggageAllowance, civilianTDYUBMove]);

  // Conditionally set the civilian TDY UB allowance warning message based on provided weight being in the 351 to 2000 lb range
  const showcivilianTDYUBAllowanceWarning = values.ubAllowance > 350 && values.ubAllowance <= 2000;

  let civilianTDYUBAllowanceWarning = '';
  if (showcivilianTDYUBAllowanceWarning) {
    civilianTDYUBAllowanceWarning = (
      <div className={styles.civilianUBAllowanceWarning}>
        350 lbs. is the maximum UB weight allowance for a civilian TDY move unless stated otherwise on the orders.
      </div>
    );
  }

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
      {enableUB && civilianTDYUBMove && (
        <MaskedTextField
          data-testid="civilianTdyUbAllowance"
          warning={civilianTDYUBAllowanceWarning}
          defaultValue="0"
          name="ubAllowance"
          id="civilianTdyUbAllowance"
          mask={Number}
          scale={0}
          signed={false}
          thousandsSeparator=","
          lazy={false}
          isDisabled={formIsDisabled}
          label={
            <Label className={styles.labelwithToolTip}>
              If the customer&apos;s orders specify a UB weight allowance, enter it here.
              <ToolTip
                text={
                  <span className={styles.toolTipText}>
                    Optional. If you do not specify a UB weight allowance, the default of 0 lbs will be used.
                  </span>
                }
                position="left"
                icon="info-circle"
                color="blue"
                data-testid="civilianTDYUBAllowanceToolTip"
                closeOnLeave
              />
            </Label>
          }
        />
      )}
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
          name="weightRestriction"
          label="Weight Restriction (lbs)"
          mask={Number}
          scale={0}
          signed={false}
          thousandsSeparator=","
          lazy={false}
          isDisabled={formIsDisabled}
        />
      )}
      <div className={styles.wrappedCheckbox}>
        <CheckboxField
          data-testid="adminUBWeightLocation"
          id="adminUBWeightLocation"
          name="adminRestrictedUBWeightLocation"
          label="Admin restricted UB weight location"
          isDisabled={formIsDisabled}
          onChange={handleAdminUBWeightLocationChange}
          checked={isAdminUBWeightLocationChecked}
        />
      </div>
      {isAdminUBWeightLocationChecked && (
        <MaskedTextField
          data-testid="ubWeightRestrictionInput"
          id="ubWeightRestrictionId"
          name="ubWeightRestriction"
          label="UB Weight Restriction (lbs)"
          mask={Number}
          scale={0}
          signed={false}
          thousandsSeparator=","
          lazy={false}
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
  civilianTDYUBMove: PropTypes.bool,
};

AllowancesDetailForm.defaultProps = {
  header: null,
  formIsDisabled: false,
  civilianTDYUBMove: false,
};

export default AllowancesDetailForm;
