import React from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, Form, FormGroup, Radio } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import * as Yup from 'yup';

import styles from './AboutForm.module.scss';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { DatePickerInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import formStyles from 'styles/form.module.scss';
import { ShipmentShape } from 'types/shipment';
import { uploadShape } from 'types/uploads';

const validationSchema = Yup.object().shape({
  receiptType: Yup.string().oneOf([
    'Contracted expense, Oil, Packing materials, Rental equipment, Storage, Tolls, Weighing fee, Other',
  ]),
  description: Yup.string().required('Required'),
  paidWithGTCC: Yup.boolean().required('Required'),
  amount: Yup.number().required('Required'),
  noReceipt: Yup.boolean().required('Required'),
  receiptDocument: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  startDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('receiptType', {
      is: 'Storage',
      then: (schema) => schema.required('Required'),
    }),
  endDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .when('receiptType', {
      is: 'Storage',
      then: (schema) => schema.required('Required'),
    }),
});

const ExpenseForm = ({ expense, onBack, onSubmit }) => {
  const {
    receiptType,
    description,
    paidWithGTCC,
    amount,
    noReceipt,
    receiptDocument,
    startDate,
    endDate
  } = expense || {};

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

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <div className={classnames(ppmStyles.formContainer)}>
            <Form className={classnames(formStyles.form, ppmStyles.form)}>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Departure date</h2>
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
  mtoShipment: ShipmentShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

export default ExpenseForm;
