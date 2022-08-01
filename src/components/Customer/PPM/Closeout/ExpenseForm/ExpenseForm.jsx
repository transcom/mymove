import React, { createRef } from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, ErrorMessage, Form, FormGroup, Radio, Label } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import * as Yup from 'yup';

// import styles from './AboutForm.module.scss';

import { dropdownInputOptions } from 'utils/formatters';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import FileUpload from 'components/FileUpload/FileUpload';
import formStyles from 'styles/form.module.scss';
// import { ShipmentShape } from 'types/shipment';
import { uploadShape } from 'types/uploads';
import { CheckboxField, DatePickerInput, DropdownInput } from 'components/form/fields';
import { DocumentAndImageUploadInstructions, UploadDropZoneLabel, UploadDropZoneLabelMobile } from 'content/uploads';
import UploadsTable from 'components/UploadsTable/UploadsTable';

const validationSchema = Yup.object().shape({
  receiptType: Yup.string().oneOf([
    'Contracted expense, Oil, Packing materials, Rental equipment, Storage, Tolls, Weighing fee, Other',
  ]),
  description: Yup.string().required('Required'),
  paidWithGTCC: Yup.boolean().required('Required'),
  amount: Yup.number().required('Required'),
  noReceipt: Yup.boolean().required('Required'),
  receiptDocument: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  sitStartDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('receiptType', {
      is: 'Storage',
      then: (schema) => schema.required('Required'),
    }),
  sitEndDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('receiptType', {
      is: 'Storage',
      then: (schema) => schema.required('Required'),
    }),
});

const ExpenseForm = ({ expense, onBack, onSubmit, onCreateUpload, onUploadComplete, onUploadDelete }) => {
  const { receiptType, description, paidWithGTCC, amount, noReceipt, receiptDocument, startDate, endDate } =
    expense || {};

  const initialValues = {
    receiptType: receiptType || '',
    description: description || '',
    paidWithGTCC: paidWithGTCC ? 'true' : 'false',
    amount: amount ? amount.toString() : '',
    noReceipt: noReceipt ? 'true' : 'false',
    receiptDocument: receiptDocument || [],
    startDate: startDate || '',
    endDate: endDate || '',
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
      {({ isValid, isSubmitting, handleSubmit, values, ...formikProps }) => {
        return (
          <div className={classnames(ppmStyles.formContainer)}>
            <Form className={classnames(formStyles.form, ppmStyles.form)}>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                {/* TODO: ADD WEIGHT TICKET NUMBER */}
                <DropdownInput label="Select type" name="receiptType" options={expenseOptions} id="receiptType" />
                <h2>Description</h2>
                <TextField label="What did you buy?" id="description" name="description" />
                <Hint className={ppmStyles.hint}>Add a brief description of the expense.</Hint>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">Did you pay with your GTCC? (Government Travel Charge Card)</legend>
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
                  <h2>Amount</h2>
                  <MaskedTextField
                    defaultValue="0"
                    name="advanceAmountRequested"
                    label="Amount requested"
                    id="advanceAmountRequested"
                    mask={Number}
                    scale={0} // digits after point, 0 for integers
                    signed={false} // disallow negative
                    thousandsSeparator=","
                    lazy={false} // immediate masking evaluation
                    prefix="$"
                    hintClassName={ppmStyles.innerHint}
                  />
                  <Hint>
                    Enter the total unit price for all items on the receipt that you&apos;re claiming as part of your
                    PPM moving expenses.
                  </Hint>
                  <CheckboxField id="missingReceipt" name="missingReceipt" label="I don't have this receipt" />
                  <div className="labelWrapper">
                    <Label
                      error={
                        formikProps.touched?.proofOfTrailerOwnershipDocument &&
                        formikProps.errors?.proofOfTrailerOwnershipDocument
                      }
                      htmlFor="receiptDocument"
                    >
                      Upload receipt
                    </Label>
                  </div>
                  {formikProps.touched?.proofOfTrailerOwnershipDocument &&
                    formikProps.errors?.proofOfTrailerOwnershipDocument && (
                      <ErrorMessage>{formikProps.errors?.proofOfTrailerOwnershipDocument}</ErrorMessage>
                    )}
                  <Hint>
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
                {values.receiptType === 'storage' && (
                  <FormGroup>
                    <h2>Dates</h2>
                    <DatePickerInput name="sitStartDate" label="Start date" />
                    <DatePickerInput name="sitEndDate" label="End date" />
                    <h3>Days in storage:</h3>
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
  // mtoShipment: ShipmentShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  onCreateUpload: func.isRequired,
  onUploadComplete: func.isRequired,
  onUploadDelete: func.isRequired,
};

export default ExpenseForm;
