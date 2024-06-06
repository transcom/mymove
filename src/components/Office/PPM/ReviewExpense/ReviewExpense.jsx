import React, { useEffect, useCallback } from 'react';
import { useMutation } from '@tanstack/react-query';
import { func, number, object } from 'prop-types';
import { Formik } from 'formik';
import classnames from 'classnames';
import { FormGroup, Label, Radio, Textarea } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import moment from 'moment';

import PPMHeaderSummary from '../PPMHeaderSummary/PPMHeaderSummary';

import styles from './ReviewExpense.module.scss';

import { formatCents, formatDate, dropdownInputOptions, toDollarString } from 'utils/formatters';
import { ExpenseShape } from 'types/shipment';
import Fieldset from 'shared/Fieldset';
import { DatePickerInput, DropdownInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import formStyles from 'styles/form.module.scss';
import approveRejectStyles from 'styles/approveRejectControls.module.scss';
import ppmDocumentStatus from 'constants/ppms';
import { expenseTypes, ppmExpenseTypes, getExpenseTypeValue, llvmExpenseTypes } from 'constants/ppmExpenseTypes';
import { ErrorMessage, Form } from 'components/form';
import { patchExpense, patchPPMSIT } from 'services/ghcApi';
import { convertDollarsToCents } from 'shared/utils';
import TextField from 'components/form/fields/TextField/TextField';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import { useGetPPMSITEstimatedCostQuery } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const sitLocationOptions = dropdownInputOptions(LOCATION_TYPES);

const validationSchema = (allowableWeight) => {
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
        schema.required('Required').max(allowableWeight, `Enter an amount ${allowableWeight} lbs or less`),
    }),
    sitLocation: Yup.mixed()
      .oneOf(sitLocationOptions.map((i) => i.key))
      .required('Required'),
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

  const { mutate: patchPPMSITMutation } = useMutation(patchPPMSIT);

  const allowableWeight = ppmShipmentInfo.estimatedWeight;
  const [ppmSITLocation, setSITLocation] = React.useState(sitLocation?.toString() || '');
  const { estimatedCost, isLoading, isError } = useGetPPMSITEstimatedCostQuery(ppmShipmentInfo.id, ppmSITLocation);

  const initialValues = {
    movingExpenseType: movingExpenseType || '',
    description: description || '',
    amount: amount ? `${formatCents(amount)}` : '',
    paidWithGtcc: paidWithGtcc ? 'true' : 'false',
    sitStartDate: sitStartDate ? formatDate(sitStartDate, 'YYYY-MM-DD', 'DD MMM YYYY') : '',
    sitEndDate: sitEndDate ? formatDate(sitEndDate, 'YYYY-MM-DD', 'DD MMM YYYY') : '',
    weightStored: weightStored?.toString() || '',
    status: status || '',
    reason: reason || '',
    actualWeight: ppmShipmentInfo?.actualWeight?.toString() || '',
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

  const handleSubmit = (values) => {
    if (values.movingExpenseType === 'STORAGE') {
      const ppmSitPayload = {
        sitLocation: ppmSITLocation,
      };
      patchPPMSITMutation({
        ppmShipmentId: expense.ppmShipmentId,
        payload: ppmSitPayload,
        eTag: ppmShipmentInfo.eTag,
      });
    }
    const payload = {
      ppmShipmentId: expense.ppmShipmentId,
      movingExpenseType: llvmExpenseTypes[selectedExpenseType],
      description: values.description,
      amount: convertDollarsToCents(values.amount),
      paidWithGtcc: values.paidWithGtcc,
      sitStartDate: formatDate(values.sitStartDate, 'DD MMM YYYY', 'YYYY-MM-DD'),
      sitEndDate: formatDate(values.sitEndDate, 'DD MMM YYYY', 'YYYY-MM-DD'),
      reason: values.status === ppmDocumentStatus.APPROVED ? null : values.reason,
      status: values.status,
      weightStored: Number.parseInt(values.weightStored, 10),
      sitLocation: ppmSITLocation,
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
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;
  return (
    <div className={classnames(styles.container, 'container--accent--ppm')}>
      <Formik
        initialValues={initialValues}
        validationSchema={() => validationSchema(allowableWeight)}
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

          const handleSITLocationChange = (event) => {
            setSITLocation(event.target.value);
            setSamePage(true);
            const count = computeCurrentCategoryIndex(event.target.value);
            setCurrentCategoryIndex(count + 1);
          };

          const daysInSIT =
            values.sitStartDate && values.sitEndDate && !errors.sitStartDate && !errors.sitEndDate
              ? moment(values.sitEndDate, 'DD MMM YYYY')
                  .add(1, 'days')
                  .diff(moment(values.sitStartDate, 'DD MMM YYYY'), 'days')
              : '##';

          return (
            <Form className={classnames(formStyles.form, styles.ReviewExpense)}>
              <PPMHeaderSummary ppmShipmentInfo={ppmShipmentInfo} ppmNumber={ppmNumber} showAllFields={false} />
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
                onChange={(e) => {
                  setSelectedExpenseType(e.target.value);
                  setSamePage(true);
                  const count = computeCurrentCategoryIndex(e.target.value);
                  setCurrentCategoryIndex(count + 1);
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
              />
              {llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE && (
                <>
                  <div className="labelWrapper">
                    <Label htmlFor="sitLocationInput">SIT Location</Label>
                  </div>
                  <select
                    label="SIT Location"
                    name="sitLocation"
                    id="sitLocationInput"
                    required
                    className={classnames('usa-select')}
                    value={ppmSITLocation}
                    onChange={(e) => {
                      handleSITLocationChange(e);
                    }}
                  >
                    {sitLocationOptions.map((x) => (
                      <option key={x.key}>{x.value}</option>
                    ))}
                  </select>
                  {/* <DropdownInput
                    label="SIT Location"
                    id="sitLocationInput"
                    name="sitLocation"
                    options={sitLocationOptions}
                    onChange={(e) => {
                      handleSITLocationChange(e);
                    }}
                  /> */}
                  <legend className={classnames('usa-label', styles.label)}>Cost</legend>
                  <div className={styles.displayValue}>
                    {toDollarString(formatCents(estimatedCost?.estimatedCost || 0))}
                  </div>
                </>
              )}
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
              {llvmExpenseTypes[selectedExpenseType] === expenseTypes.STORAGE && (
                <>
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
                  />
                  <MaskedTextField
                    defaultValue="0"
                    name="actualWeight"
                    label="Actual Weight"
                    id="actualWeight"
                    mask={Number}
                    scale={0} // digits after point, 0 for integers
                    signed={false} // disallow negative
                    thousandsSeparator=","
                    lazy={false} // immediate masking evaluation
                    suffix="lbs"
                  />
                  <DatePickerInput name="sitStartDate" label="Start date" />
                  <DatePickerInput name="sitEndDate" label="End date" />
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
                      <ErrorMessage display={!!errors?.reason && !!touched?.reason}>{errors.reason}</ErrorMessage>
                      <Textarea
                        id={`excludeReason-${expense?.id}`}
                        name="reason"
                        onChange={handleChange}
                        value={values.reason}
                        placeholder="Type something"
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
                      />
                      <div className={styles.hint}>{500 - values.reason.length} characters</div>
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
  tripNumber: number.isRequired,
  ppmNumber: number.isRequired,
  onSuccess: func,
  formRef: object,
};

ReviewExpense.defaultProps = {
  expense: undefined,
  onSuccess: null,
  formRef: null,
};
