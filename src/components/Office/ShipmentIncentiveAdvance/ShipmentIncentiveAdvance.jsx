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

const useAdvanceRequestedField = () => {
  const [componentProps, ...remaining] = useField({
    name: 'advanceRequested',
  });
  return [
    {
      requested: { id: 'hasRequestedAdvanceYes', label: 'Yes', value: 'Yes', title: 'Yes' },
      notRequested: { id: 'hasRequestedAdvanceNo', label: 'No', value: 'No', title: 'No' },
      ...componentProps,
    },
    ...remaining,
  ];
};

const useAdvanceStatusField = () => {
  const [componentProps, ...remaining] = useField({
    name: 'advanceStatus',
  });
  return [
    {
      approved: {
        id: 'approveAdvanceRequest',
        type: 'checkbox',
        label: 'Approve',
        title: 'Approve',
        value: ADVANCE_STATUSES.APPROVED.apiValue,
      },
      rejected: {
        id: 'rejectAdvanceRequest',
        type: 'checkbox',
        label: 'Reject',
        title: 'Reject',
        value: ADVANCE_STATUSES.REJECTED.apiValue,
      },
      ...componentProps,
    },
    ...remaining,
  ];
};

const ShipmentIncentiveAdvance = ({ estimatedIncentive }) => {
  const [advanceAmountProps, { initialValue: initialAdvanceAmount }, amountHelper] = useAdvanceAmountField();

  const [
    { requested: yesAdvanceRequestedProps, notRequested: noAdvanceRequestedProps },
    { initialValue: initialHasRequestedAdvance },
    advanceHelper,
  ] = useAdvanceRequestedField();

  const [{ approved: approvedStatusProps, rejected: rejectedStatusProps, ...statusInput }, , statusHelper] =
    useAdvanceStatusField();

  const [, , remarksHelper] = useField('customerRemarks');
  const [advanceRequested, setDidRequestAnAdvance] = useState(initialHasRequestedAdvance);

  const setAdvanceValueCallback = useCallback(
    (advanceAmountValue) => {
      advanceHelper.setValue(advanceRequested);
      if (!advanceRequested) {
        amountHelper.setValue(advanceAmountValue);
        remarksHelper.setTouched(false, true);
      }
    },
    [advanceRequested, advanceHelper, amountHelper, remarksHelper],
  );

  useEffect(() => {
    setTimeout(() => setAdvanceValueCallback(initialAdvanceAmount));
  }, [setAdvanceValueCallback, initialAdvanceAmount, advanceRequested]);

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
                {...yesAdvanceRequestedProps}
                checked={advanceRequested}
                onChange={handleHasRequestedAdvanceChange}
              />
              <Radio
                {...noAdvanceRequestedProps}
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
                <FormGroup>
                  <h3 className={styles.NoSpacing}>Review the advance (AOA) request:</h3>
                  <Label className={styles.Label}>Advance request status:</Label>
                  <Radio
                    {...approvedStatusProps}
                    checked={statusCheckedYes} // defaults to false if advanceStatus has a null value
                    onChange={handleAdvanceRequestStatusChange}
                  />
                  <Radio
                    {...rejectedStatusProps}
                    checked={!statusCheckedYes} // defaults to false if advanceStatus has a null value
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
};

ShipmentIncentiveAdvance.defaultProps = {
  estimatedIncentive: 0,
};
