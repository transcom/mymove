import React, { createRef } from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, ErrorMessage, Form, FormGroup, Radio, Label, Alert } from '@trussworks/react-uswds';
import { func, string } from 'prop-types';
import * as Yup from 'yup';

import styles from './ExpenseForm.module.scss';

import numOfDaysBetweenDates from 'utils/dates';
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
  amount: Yup.number().required('Required'),
  missingReceipt: Yup.boolean().required('Required'),
  receiptDocument: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  sitStartDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('expenseType', {
      is: 'storage',
      then: (schema) => schema.required('Required'),
    }),
  sitEndDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('expenseType', {
      is: 'storage',
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
  const { expenseType, description, paidWithGTCC, amount, missingReceipt, receiptDocument, sitStartDate, sitEndDate } =
    expense || {};

  const initialValues = {
    expenseType: expenseType || '',
    description: description || '',
    paidWithGTCC: paidWithGTCC ? 'true' : 'false',
    amount: amount ? `${amount}` : '',
    missingReceipt: !!missingReceipt,
    receiptDocument: receiptDocument?.uploads || [],
    sitStartDate: sitStartDate || '',
    sitEndDate: sitEndDate || '',
  };

  const receiptDocumentRef = createRef();
  const expenseOptions = [
    { value: 'Contracted expense', key: 'contracted_expense' },
    { value: 'Oil', key: 'oil' },
    { value: 'Packing materials', key: 'packing_materials' },
    { value: 'Rental equipment', key: 'rental_equipment' },
    { value: 'Storage', key: 'storage' },
    { value: 'Tolls', key: 'tolls' },
    { value: 'Weighing fee', key: 'weighing_fee' },
    { value: 'Other', key: 'other' },
  ];
  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values, errors, ...formikProps }) => {
        return (
          <div className={classnames(ppmStyles.formContainer)}>
            <Form className={classnames(formStyles.form, ppmStyles.form, styles.ExpenseForm)}>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>{`Receipt ${receiptNumber}`}</h2>
                <FormGroup>
                  <DropdownInput label="Select type" name="expenseType" options={expenseOptions} id="expenseType" />
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
                        scale={0} // digits after point, 0 for integers
                        signed={false} // disallow negative
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
                        <Label
                          error={formikProps.touched?.receiptDocument && formikProps.errors?.receiptDocument}
                          htmlFor="receiptDocument"
                        >
                          Upload receipt
                        </Label>
                      </div>
                      {formikProps.touched?.receiptDocument && formikProps.errors?.receiptDocument && (
                        <ErrorMessage>{formikProps.errors?.receiptDocument}</ErrorMessage>
                      )}
                      <Hint className={styles.uploadInstructions}>
                        <p>{DocumentAndImageUploadInstructions}</p>
                      </Hint>
                      <UploadsTable
                        // className={styles.uploadsTable}
                        uploads={values.receiptDocument}
                        onDelete={(uploadId) =>
                          onUploadDelete(
                            uploadId,
                            'receiptDocument',
                            formikProps.setFieldTouched,
                            formikProps.setFieldValue,
                          )
                        }
                      />
                      <FileUpload
                        name="receiptDocument"
                        className="receiptDocument"
                        createUpload={(file) => onCreateUpload('receiptDocument', file)}
                        labelIdle={UploadDropZoneLabel}
                        labelIdleMobile={UploadDropZoneLabelMobile}
                        onChange={(err, upload) => {
                          formikProps.setFieldTouched('receiptDocument', true);
                          onUploadComplete(err);
                          receiptDocumentRef.current.removeFile(upload.id);
                        }}
                        ref={receiptDocumentRef}
                      />
                    </FormGroup>
                  </>
                )}
                {values.expenseType === 'storage' && (
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
                  Finish Later
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
  receiptNumber: string,
  expense: ExpenseShape,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  onCreateUpload: func.isRequired,
  onUploadComplete: func.isRequired,
  onUploadDelete: func.isRequired,
};

ExpenseForm.defaultProps = {
  expense: undefined,
  receiptNumber: '1',
};

export default ExpenseForm;
