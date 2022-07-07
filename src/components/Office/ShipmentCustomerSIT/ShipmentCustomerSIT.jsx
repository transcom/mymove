import React from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label, Button } from '@trussworks/react-uswds';
import { useField } from 'formik';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { DatePickerInput } from 'components/form/fields';

const ShipmentCustomerSIT = ({ onCalculateClick }) => {
  const [sitExpectedInput, , sitExpectedHelper] = useField('sitExpected');
  const sitExpected = sitExpectedInput.value === true;

  const handleSITExpected = (event) => {
    sitExpectedHelper.setValue(event.target.value === 'yes');
  };

  const [sitLocationInput, , sitLocationHelper] = useField('sitLocation');
  const sitLocationValue = sitLocationInput.value || 'DESTINATION';

  const [weightInput] = useField('sitEstimatedWeight');
  const weightValue = weightInput.value;

  const [startInput] = useField('sitEstimatedEntryDate');
  const startValue = startInput.value;

  const [endInput] = useField('sitEstimatedDepartureDate');
  const endValue = endInput.value;

  const handleSITLocation = (event) => {
    sitLocationHelper.setValue(event.target.value);
  };

  const isCalculationEnabled =
    !!startValue && startValue !== 'Invalid date' && !!endValue && endValue !== 'Invalid date' && !!weightValue;

  const handleCalculateClick = async () => {
    /*
       The SIT calcuation work is being left for a later sprint.
       We'll want to do _something_ with the results of this function.
     */
    onCalculateClick({ location: sitLocationValue, start: startValue, end: endValue, weight: weightValue });
  };

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Storage in transit (SIT)</h2>

        <Grid row gap>
          <Grid col={12}>
            <FormGroup>
              <Label className={styles.Label}>Does the customer plan to use SIT?</Label>
              <Radio
                id="sitExpectedYes"
                label="Yes"
                name="sitExpected"
                value="yes"
                title="Yes"
                checked={sitExpected}
                onChange={handleSITExpected}
              />
              <Radio
                id="sitExpectedNo"
                label="No"
                name="sitExpected"
                value="no"
                title="No"
                checked={!sitExpected}
                onChange={handleSITExpected}
              />
            </FormGroup>

            {sitExpected && (
              <>
                <FormGroup>
                  <Label className={styles.Label}>Origin or destination?</Label>
                  <Radio
                    id="sitLocationOrigin"
                    label="Origin"
                    name="sitLocation"
                    value="ORIGIN"
                    title="Origin"
                    checked={sitLocationValue === 'ORIGIN'}
                    onChange={handleSITLocation}
                  />
                  <Radio
                    id="sitLocationDestination"
                    label="Destination"
                    name="sitLocation"
                    value="DESTINATION"
                    title="Destination"
                    checked={sitLocationValue === 'DESTINATION'}
                    onChange={handleSITLocation}
                  />
                </FormGroup>

                <MaskedTextField
                  name="sitEstimatedWeight"
                  label="Estimated SIT weight"
                  id="sitEstimatedWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                />

                <DatePickerInput name="sitEstimatedEntryDate" label="Estimated storage start" />

                <DatePickerInput name="sitEstimatedDepartureDate" label="Estimated storage end" />

                <Button type="button" secondary disabled={!isCalculationEnabled} onClick={handleCalculateClick}>
                  Calculate SIT
                </Button>
              </>
            )}
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default ShipmentCustomerSIT;

ShipmentCustomerSIT.propTypes = {
  onCalculateClick: PropTypes.func,
};

ShipmentCustomerSIT.defaultProps = {
  onCalculateClick: () => {},
};
