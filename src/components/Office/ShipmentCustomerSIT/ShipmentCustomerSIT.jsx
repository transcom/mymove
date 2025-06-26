import React from 'react';
import PropTypes from 'prop-types';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import { DatePickerInput } from 'components/form/fields';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const ShipmentCustomerSIT = ({ sitEstimatedWeight, sitEstimatedEntryDate, sitEstimatedDepartureDate }) => {
  const [sitExpectedInput, , sitExpectedHelper] = useField('sitExpected');
  const sitExpected = sitExpectedInput.value === true;
  const [, , sitEstimatedWeightHelper] = useField('sitEstimatedWeight');
  const [, , sitEstimatedEntryDateHelper] = useField('sitEstimatedEntryDate');
  const [, , sitEstimatedDepartureDateHelper] = useField('sitEstimatedDepartureDate');

  const handleSITEstimatedWeight = (event) => {
    sitEstimatedWeightHelper.setValue(event.target.value);
    sitEstimatedWeightHelper.setTouched(true);
  };

  const handleSITExpected = (event) => {
    sitExpectedHelper.setValue(event.target.value === 'yes');

    // Handle yes/no select with respect to validation schema.
    if (event.target.value === 'no') {
      // Timeout callback handler to overcome racing condition
      // between validator and onchange event. Doing this ensures
      // schema validation behaves correctly when NO is selected.
      // If not done, schema validation would still maintain schema state when
      // YES is selected. For example, if schema validation fails
      // in YES state, we want to reset form validation back to NO state.
      setTimeout(() => {
        if (!(sitEstimatedWeight === undefined || sitEstimatedWeight === null)) {
          // restore to persisted/defaulted value
          sitEstimatedWeightHelper.setValue(sitEstimatedWeight.toString());
        } else {
          // reset input to default empty value if something was typed in.
          // this will clear out the control.
          sitEstimatedWeightHelper.setValue('');
        }
        // restore to persisted/default values
        sitEstimatedEntryDateHelper.setValue(sitEstimatedEntryDate);
        sitEstimatedDepartureDateHelper.setValue(sitEstimatedDepartureDate);
      }, 1);
    } else {
      // Timeout callback handler to overcome racing condition
      // between validator and onchange event for YES state.
      setTimeout(() => {
        // Set touched to force required message to display if default values
        // are null ensuring schema validition for YES state.
        // This is for consistently purposes with NO/YES toggling.
        sitEstimatedWeightHelper.setTouched(true);
        sitEstimatedEntryDateHelper.setTouched(true);
        sitEstimatedDepartureDateHelper.setTouched(true);
      }, 1);
    }
  };

  const [sitLocationInput, , sitLocationHelper] = useField('sitLocation');
  const sitLocationValue = sitLocationInput.value || 'DESTINATION';

  const handleSITLocation = (event) => {
    sitLocationHelper.setValue(event.target.value);
  };

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h3 className={styles.SectionHeader}>Storage in transit (SIT)</h3>

        <Grid row gap>
          <Grid col={12}>
            <FormGroup>
              <Label className={styles.Label} htmlFor="sitExpected">
                Does the customer plan to use SIT?
              </Label>
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

                {requiredAsteriskMessage}
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
                  onChange={handleSITEstimatedWeight}
                  showRequiredAsterisk
                  required
                />

                <DatePickerInput
                  name="sitEstimatedEntryDate"
                  label="Estimated storage start"
                  showRequiredAsterisk
                  required
                />

                <DatePickerInput
                  name="sitEstimatedDepartureDate"
                  label="Estimated storage end"
                  showRequiredAsterisk
                  required
                />
              </>
            )}
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

ShipmentCustomerSIT.propTypes = {
  sitEstimatedWeight: PropTypes.number.isRequired,
  sitEstimatedEntryDate: PropTypes.string.isRequired,
  sitEstimatedDepartureDate: PropTypes.string.isRequired,
};

export default ShipmentCustomerSIT;
