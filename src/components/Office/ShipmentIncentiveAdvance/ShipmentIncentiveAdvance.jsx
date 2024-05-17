/* eslint-disable react/jsx-props-no-spreading */
import React, { useCallback, useEffect, useState } from 'react';
import { Fieldset, FormGroup, Radio, Grid, Label } from '@trussworks/react-uswds';
import { useField } from 'formik';
import PropTypes from 'prop-types';

import formStyles from 'styles/form.module.scss';
import styles from 'components/Office/ShipmentForm/ShipmentForm.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { calculateMaxAdvanceAndFormatAdvanceAndIncentive } from 'utils/incentives';
import { ADVANCE_STATUSES } from 'constants/ppms';

// first classing the field, and giving it an id
const useAdvanceAmountField = () => {
  const [componentProps, ...remaining] = useField({
    name: 'advance',
  });
  return [{ id: 'advanceAmountRequested', type: 'text', label: 'Amount requested', ...componentProps }, ...remaining];
};

const ShipmentIncentiveAdvance = ({ estimatedIncentive }) => {
  const [advanceAmountProps, { initialValue: initialAdvanceAmount }, amountHelper] = useAdvanceAmountField();
  const [, { initialValue: initialHasRequestedAdvance }, advanceHelper] = useField('advanceRequested');
  const [statusInput, , statusHelper] = useField('advanceStatus');
  const [, , remarksHelper] = useField('customerRemarks');
  const [advanceRequested, setDidRequestAnAdvance] = useState(initialHasRequestedAdvance);
  const setAdvanceValueCallback = useCallback(
    (advanceAmountValue) => {
      // console.log('theError', theError, advanceRequested);
      advanceHelper.setValue(advanceRequested);
      if (!advanceRequested) {
        amountHelper.setValue(advanceAmountValue);
        setTimeout(() => remarksHelper.setTouched(false, true));
      }
    },
    [advanceRequested, advanceHelper, amountHelper, remarksHelper],
  );

  // useEffect has advanceRequested as a dependency, so if the user clicks this button, it triggers a re-render and batches
  useEffect(() => {
    setAdvanceValueCallback(initialAdvanceAmount);
  }, [setAdvanceValueCallback, initialAdvanceAmount]);

  const { formattedMaxAdvance, formattedIncentive } =
    calculateMaxAdvanceAndFormatAdvanceAndIncentive(estimatedIncentive);

  const handleHasRequestedAdvanceChange = (event) => setDidRequestAnAdvance(() => event.target?.value === 'Yes');
  const handleAdvanceRequestStatusChange = (event) => statusHelper.setValue(event.target.value);
  const statusCheckedYes = statusInput.value !== ADVANCE_STATUSES.REJECTED.apiValue;
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
                    {...advanceAmountProps}
                    mask={Number}
                    scale={0} // digits after point, 0 for integers
                    signed={false} // disallow negative
                    thousandsSeparator=","
                    lazy={false} // immediate masking evaluation
                    prefix="$"
                  />
                </FormGroup>

                <FormGroup>
                  <div className={styles.AdvanceText}>Maximum advance: ${formattedMaxAdvance}</div>
                </FormGroup>
                {
                  <FormGroup>
                    <h3 className={styles.NoSpacing}>Review the advance (AOA) request:</h3>
                    <Label className={styles.Label}>Advance request status:</Label>
                    <Radio
                      id="approveAdvanceRequest"
                      label="Approve"
                      name="advanceStatus"
                      value={ADVANCE_STATUSES.APPROVED.apiValue}
                      title="Approve"
                      checked={statusCheckedYes} // defaults to false if advanceStatus has a null value
                      onChange={handleAdvanceRequestStatusChange}
                    />
                    <Radio
                      id="rejectAdvanceRequest"
                      label="Reject"
                      name="advanceStatus"
                      value={ADVANCE_STATUSES.REJECTED.apiValue}
                      title="Reject"
                      checked={!statusCheckedYes} // defaults to false if advanceStatus has a null value
                      onChange={handleAdvanceRequestStatusChange}
                    />
                  </FormGroup>
                }
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
};

ShipmentIncentiveAdvance.defaultProps = {
  estimatedIncentive: 0,
};
