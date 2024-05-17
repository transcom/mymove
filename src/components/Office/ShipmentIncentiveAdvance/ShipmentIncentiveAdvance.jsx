import React from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { calculateMaxAdvanceAndFormatAdvanceAndIncentive } from 'utils/incentives';
import { ADVANCE_STATUSES } from 'constants/ppms';

const ShipmentIncentiveAdvance = ({ estimatedIncentive, advanceAmountRequested }) => {
  const [advanceInput, , advanceHelper] = useField('advanceRequested');
  const [statusInput, , statusHelper] = useField('advanceStatus');
  const [, , advanceAmountRequestedProps] = useField('advance');

  const advanceRequested = String(advanceInput.value) === 'true';
  const advanceRequestStatus =
    statusInput.value === ADVANCE_STATUSES.APPROVED.apiValue || statusInput.value === ADVANCE_STATUSES.EDITED.apiValue;

  const { formattedMaxAdvance, formattedIncentive } =
    calculateMaxAdvanceAndFormatAdvanceAndIncentive(estimatedIncentive);

  const handleHasRequestedAdvanceChange = (event) => {
    const selected = event.target.value;
    advanceHelper.setValue(selected === 'Yes');
  };

  // Current '0' for scale value indicate zero decimal places.
  // Input must be an integer.
  const advanceAmountRequestedMaskedTextFieldScale = 0;

  // Denominator calculation to convert raw advanceAmountRequested value to mask value. (base ^ exponent)
  // ie..takes 10100 to 101 for display with a scale of 0.
  const advanceAmountRequestedMaskTaskFieldValueDenominator = 10 ** (2 - advanceAmountRequestedMaskedTextFieldScale);

  const handleAdvanceRequestStatusChange = (event) => {
    const selected = event.target.value;
    statusHelper.setValue(selected);
    if (selected === ADVANCE_STATUSES.REJECTED.apiValue) {
      // Wrap in timout callback due to racing condition with state changes
      // with respect to form validator. Doing this ensures setValue uses correct AdvanceStatus.
      setTimeout(() => {
        // Programmatically undo unsaved input to persisted value
        advanceAmountRequestedProps.setValue(
          `${advanceAmountRequested / advanceAmountRequestedMaskTaskFieldValueDenominator}`,
        );
      }, 100);
    }
  };

  return (
    <SectionWrapper className={formStyles.formSection}>
      <Fieldset className={styles.Fieldset}>
        <h2 className={styles.SectionHeader}>Incentive &amp; advance</h2>
        <h3 className={styles.NoSpacing}>Estimated incentive: ${formattedIncentive}</h3>

        <Grid row>
          <Grid col={12}>
            <FormGroup>
              <Label className={styles.Label}>Advance (AOA) requested?</Label>
              <Radio
                id="hasRequestedAdvanceYes"
                label="Yes"
                name="advanceRequested"
                value="Yes"
                title="Yes"
                checked={advanceRequested}
                onChange={handleHasRequestedAdvanceChange}
              />
              <Radio
                id="hasRequestedAdvanceNo"
                label="No"
                name="advanceRequested"
                value="No"
                title="No"
                checked={!advanceRequested}
                onChange={handleHasRequestedAdvanceChange}
              />
            </FormGroup>

            {advanceRequested && (
              <>
                <FormGroup>
                  <MaskedTextField
                    defaultValue="0"
                    name="advance"
                    label="Amount requested"
                    id="advanceAmountRequested"
                    mask={Number}
                    scale={advanceAmountRequestedMaskedTextFieldScale} // digits after point, 0 for integers
                    signed={false} // disallow negative
                    thousandsSeparator=","
                    lazy={false} // immediate masking evaluation
                    prefix="$"
                  />
                </FormGroup>

                <FormGroup>
                  <div className={styles.AdvanceText}>Maximum advance: ${formattedMaxAdvance}</div>
                </FormGroup>

                <FormGroup>
                  <h3 className={styles.NoSpacing}>Review the advance (AOA) request:</h3>
                  <Label className={styles.Label}>Advance request status:</Label>
                  <Radio
                    id="approveAdvanceRequest"
                    label="Approve"
                    name="advanceStatus"
                    value={ADVANCE_STATUSES.APPROVED.apiValue}
                    title="Approve"
                    checked={!!statusInput.value && advanceRequestStatus} // defaults to false if advanceStatus has a null value
                    onChange={handleAdvanceRequestStatusChange}
                  />
                  <Radio
                    id="rejectAdvanceRequest"
                    label="Reject"
                    name="advanceStatus"
                    value={ADVANCE_STATUSES.REJECTED.apiValue}
                    title="Reject"
                    checked={!!statusInput.value && !advanceRequestStatus} // defaults to false if advanceStatus has a null value
                    onChange={handleAdvanceRequestStatusChange}
                  />
                </FormGroup>
              </>
            )}
          </Grid>
        </Grid>
      </Fieldset>
    </SectionWrapper>
  );
};

export default ShipmentIncentiveAdvance;

ShipmentIncentiveAdvance.propTypes = {
  estimatedIncentive: PropTypes.number,
  advanceAmountRequested: PropTypes.number,
};

ShipmentIncentiveAdvance.defaultProps = {
  estimatedIncentive: 0,
  advanceAmountRequested: 0,
};
