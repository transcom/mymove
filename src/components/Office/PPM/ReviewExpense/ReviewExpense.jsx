import React from 'react';
import { number } from 'prop-types';
import { Formik } from 'formik';
import classnames from 'classnames';
import { Form, FormGroup, Label, Radio, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import moment from 'moment';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewExpense.module.scss';

import { formatCents, formatDate } from 'utils/formatters';
import { PPMShipmentShape, ExpenseShape } from 'types/shipment';
import Fieldset from 'shared/Fieldset';
import { DatePickerInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import ppmDocumentStatus from 'constants/ppms';
import { expenseTypeLabels, expenseTypes } from 'constants/ppmExpenseTypes';

const validationSchema = Yup.object().shape({
  amount: Yup.number().required('Enter the expense amount').min(1, 'Enter an expense amount greater than $0.00'),
  sitStartDate: Yup.date().when('movingExpenseType', {
    is: expenseTypes.STORAGE,
    then: (schema) =>
      schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
  }),
  sitEndDate: Yup.date().when('movingExpenseType', {
    is: expenseTypes.STORAGE,
    then: (schema) =>
      schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
  }),
  reason: Yup.string()
    .when('status', {
      is: ppmDocumentStatus.REJECTED,
      then: (schema) => schema.required('Add a reason why this receipt is rejected'),
    })
    .when('status', {
      is: ppmDocumentStatus.EXCLUDED,
      then: (schema) => schema.required('Add a reason why this receipt is excluded'),
    }),
  status: Yup.string().required('Reviewing this receipt is required'),
});

export default function ReviewExpense({ ppmShipment, expense, expenseNumber, ppmNumber }) {
  const { movingExpenseType, description, amount, paidWithGtcc, sitStartDate, sitEndDate, status, reason } =
    expense || {};

  const initialValues = {
    movingExpenseType: movingExpenseType || '',
    description: description || '',
    amount: amount ? `${formatCents(amount)}` : '',
    paidWithGtcc: paidWithGtcc ? 'true' : 'false',
    sitStartDate: sitStartDate ? formatDate(sitStartDate, 'YYYY-MM-DD', 'DD MMM YYYY') : '',
    sitEndDate: sitEndDate ? formatDate(sitEndDate, 'YYYY-MM-DD', 'DD MMM YYYY') : '',
    status: status || '',
    reason: reason || '',
  };
  const expenseName = movingExpenseType === expenseTypes.STORAGE ? 'Storage' : 'Receipt';
  return (
    <div className={classnames(styles.container, 'container--accent--ppm')}>
      <Formik initialValues={initialValues} validationSchema={validationSchema}>
        {({ handleChange, errors, values }) => {
          const daysInSIT =
            values.sitStartDate && values.sitEndDate && !errors.sitStartDate && !errors.sitEndDate
              ? moment(values.sitEndDate, 'DD MMM YYYY').diff(moment(values.sitStartDate, 'DD MMM YYYY'), 'days')
              : '##';
          return (
            <Form className={classnames(formStyles.form, styles.ReviewExpense)}>
              <PPMHeaderSummary ppmShipment={ppmShipment} ppmNumber={ppmNumber} />
              <hr />
              <h3 className={styles.expenseNumber}>
                {expenseName} {expenseNumber}
              </h3>
              <legend className={classnames('usa-label', styles.label)}>Expense type</legend>
              <div className={styles.displayValue}>{expenseTypeLabels[movingExpenseType]}</div>
              <legend className={classnames('usa-label', styles.label)}>Description</legend>
              <div className={styles.displayValue}>{description}</div>
              <MaskedTextField
                defaultValue="0"
                name="amount"
                label="Amount"
                id="amount"
                mask={Number}
                scale={2} // digits after point, 0 for integers
                radix="." // fractional delimiter
                mapToRadix={['.']} // symbols to process as radix
                padFractionalZeros // if true, then pads zeros at end to the length of scale
                signed={false} // disallow negative
                thousandsSeparator=","
                lazy={false} // immediate masking evaluation
                prefix="$"
              />
              {movingExpenseType === expenseTypes.STORAGE && (
                <>
                  <DatePickerInput name="sitStartDate" label="Start date" />
                  <DatePickerInput name="sitEndDate" label="End date" />
                  <legend className={classnames('usa-label', styles.label)}>Total days in SIT</legend>
                  <div className={styles.displayValue} data-testid="days-in-sit">
                    {daysInSIT}
                  </div>
                </>
              )}
              <h3 className={styles.reviewHeader}>
                Review {expenseName.toLowerCase()} {expenseNumber}
              </h3>
              <p>Add a review for this {expenseName.toLowerCase()}</p>
              <Fieldset>
                <div
                  className={classnames(approveRejectStyles.statusOption, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.APPROVED,
                  })}
                >
                  <Radio
                    id={`accept-${expense?.id}`}
                    checked={values.status === ppmDocumentStatus.APPROVED}
                    value={ppmDocumentStatus.APPROVED}
                    name="status"
                    label="Accept"
                    onChange={handleChange}
                    data-testid="acceptRadio"
                  />
                </div>
                <div
                  className={classnames(approveRejectStyles.statusOption, styles.exclude, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.EXCLUDED,
                  })}
                >
                  <Radio
                    id={`exclude-${expense?.id}`}
                    checked={values.status === ppmDocumentStatus.EXCLUDED}
                    value={ppmDocumentStatus.EXCLUDED}
                    name="status"
                    label="Exclude"
                    onChange={handleChange}
                    data-testid="excludeRadio"
                  />

                  {values.status === ppmDocumentStatus.EXCLUDED && (
                    <FormGroup className={styles.reason}>
                      <Label htmlFor={`excludeReason-${expense?.id}`}>Reason</Label>
                      <Textarea
                        id={`excludeReason-${expense?.id}`}
                        name="reason"
                        onChange={handleChange}
                        value={values.reason}
                        placeholder="Type something"
                      />
                    </FormGroup>
                  )}
                </div>
                <div
                  className={classnames(approveRejectStyles.statusOption, styles.reject, {
                    [approveRejectStyles.selected]: values.status === ppmDocumentStatus.REJECTED,
                  })}
                >
                  <Radio
                    id={`reject-${expense?.id}`}
                    checked={values.status === ppmDocumentStatus.REJECTED}
                    value={ppmDocumentStatus.REJECTED}
                    name="status"
                    label="Reject"
                    onChange={handleChange}
                    data-testid="rejectRadio"
                  />

                  {values.status === ppmDocumentStatus.REJECTED && (
                    <FormGroup className={styles.reason}>
                      <Label htmlFor={`rejectReason-${expense?.id}`}>Reason</Label>
                      <Textarea
                        id={`rejectReason-${expense?.id}`}
                        name="reason"
                        onChange={handleChange}
                        value={values.reason}
                        placeholder="Type something"
                      />
                      <div className={styles.hint}>500 characters</div>
                    </FormGroup>
                  )}
                </div>
              </Fieldset>
            </Form>
          );
        }}
      </Formik>
    </div>
  );
}

ReviewExpense.propTypes = {
  expense: ExpenseShape,
  ppmShipment: PPMShipmentShape,
  expenseNumber: number.isRequired,
  ppmNumber: number.isRequired,
};

ReviewExpense.defaultProps = {
  expense: undefined,
  ppmShipment: undefined,
};
