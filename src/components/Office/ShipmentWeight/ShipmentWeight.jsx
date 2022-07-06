import React from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import SectionWrapper from 'components/Customer/SectionWrapper';

const ShipmentWeight = () => {
  const [proGearInput, , hasProGearHelper] = useField('hasProGear');
  const [, , estimatedWeightHelper] = useField('estimatedWeight');

  const hasProGear = proGearInput.value === true;

  const handleProGear = (event) => {
    hasProGearHelper.setValue(event.target.value === 'yes');
  };
  const handleEstimatedWeight = (event) => {
    estimatedWeightHelper.setValue(event.target.value);
  };

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Weight</h2>

        <Grid row gap>
          <Grid col={6}>
            <MaskedTextField
              name="estimatedWeight"
              label="Estimated PPM weight"
              id="estimatedWeight"
              mask={Number}
              scale={0} // digits after point, 0 for integers
              signed={false} // disallow negative
              thousandsSeparator=","
              lazy={false} // immediate masking evaluation
              suffix="lbs"
              onChange={handleEstimatedWeight}
            />
            <Label className={styles.Label}>Pro-gear?</Label>
            <FormGroup>
              <Radio
                id="hasProGearYes"
                label="Yes"
                name="hasProGear"
                value="yes"
                title="Yes"
                checked={hasProGear}
                onChange={handleProGear}
              />
              <Radio
                id="proGearNo"
                label="No"
                name="proGear"
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
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default ShipmentWeight;
