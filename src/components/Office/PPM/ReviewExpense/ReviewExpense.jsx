import React, { useEffect, useCallback } from 'react';
import { useMutation } from '@tanstack/react-query';
import { func, number, object } from 'prop-types';
import { Formik, Field } from 'formik';
import classnames from 'classnames';
import { FormGroup, Label, Radio, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import moment from 'moment';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewExpense.module.scss';

import {
  formatCents,
  formatDate,
  formatWeight,
  dropdownInputOptions,
  removeCommas,
  toDollarString,
} from 'utils/formatters';
import { ExpenseShape } from 'types/shipment';
import { OrderShape } from 'types/order';
import Fieldset from 'shared/Fieldset';
import { DatePickerInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import ppmDocumentStatus from 'constants/ppms';
import { expenseTypes, ppmExpenseTypes, getExpenseTypeValue, llvmExpenseTypes } from 'constants/ppmExpenseTypes';
import { ErrorMessage, Form } from 'components/form';
import { patchExpense } from 'services/ghcApi';
import { convertDollarsToCents } from 'shared/utils';
import TextField from 'components/form/fields/TextField/TextField';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import SitCost from 'components/Office/PPM/SitCost/SitCost';
import { useGetPPMSITEstimatedCostQuery } from 'hooks/queries';

const sitLocationOptions = dropdownInputOptions(LOCATION_TYPES);

const validationSchema = (maxWeight) => {
  return Yup.object().shape({
    amount: Yup.string()
      .required('Enter the expense amount')
      .notOneOf(['0', '0.00'], 'Enter an expense amount greater than $0.00'),
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
    weightStored: Yup.number().when('movingExpenseType', {
      is: expenseTypes.STORAGE,
      then: (schema) =>
        schema
          .required('Required')
          .max(maxWeight, `Weight must be less than total PPM weight of ${formatWeight(maxWeight)}`)
          .min(1, `Enter a weight greater than 0 lbs`),
    }),
    sitLocation: Yup.mixed().when('movingExpenseType', {
      is: expenseTypes.STORAGE,
      then: (schema) => schema.oneOf(sitLocationOptions.map((i) => i.key)).required('Required'),
    }),
  });
};

export default function ReviewExpense({
  ppmShipmentInfo,
  expense,
  documentSets,
  documentSetIndex,
  categoryIndex,
  tripNumber,
  ppmNumber,
  onError,
  onSuccess,
  formRef,
  readOnly,
  order,
}) {
  const {
    movingExpenseType,
    description,
    amount,
    paidWithGtcc,
    sitStartDate,
    sitEndDate,
    status,
    reason,
    weightStored,
    sitLocation,
  } = expense || {};

  const { mutate: patchExpenseMutation } = useMutation(patchExpense, {
    onSuccess,
    onError,
  });

  const [descriptionString, setDescriptionString] = React.useState(description || '');
  const actualWeight = ppmShipmentInfo?.actualWeight || '';
  const [amountValue, setAmountValue] = React.useState(amount?.toString() || '');
  const [weightStoredValue, setWeightStoredValue] = React.useState(weightStored);
  const [ppmSITLocation, setSITLocation] = React.useState(sitLocation?.toString() || 'DESTINATION');
  const [sitStartDateValue, setSitStartDateValue] = React.useState(sitStartDate != null ? sitStartDate : '');
  const [sitEndDateValue, setSitEndDateValue] = React.useState(sitEndDate != null ? sitEndDate : '');
  const displaySitCost =
    ppmSITLocation !== '' && sitStartDateValue !== '' && sitEndDateValue !== '' && weightStoredValue !== '';
  const [estimatedCost, setEstimatedCost] = React.useState(0);
  const [actualSITReimbursed, setActualSITReimbursed] = React.useState(
    amountValue < estimatedCost ? amountValue : estimatedCost,
  );
  const initialValues = {
    movingExpenseType: movingExpenseType || '',
    description: descriptionString,
    amount: amountValue ? `${formatCents(amountValue)}` : '',
    paidWithGtcc: paidWithGtcc ? 'true' : 'false',
    sitStartDate: sitStartDateValue,
    sitEndDate: sitEndDateValue,
    status: status || '',
    reason: reason || '',
    weightStored: weightStoredValue?.toString() || '',
    sitLocation: ppmSITLocation,
  };

  const [selectedExpenseType, setSelectedExpenseType] = React.useState(getExpenseTypeValue(movingExpenseType)); // Set initial expense type via value received from backend
  const [currentCategoryIndex, setCurrentCategoryIndex] = React.useState(categoryIndex);
  const [samePage, setSamePage] = React.useState(false); // Helps track if back button was used or not
  /**
   * Gets the current index for the receipt type, i.e. if we've already reviewed two "Oil" expense receipts, and user chooses "Oil" for expense type,
   * then this will display "Oil #3" at bottom of page.
   * * */
  const computeCurrentCategoryIndex = useCallback(
    (expenseType) => {
      let count = 0;
      const expenseDocs = documentSets.filter((docSet) => docSet.documentSetType === 'MOVING_EXPENSE'); // documentSets includes Trip weight tickets, progear, etc. that we don't need
      const docsFiltered = documentSets.length - expenseDocs.length; // Reduce count/index by number of docs filtered out
      for (let i = 0; i < documentSetIndex - docsFiltered; i += 1) {
        if (expenseDocs[i].documentSet.movingExpenseType === expenseType) count += 1;
      }
      return count + 1;
    },
    [documentSetIndex, documentSets],
  );

  useEffect(() => {
    // Don't update from parent component if user just changed the dropdown field. I.e. this only fires on submit or back button
    if (!samePage) setSelectedExpenseType(getExpenseTypeValue(movingExpenseType));

    const selectedExpenseTypeKey = llvmExpenseTypes[selectedExpenseType]; // Convert nice "stringified" value back into an enum key for ppmExpenseTypes
    const index = computeCurrentCategoryIndex(selectedExpenseTypeKey); // Get index for number at bottom of page (e.x. "Contracted Expense #2")
    setCurrentCategoryIndex(index);
  }, [movingExpenseType, tripNumber, computeCurrentCategoryIndex, selectedExpenseType, samePage]);

  // If parent state updates to show that we've moved onto another document, then user must've used back or submit button
  useEffect(() => {
    setSamePage(false);
  }, [documentSetIndex]);

  useEffect(() => {
    if (displaySitCost) {
      const value = parseInt(removeCommas(amountValue), 10);
      setActualSITReimbursed(value < estimatedCost ? value : estimatedCost);
    }
  }, [estimatedCost, amountValue, displaySitCost]);

  const handleSubmit = (values) => {
    if (readOnly) {
      onSuccess();
      return;
    }

    // To prevent errors when submitting the request and for better error messages we notify the user which fields are still required.
    // This is also to done because Formik can fail to perform validation correctly when components are refreshed.
    if (selectedExpenseType.toUpperCase() === expenseTypes.STORAGE) {
      let errorMessage = '';

      if (sitStartDateValue === '') {
        errorMessage += 'SIT Start Date is required.\n';
      }

      if (sitEndDateValue === '') {
        errorMessage += 'SIT End Date is required.\n';
      }

      if (weightStoredValue === null) {
        errorMessage += 'Weight Stored is required.\n';
      }

      if (ppmSITLocation === '') {
        errorMessage += 'SIT Location is required.\n';
      }

      if (errorMessage !== '') {
        onError(errorMessage);
        return;
      }
    }

    const payload = {
      ppmShipmentId: expense.ppmShipmentId,
      movingExpenseType: llvmExpenseTypes[selectedExpenseType],
      description: values.description,
      amount: convertDollarsToCents(values.amount),
      paidWithGtcc: values.paidWithGtcc,
      sitStartDate: llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE ? sitStartDateValue : undefined,
      sitEndDate: llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE ? sitEndDateValue : undefined,
      reason: values.status === ppmDocumentStatus.APPROVED ? null : values.reason,
      status: values.status,
      weightStored: llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE ? weightStoredValue : undefined,
      sitLocation: llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE ? ppmSITLocation : undefined,
      sitReimburseableAmount:
        llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE ? actualSITReimbursed : undefined,
    };

    patchExpenseMutation({
      ppmShipmentId: expense.ppmShipmentId,
      movingExpenseId: expense.id,
      payload,
      eTag: expense.eTag,
    });
  };

  const titleCase = (input) => input.charAt(0).toUpperCase() + input.slice(1);
  const allCase = (input) => input?.split(' ').map(titleCase).join(' ') ?? '';

  return (
    <div className={classnames(styles.container, 'container--accent--ppm')}>
      <Formik
        initialValues={initialValues}
        validationSchema={() => validationSchema(actualWeight)}
        innerRef={formRef}
        onSubmit={handleSubmit}
        enableReinitialize
        validateOnMount
      >
        {({ handleChange, errors, setFieldError, setFieldTouched, setFieldValue, touched, values }) => {
          const handleApprovalChange = (event) => {
            handleChange(event);
            setFieldValue('reason', '');
            setFieldTouched('reason', false, false);
            setFieldError('reason', null);
          };

          const refreshPage = (event) => {
            setSamePage(true);
            const count = computeCurrentCategoryIndex(event.target.value);
            setCurrentCategoryIndex(count + 1);
          };

          const handleSITLocationChange = (event) => {
            setSITLocation(event.target.value);
            refreshPage(event);
          };

          const handleWeightStoredChange = (event) => {
            const weight = parseInt(removeCommas(event.target.value), 10);
            if (weight <= actualWeight && weight > 0) {
              setWeightStoredValue(weight);
              refreshPage(event);
            }
          };

          const handleSitStartDateChange = (value) => {
            const date = formatDate(value, 'DD MMM YYYY', 'YYYY-MM-DD');
            setSitStartDateValue(date);
            setSamePage(true);
            const count = computeCurrentCategoryIndex(value);
            setCurrentCategoryIndex(count + 1);
          };

          const handleSitEndDateChange = (value) => {
            const date = formatDate(value, 'DD MMM YYYY', 'YYYY-MM-DD');
            setSitEndDateValue(date);
            setSamePage(true);
            const count = computeCurrentCategoryIndex(value);
            setCurrentCategoryIndex(count + 1);
          };

          const sitAdditionalStartDate = sitStartDateValue
            ? moment(sitStartDateValue, 'YYYY-MM-DD').add(1, 'days')
            : '##';

          const daysInSIT =
            sitStartDateValue && sitEndDateValue
              ? moment(sitEndDateValue, 'YYYY-MM-DD')
                  .add(1, 'days')
                  .diff(moment(sitStartDateValue, 'YYYY-MM-DD'), 'days')
              : '##';

          return (
            <>
              <div className={classnames(formStyles.form, styles.ReviewExpense, styles.headerContainer)}>
                <PPMHeaderSummary
                  ppmShipmentInfo={ppmShipmentInfo}
                  order={order}
                  ppmNumber={ppmNumber}
                  showAllFields={false}
                  readOnly={readOnly}
                />
              </div>
              <Form className={classnames(formStyles.form, styles.ReviewExpense)}>
                <hr />
                <h3 className={styles.tripNumber}>{`Receipt ${tripNumber}`}</h3>
                <div className="labelWrapper">
                  <Label htmlFor="movingExpenseType">Expense Type</Label>
                </div>
                <select
                  label="Expense Type"
                  name="movingExpenseType"
                  id="movingExpenseType"
                  required
                  className={classnames('usa-select')}
                  value={selectedExpenseType}
                  disabled={readOnly}
                  onChange={(e) => {
                    setSelectedExpenseType(e.target.value);
                    refreshPage(e);
                  }}
                >
                  {ppmExpenseTypes.map((x) => (
                    <option key={x.key}>{x.value}</option>
                  ))}
                </select>
                <TextField
                  defaultValue={description}
                  name="description"
                  label="Description"
                  id="description"
                  className={styles.displayValue}
                  disabled={readOnly}
                  onBlur={(e) => {
                    setDescriptionString(e.target.value);
                  }}
                />
                {llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE && (
                  <>
                    <div className="labelWrapper">
                      <Label htmlFor="sitLocationInput">SIT Location</Label>
                    </div>
                    <Field
                      as={Radio}
                      id="sitLocationOrigin"
                      label="Origin"
                      name="sitLocation"
                      value="ORIGIN"
                      checked={values.sitLocation === 'ORIGIN'}
                      disabled={readOnly}
                      onChange={(e) => {
                        handleSITLocationChange(e);
                      }}
                    />
                    <Field
                      as={Radio}
                      id="sitLocationDestination"
                      label="Destination"
                      name="sitLocation"
                      value="DESTINATION"
                      checked={values.sitLocation === 'DESTINATION'}
                      disabled={readOnly}
                      onChange={(e) => {
                        handleSITLocationChange(e);
                      }}
                    />
                    {displaySitCost && (
                      <SitCost
                        ppmShipmentInfo={ppmShipmentInfo}
                        ppmSITLocation={ppmSITLocation}
                        sitStartDate={sitStartDateValue}
                        sitAdditionalStartDate={sitAdditionalStartDate}
                        sitEndDate={sitEndDateValue}
                        weightStored={weightStoredValue}
                        actualWeight={actualWeight}
                        useQueries={useGetPPMSITEstimatedCostQuery}
                        setEstimatedCost={setEstimatedCost}
                      />
                    )}
                  </>
                )}
                <MaskedTextField
                  defaultValue="0"
                  name="amount"
                  label="Amount Requested"
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
                  isDisabled={readOnly}
                  onBlur={(e) => {
                    const newAmount = e.target.value.replace(/[,.]/g, '');
                    setAmountValue(newAmount);
                  }}
                />
                {llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE && (
                  <>
                    <div>
                      <legend className={classnames('usa-label', styles.label)}>Actual SIT Reimbursement</legend>
                      <div className={styles.displayValue} data-testid="actual-sit-reimbursement">
                        {toDollarString(formatCents(actualSITReimbursed))}
                      </div>
                    </div>
                    <MaskedTextField
                      defaultValue="0"
                      name="weightStored"
                      label="Weight Stored"
                      id="weightStored"
                      mask={Number}
                      scale={0} // digits after point, 0 for integers
                      signed={false} // disallow negative
                      thousandsSeparator=","
                      lazy={false} // immediate masking evaluation
                      suffix="lbs"
                      isDisabled={readOnly}
                      onBlur={(e) => {
                        handleWeightStoredChange(e);
                      }}
                    />
                    <div>
                      <legend className={classnames('usa-label', styles.label)}>Actual PPM Weight</legend>
                      <div className={styles.displayValue}>{formatWeight(actualWeight)}</div>
                    </div>
                    <DatePickerInput
                      name="sitStartDate"
                      label="Start date"
                      required
                      disabled={readOnly}
                      onChange={(value) => {
                        handleSitStartDateChange(value);
                      }}
                    />
                    <DatePickerInput
                      name="sitEndDate"
                      label="End date"
                      required
                      disabled={readOnly}
                      onChange={(value) => {
                        handleSitEndDateChange(value);
                      }}
                    />
                    <legend className={classnames('usa-label', styles.label)}>Total days in SIT</legend>
                    <div className={styles.displayValue} data-testid="days-in-sit">
                      {daysInSIT}
                    </div>
                  </>
                )}
                <h3 className={styles.reviewHeader}>{`Review ${allCase(
                  selectedExpenseType,
                )} #${currentCategoryIndex}`}</h3>
                <p>Add a review for this {allCase(selectedExpenseType)}</p>
                <ErrorMessage display={!!errors?.status && !!touched?.status}>{errors.status}</ErrorMessage>
                <Fieldset className={styles.statusOptions}>
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
                      onChange={handleApprovalChange}
                      data-testid="acceptRadio"
                      disabled={readOnly}
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
                      disabled={readOnly}
                    />
                    {values.status === ppmDocumentStatus.EXCLUDED && (
                      <FormGroup className={styles.reason}>
                        <Label htmlFor={`excludeReason-${expense?.id}`}>Reason</Label>
                        <ErrorMessage display={!!errors?.reason && !!touched?.reason}>{errors.reason}</ErrorMessage>
                        <Textarea
                          id={`excludeReason-${expense?.id}`}
                          name="reason"
                          onChange={handleChange}
                          value={values.reason}
                          placeholder="Type something"
                          disabled={readOnly}
                        />
                        <div className={styles.hint}>{500 - values.reason.length} characters</div>
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
                      disabled={readOnly}
                    />

                    {values.status === ppmDocumentStatus.REJECTED && (
                      <FormGroup className={styles.reason}>
                        <Label htmlFor={`rejectReason-${expense?.id}`}>Reason</Label>
                        <ErrorMessage display={!!errors?.reason && !!touched?.reason}>{errors.reason}</ErrorMessage>
                        <Textarea
                          id={`rejectReason-${expense?.id}`}
                          name="reason"
                          onChange={handleChange}
                          value={values.reason}
                          placeholder="Type something"
                          disabled={readOnly}
                        />
                        <div className={styles.hint}>{500 - values.reason.length} characters</div>
                      </FormGroup>
                    )}
                  </div>
                </Fieldset>
              </Form>
            </>
          );
        }}
      </Formik>
    </div>
  );
}

ReviewExpense.propTypes = {
  expense: ExpenseShape,
  tripNumber: number.isRequired,
  ppmNumber: number.isRequired,
  onSuccess: func,
  formRef: object,
  order: OrderShape.isRequired,
};

ReviewExpense.defaultProps = {
  expense: undefined,
  onSuccess: null,
  formRef: null,
};
