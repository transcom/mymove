import React, { useEffect, useState } from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import Hint from 'components/Hint';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const ShipmentWeight = ({ onEstimatedWeightChange }) => {
  const [proGearInput, , hasProGearHelper] = useField('hasProGear');
  const [gunSafeInput, , hasGunSafeHelper] = useField('hasGunSafe');
  const [, , estimatedWeightHelper] = useField('estimatedWeight');
  const [isGunSafeEnabled, setIsGunSafeEnabled] = useState(false);

  const handleEstimatedWeightChange = (value) => {
    onEstimatedWeightChange(value);
  };

  const handleEstimatedWeight = (event) => {
    const value = event.target.value.replace(/,/g, ''); // removing comma to avoid NaN issue.
    estimatedWeightHelper.setValue(value);
    estimatedWeightHelper.setTouched(true);
    handleEstimatedWeightChange(value);
  };

  const hasProGear = proGearInput.value === true;
  const hasGunSafe = gunSafeInput.value === true;

  const handleProGear = (event) => {
    hasProGearHelper.setValue(event.target.value === 'yes');
  };
  const handleGunSafe = (event) => {
    hasGunSafeHelper.setValue(event.target.value === 'yes');
  };

  useEffect(() => {
    const fetchData = async () => {
      setIsGunSafeEnabled(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.GUN_SAFE));
    };
    fetchData();
  }, []);

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h3 className={styles.SectionHeader}>Weight</h3>
        {requiredAsteriskMessage}

        <Grid row gap>
          <Grid>
            <MaskedTextField
              name="estimatedWeight"
              label="Estimated PPM weight"
              data-testid="estimatedWeight"
              id="estimatedWeight"
              showRequiredAsterisk
              required
              mask={Number}
              scale={0} // digits after point, 0 for integers
              signed={false} // disallow negative
              thousandsSeparator=","
              lazy={false} // immediate masking evaluation
              suffix="lbs"
              onInput={handleEstimatedWeight}
            />
            <Label className={styles.radioLabel}>
              <span>Pro-gear?</span>
            </Label>
            <FormGroup className={styles.radioGroup}>
              <Radio
                id="hasProGearYes"
                label="Yes"
                name="hasProGear"
                data-testid="hasProGearYes"
                value="yes"
                title="Yes"
                checked={hasProGear}
                onChange={handleProGear}
              />
              <Radio
                id="proGearNo"
                label="No"
                name="proGear"
                data-testid="hasProGearNo"
                value="no"
                title="No"
                checked={!hasProGear}
                onChange={handleProGear}
              />
            </FormGroup>
            {hasProGear && (
              <>
                <MaskedTextField
                  defaultValue="0"
                  name="proGearWeight"
                  label="Estimated pro-gear weight"
                  id="proGearWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                />
                <MaskedTextField
                  defaultValue="0"
                  name="spouseProGearWeight"
                  label="Estimated spouse pro-gear weight"
                  id="spouseProGearWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                />
              </>
            )}
            {isGunSafeEnabled && (
              <>
                <Label className={styles.radioLabel}>Gun safe?</Label>
                <FormGroup className={styles.radioGroup}>
                  <Radio
                    id="hasGunSafeYes"
                    label="Yes"
                    name="hasGunSafe"
                    data-testid="hasGunSafeYes"
                    value="yes"
                    title="Yes"
                    checked={hasGunSafe}
                    onChange={handleGunSafe}
                  />
                  <Radio
                    id="hasGunSafeNo"
                    label="No"
                    name="hasGunSafeNo"
                    data-testid="hasGunSafeNo"
                    value="no"
                    title="No"
                    checked={!hasGunSafe}
                    onChange={handleGunSafe}
                  />
                </FormGroup>
                {hasGunSafe && (
                  <>
                    <MaskedTextField
                      defaultValue="0"
                      name="gunSafeWeight"
                      label="Estimated gun safe weight"
                      id="gunSafeWeight"
                      mask={Number}
                      scale={0} // digits after point, 0 for integers
                      thousandsSeparator=","
                      lazy={false} // immediate masking evaluation
                      suffix="lbs"
                      showRequiredAsterisk
                      required
                    />
                    <Hint>
                      The government authorizes the shipment of a gun safe up to 500 lbs. The weight entitlement is
                      charged for any weight over 500 lbs. The additional 500 lbs gun safe weight entitlement cannot be
                      applied if a customer&apos;s overall entitlement is already at the 18,000 lbs maximum.
                    </Hint>
                  </>
                )}
              </>
            )}
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default ShipmentWeight;
