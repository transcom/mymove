import React, { createRef } from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, ErrorMessage, Form, FormGroup, Radio, Label, Alert } from '@trussworks/react-uswds';
import { func, number } from 'prop-types';
import * as Yup from 'yup';

import SmallPackageForm from '../SmallPackageForm/SmallPackageForm';

import styles from './ExpenseForm.module.scss';

import { formatCents } from 'utils/formatters';
import { numOfDaysBetweenDates } from 'utils/dates';
import { expenseTypes, ppmExpenseTypes } from 'constants/ppmExpenseTypes';
import { ExpenseShape } from 'types/shipment';
import ppmStyles from 'components/Shared/PPM/PPM.module.scss';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import FileUpload from 'components/FileUpload/FileUpload';
import formStyles from 'styles/form.module.scss';
import { uploadShape } from 'types/uploads';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import { DocumentAndImageUploadInstructions, UploadDropZoneLabel, UploadDropZoneLabelMobile } from 'content/uploads';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import { PPM_TYPES } from 'shared/constants';
import { APP_NAME } from 'constants/apps';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const validationSchema = Yup.object().shape({
  expenseType: Yup.string().required('Required'),
  description: Yup.string().when('expenseType', {
    is: (expenseType) => expenseType !== expenseTypes.SMALL_PACKAGE,
    then: (schema) => schema.required('Required'),
  }),
  paidWithGTCC: Yup.boolean().required('Required'),
  amount: Yup.string().notOneOf(['0', '0.00'], 'Please enter a non-zero amount').required('Required'),
  missingReceipt: Yup.boolean().required('Required'),
  document: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  sitStartDate: Yup.date()
    .nullable()
    .transform((value, originalValue) => (originalValue === '' ? null : value))
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('expenseType', {
      is: expenseTypes.STORAGE,
      then: (schema) => schema.required('Required').max(Yup.ref('sitEndDate'), 'Start date must be before end date.'),
    }),
  sitEndDate: Yup.date()
    .nullable()
    .transform((value, originalValue) => (originalValue === '' ? null : value))
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('expenseType', {
      is: expenseTypes.STORAGE,
      then: (schema) => schema.required('Required'),
    }),
  sitLocation: Yup.string().when('expenseType', {
    is: expenseTypes.STORAGE,
    then: (schema) => schema.required('Required'),
  }),
  sitWeight: Yup.number()
    .nullable()
    .transform((value, originalValue) => (originalValue === '' ? null : value))
    .when('expenseType', {
      is: expenseTypes.STORAGE,
      then: (schema) => schema.required('Required').moreThan(0, 'Weight stored must be at least 1 lb.'),
    }),
  weightShipped: Yup.number().when('expenseType', {
    is: expenseTypes.SMALL_PACKAGE,
    then: (schema) => schema.required('Required').min(0, 'Weight shipped must be at least 1 lb.'),
  }),
  isProGear: Yup.string().when('expenseType', {
    is: expenseTypes.SMALL_PACKAGE,
    then: (schema) => schema.required('Required'),
  }),
  proGearBelongsToSelf: Yup.string()
    .nullable()
    .when('isProGear', {
      is: 'true',
      then: (schema) => schema.required('Required'),
      otherwise: (schema) => schema.strip(),
    }),
  proGearDescription: Yup.string().when('isProGear', {
    is: 'true',
    then: (schema) => schema.required('Required'),
    otherwise: (schema) => schema.strip(),
  }),
});

const ExpenseForm = ({
  ppmType,
  expense,
  receiptNumber,
  onBack,
  onSubmit,
  onCreateUpload,
  onUploadComplete,
  onUploadDelete,
  appName,
}) => {
  const {
    movingExpenseType,
    description,
    paidWithGtcc,
    amount,
    missingReceipt,
    document,
    sitStartDate,
    sitEndDate,
    sitLocation,
    weightStored,
    trackingNumber,
    weightShipped,
    isProGear,
    proGearBelongsToSelf,
    proGearDescription,
  } = expense || {};

  const initialValues = {
    expenseType:
      !movingExpenseType && ppmType === PPM_TYPES.SMALL_PACKAGE ? expenseTypes.SMALL_PACKAGE : movingExpenseType,
    description: description || '',
    paidWithGTCC: paidWithGtcc ? 'true' : 'false',
    amount: amount ? `${formatCents(amount)}` : '',
    missingReceipt: !!missingReceipt,
    document: document?.uploads || [],
    sitStartDate: sitStartDate || '',
    sitEndDate: sitEndDate || '',
    sitLocation: sitLocation || undefined,
    sitWeight: weightStored ? `${weightStored}` : '',
    trackingNumber: trackingNumber || '',
    weightShipped: weightShipped ? `${weightShipped}` : '',
    isProGear: isProGear ? 'true' : 'false',
    ...(isProGear && {
      proGearBelongsToSelf: proGearBelongsToSelf ? 'true' : 'false',
      proGearDescription: proGearDescription || '',
    }),
  };

  const documentRef = createRef();

  const availableExpenseTypes =
    ppmType === PPM_TYPES.SMALL_PACKAGE
      ? [{ value: 'Small package reimbursement', key: 'SMALL_PACKAGE' }]
      : ppmExpenseTypes;

  const isCustomerPage = appName === APP_NAME.MYMOVE;

  return (
    <>
      <div className={styles.introSection}>
        <p>
          Document your qualified expenses by uploading receipts. They should include a description of the item, the
          price you paid, the date of purchase, and the business name. All documents must be legible and unaltered.
        </p>
        <p>Your finance office will make the final decision about which expenses are deductible or reimbursable.</p>
        <p>Upload one receipt at a time. Please do not put multiple receipts in one image.</p>
      </div>
      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {({ isValid, isSubmitting, handleSubmit, values, errors, ...formikProps }) => {
          return (
            <div className={classnames(ppmStyles.formContainer)}>
              <Form className={classnames(formStyles.form, ppmStyles.form, styles.ExpenseForm)}>
                <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                  <h2>
                    {ppmType !== PPM_TYPES.SMALL_PACKAGE ? `Receipt ` : `Package `}
                    {receiptNumber}
                  </h2>
                  {values.expenseType === expenseTypes.SMALL_PACKAGE && (
                    <Hint data-testid="smallPackageInfo">
                      Receipts from the package carrier should include the weight, cost, and tracking number (optional).
                      Receipts must be legible and unaltered. Files must be 25MB or smaller. You must upload at least
                      one package carrier receipt to get paid for your Small Package Reimbursement PPM.
                    </Hint>
                  )}
                  {requiredAsteriskMessage}
                  <FormGroup className={styles.dropdown}>
                    <DropdownInput
                      label="Select type"
                      name="expenseType"
                      options={availableExpenseTypes}
                      id="expenseType"
                      isDisabled={ppmType === PPM_TYPES.SMALL_PACKAGE}
                      showRequiredAsterisk
                      required
                    />
                  </FormGroup>
                  {values.expenseType && (
                    <>
                      <FormGroup>
                        {values.expenseType !== expenseTypes.SMALL_PACKAGE && (
                          <>
                            <h3>Description</h3>
                            <TextField
                              label="What did you buy or rent?"
                              id="description"
                              name="description"
                              showRequiredAsterisk
                              required
                            />
                            <Hint>Add a brief description of the expense.</Hint>
                          </>
                        )}
                        {values.expenseType === expenseTypes.STORAGE && (
                          <FormGroup>
                            <legend className="usa-label" aria-label="Required: Where did you store your items?">
                              <span required>
                                Where did you store your items? <RequiredAsterisk />
                              </span>
                            </legend>
                            <Field
                              as={Radio}
                              id="sitLocationOrigin"
                              label="Origin"
                              name="sitLocation"
                              value="ORIGIN"
                              checked={values.sitLocation === 'ORIGIN'}
                            />
                            <Field
                              as={Radio}
                              id="sitLocationDestination"
                              label="Destination"
                              name="sitLocation"
                              value="DESTINATION"
                              checked={values.sitLocation === 'DESTINATION'}
                            />
                            <MaskedTextField
                              defaultValue="0"
                              name="sitWeight"
                              label="Weight Stored"
                              id="sitWeightInput"
                              mask={Number}
                              scale={0} // digits after point, 0 for integers
                              signed={false} // disallow negative
                              thousandsSeparator=","
                              lazy={false} // immediate masking evaluation
                              showRequiredAsterisk
                              required
                            >
                              {'  '} lbs
                            </MaskedTextField>
                            <Hint>Enter the weight of the items that were stored during your PPM.</Hint>
                          </FormGroup>
                        )}

                        <Fieldset>
                          <legend
                            className="usa-label"
                            aria-label="Required: Did you pay with your GTCC (Government Travel Charge Card)?"
                          >
                            <span required>
                              Did you pay with your GTCC (Government Travel Charge Card)? <RequiredAsterisk />
                            </span>
                          </legend>
                          <Field
                            as={Radio}
                            id="yes-used-gtcc"
                            label="Yes"
                            name="paidWithGTCC"
                            value="true"
                            checked={values.paidWithGTCC === 'true'}
                            showRequiredAsterisk
                            required
                          />
                          <Field
                            as={Radio}
                            id="no-did-not-use-gtcc"
                            label="No"
                            name="paidWithGTCC"
                            value="false"
                            checked={values.paidWithGTCC === 'false'}
                          />
                        </Fieldset>
                      </FormGroup>
                      <FormGroup>
                        {values.expenseType !== expenseTypes.SMALL_PACKAGE ? (
                          <>
                            <h3>Amount</h3>
                            <MaskedTextField
                              name="amount"
                              label="Amount"
                              id="amount"
                              mask={Number}
                              scale={2} // digits after point, 0 for integers
                              signed={false} // disallow negative
                              radix="." // fractional delimiter
                              mapToRadix={['.']} // symbols to process as radix
                              padFractionalZeros // if true, then pads zeros at end to the length of scale
                              thousandsSeparator=","
                              lazy={false} // immediate masking evaluation
                              prefix="$"
                              hintClassName={ppmStyles.innerHint}
                              showRequiredAsterisk
                              required
                            />
                            <Hint>
                              Enter the total unit price for all items on the receipt that you&apos;re claiming as part
                              of your PPM moving expenses.
                            </Hint>
                          </>
                        ) : (
                          <SmallPackageForm />
                        )}
                        <CheckboxField id="missingReceipt" name="missingReceipt" label="I don't have this receipt" />
                        {values.missingReceipt && values.expenseType === expenseTypes.SMALL_PACKAGE && (
                          <Alert type="info" className={styles.uploadInstructions}>
                            {values.expenseType === expenseTypes.SMALL_PACKAGE &&
                              'If you do not upload legible package receipts your PPM reimbursement could be affected.'}
                          </Alert>
                        )}
                        {values.missingReceipt && (
                          <Alert type="info">
                            If you can, get a replacement copy of your receipt and upload that. If that is not possible,
                            write and sign a statement that explains why this receipt is missing. Include details about
                            where and when you purchased this item. Upload that statement. Your reimbursement for this
                            expense will be based on the information you provide.
                          </Alert>
                        )}
                        <div className={styles.labelWrapper}>
                          <Label
                            error={formikProps.touched?.document && formikProps.errors?.document}
                            htmlFor="document"
                          >
                            <span>
                              Upload receipt <RequiredAsterisk />
                            </span>
                          </Label>
                        </div>
                        {formikProps.touched?.document && formikProps.errors?.document && (
                          <ErrorMessage>{formikProps.errors?.document}</ErrorMessage>
                        )}
                        <Hint className={styles.uploadInstructions}>
                          <p>{DocumentAndImageUploadInstructions}</p>
                        </Hint>
                        <UploadsTable
                          uploads={values.document}
                          onDelete={(uploadId) =>
                            onUploadDelete(uploadId, 'document', formikProps.setFieldTouched, formikProps.setFieldValue)
                          }
                        />
                        <FileUpload
                          name="document"
                          className="receiptDocument"
                          createUpload={(file) => onCreateUpload('document', file, formikProps.setFieldTouched)}
                          labelIdle={UploadDropZoneLabel}
                          labelIdleMobile={UploadDropZoneLabelMobile}
                          onChange={(err, upload) => {
                            formikProps.setFieldTouched('document', true);
                            onUploadComplete(err);
                            documentRef?.current?.removeFile(upload.id);
                          }}
                          ref={documentRef}
                        />
                      </FormGroup>
                    </>
                  )}
                  {values.expenseType === 'STORAGE' && (
                    <FormGroup>
                      <h3>Dates</h3>
                      <DatePickerInput name="sitStartDate" label="Start date" showRequiredAsterisk required />
                      <DatePickerInput name="sitEndDate" label="End date" showRequiredAsterisk required />
                      <h3 className={styles.storageTotal}>
                        Days in storage:{' '}
                        {values.sitStartDate && values.sitEndDate && !errors.sitStartDate && !errors.sitEndDate
                          ? 1 + numOfDaysBetweenDates(values.sitStartDate, values.sitEndDate)
                          : ''}
                      </h3>
                    </FormGroup>
                  )}
                </SectionWrapper>
                <div
                  className={`${
                    isCustomerPage ? ppmStyles.buttonContainer : `${formStyles.formActions} ${ppmStyles.buttonGroup}`
                  }`}
                >
                  <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                    Cancel
                  </Button>
                  <Button
                    className={ppmStyles.saveButton}
                    type="button"
                    onClick={handleSubmit}
                    disabled={!isValid || isSubmitting}
                  >
                    Save & Continue
                  </Button>
                </div>
              </Form>
            </div>
          );
        }}
      </Formik>
    </>
  );
};

ExpenseForm.propTypes = {
  receiptNumber: number,
  expense: ExpenseShape,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  onCreateUpload: func.isRequired,
  onUploadComplete: func.isRequired,
  onUploadDelete: func.isRequired,
};

ExpenseForm.defaultProps = {
  expense: undefined,
  receiptNumber: 1,
};

export default ExpenseForm;
