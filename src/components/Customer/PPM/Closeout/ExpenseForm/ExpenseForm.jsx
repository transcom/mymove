import React, { createRef } from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, ErrorMessage, Form, FormGroup, Radio, Label, Alert } from '@trussworks/react-uswds';
import { func, number } from 'prop-types';
import * as Yup from 'yup';

import styles from './ExpenseForm.module.scss';

import { formatCents } from 'utils/formatters';
import numOfDaysBetweenDates from 'utils/dates';
import { ppmExpenseTypes } from 'constants/ppmExpenseTypes';
import { ExpenseShape } from 'types/shipment';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
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

const validationSchema = Yup.object().shape({
  expenseType: Yup.string().required('Required'),
  description: Yup.string().required('Required'),
  paidWithGTCC: Yup.boolean().required('Required'),
  amount: Yup.string().notOneOf(['0', '0.00'], 'Please enter a non-zero amount').required('Required'),
  missingReceipt: Yup.boolean().required('Required'),
  document: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  sitStartDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('expenseType', {
      is: ppmExpenseTypes.STORAGE,
      then: (schema) => schema.required('Required').max(Yup.ref('sitEndDate'), 'Start date must be before end date.'),
    }),
  sitEndDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('expenseType', {
      is: ppmExpenseTypes.STORAGE,
      then: (schema) => schema.required('Required'),
    }),
});

const ExpenseForm = ({
  expense,
  receiptNumber,
  onBack,
  onSubmit,
  onCreateUpload,
  onUploadComplete,
  onUploadDelete,
}) => {
  const { movingExpenseType, description, paidWithGtcc, amount, missingReceipt, document, sitStartDate, sitEndDate } =
    expense || {};

  const initialValues = {
    expenseType: movingExpenseType || '',
    description: description || '',
    paidWithGTCC: paidWithGtcc ? 'true' : 'false',
    amount: amount ? `${formatCents(amount)}` : '',
    missingReceipt: !!missingReceipt,
    document: document?.uploads || [],
    sitStartDate: sitStartDate || '',
    sitEndDate: sitEndDate || '',
  };

  const documentRef = createRef();

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values, errors, ...formikProps }) => {
        return (
          <div className={classnames(ppmStyles.formContainer)}>
            <Form className={classnames(formStyles.form, ppmStyles.form, styles.ExpenseForm)}>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>{`Receipt ${receiptNumber}`}</h2>
                <FormGroup>
                  <DropdownInput label="Select type" name="expenseType" options={ppmExpenseTypes} id="expenseType" />
                </FormGroup>
                {values.expenseType && (
                  <>
                    <FormGroup>
                      <h3>Description</h3>
                      <TextField label="What did you buy?" id="description" name="description" />
                      <Hint>Add a brief description of the expense.</Hint>
                      <Fieldset>
                        <legend className="usa-label">
                          Did you pay with your GTCC (Government Travel Charge Card)?
                        </legend>
                        <Field
                          as={Radio}
                          id="yes-used-gtcc"
                          label="Yes"
                          name="paidWithGTCC"
                          value="true"
                          checked={values.paidWithGTCC === 'true'}
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
                      />
                      <Hint>
                        Enter the total unit price for all items on the receipt that you&apos;re claiming as part of
                        your PPM moving expenses.
                      </Hint>
                      <CheckboxField id="missingReceipt" name="missingReceipt" label="I don't have this receipt" />
                      {values.missingReceipt && (
                        <Alert type="info">
                          {`If you can, get a replacement copy of your receipt and upload that. \nIf that is not possible, write and sign a statement that explains why this receipt is missing. Include details about where and when you purchased this item. Upload that statement. Your reimbursement for this expense will be based on the information you provide.`}
                        </Alert>
                      )}
                      <div className={styles.labelWrapper}>
                        <Label error={formikProps.touched?.document && formikProps.errors?.document} htmlFor="document">
                          Upload receipt
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
                          documentRef.current.removeFile(upload.id);
                        }}
                        ref={documentRef}
                      />
                    </FormGroup>
                  </>
                )}
                {values.expenseType === 'STORAGE' && (
                  <FormGroup>
                    <h3>Dates</h3>
                    <DatePickerInput name="sitStartDate" label="Start date" />
                    <DatePickerInput name="sitEndDate" label="End date" />
                    <h3>
                      Days in storage:{' '}
                      {values.sitStartDate && values.sitEndDate && !errors.sitStartDate && !errors.sitEndDate
                        ? 1 + numOfDaysBetweenDates(values.sitStartDate, values.sitEndDate)
                        : ''}
                    </h3>
                  </FormGroup>
                )}
              </SectionWrapper>
              <div className={ppmStyles.buttonContainer}>
                <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Return To Homepage
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
