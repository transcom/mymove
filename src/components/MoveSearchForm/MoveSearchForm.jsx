import React from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, Radio } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import PropTypes from 'prop-types';

import styles from './MoveSearchForm.module.scss';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import formStyles from 'styles/form.module.scss';
import { roleTypes } from 'constants/userRoles';

const baseSchema = Yup.object().shape({
  searchType: Yup.string().required('searchtype error'),
});
const moveCodeSchema = baseSchema.concat(
  Yup.object().shape({
    searchText: Yup.string().trim().length(6, 'Move Code must be exactly 6 characters'),
  }),
);
const dodIDSchema = baseSchema.concat(
  Yup.object().shape({
    searchText: Yup.string().trim().length(10, 'DOD ID must be exactly 10 characters'),
  }),
);
const customerNameSchema = baseSchema.concat(
  Yup.object().shape({
    searchText: Yup.string().trim().min(1, 'Customer search must contain a value'),
  }),
);
const paymentRequestCodeSchema = baseSchema.concat(
  Yup.object().shape({
    searchText: Yup.string()
      .trim()
      .length(
        11,
        'Payment request number must be a 11 character string (9 numbers with hyphens after every 4th number, e.g. 1234-5678-9)',
      )
      .matches(/(\d{4})(-{1})(\d{4})(-{1})(\d{1})/g, {
        message:
          'Payment request number must be a 11 character string (9 numbers with hyphens after every 4th number, e.g. 1234-5678-9)',
      }),
  }),
);

const MoveSearchForm = ({ onSubmit, role }) => {
  const getValidationSchema = (values) => {
    switch (values.searchType) {
      case 'moveCode':
        return moveCodeSchema;
      case 'dodID':
        return dodIDSchema;
      case 'customerName':
        return customerNameSchema;
      case 'paymentRequestCode':
        return paymentRequestCodeSchema;
      default:
        return Yup.object().shape({
          searchType: Yup.string().required('Search option must be selected'),
          searchText: Yup.string().required('Required'),
        });
    }
  };
  return (
    <Formik
      initialValues={{ searchType: 'moveCode', searchText: '' }}
      onSubmit={onSubmit}
      validateOnChange
      data-testid="move-search-form"
      // adding a return will break the validation
      // RA Validator Status: RA Accepted
      // eslint-disable-next-line consistent-return
      validate={(values) => {
        const schema = getValidationSchema(values);
        try {
          schema.validateSync(values, { abortEarly: false });
        } catch (error) {
          return error.inner.reduce((acc, { path, message }) => ({ ...acc, [path]: message }), {});
        }
      }}
    >
      {(formik) => {
        return (
          <Form
            className={classnames(formStyles.form, styles.MoveSearchForm)}
            onSubmit={formik.handleSubmit}
            role="search"
          >
            <legend className="usa-label">What do you want to search for?</legend>
            <div role="group" className={formStyles.radioGroup}>
              <Field
                as={Radio}
                id="radio-picked-movecode"
                data-testid="moveCode"
                type="radio"
                name="searchType"
                value="moveCode"
                title="Move Code"
                label="Move Code"
                onChange={(e) => {
                  formik.setFieldValue('searchType', e.target.value);
                  formik.setFieldValue('searchText', '', false); // Clear TextField
                  formik.setFieldTouched('searchText', false, false);
                }}
              />
              <Field
                as={Radio}
                id="radio-picked-dodid"
                data-testid="dodID"
                type="radio"
                name="searchType"
                value="dodID"
                title="DOD ID"
                label="DOD ID"
                onChange={(e) => {
                  formik.setFieldValue('searchType', e.target.value);
                  formik.setFieldValue('searchText', '', false); // Clear TextField
                  formik.setFieldTouched('searchText', false, false);
                }}
              />
              <Field
                as={Radio}
                id="radio-picked-customername"
                data-testid="customerName"
                type="radio"
                name="searchType"
                value="customerName"
                title="Customer Name"
                label="Customer Name"
                onChange={(e) => {
                  formik.setFieldValue('searchType', e.target.value);
                  formik.setFieldValue('searchText', '', false); // Clear TextField
                  formik.setFieldTouched('searchText', false, false);
                }}
              />
              {role !== roleTypes.SERVICES_COUNSELOR && (
                <Field
                  as={Radio}
                  id="radio-picked-paymentRequestCode"
                  data-testid="paymentRequestCode"
                  type="radio"
                  name="searchType"
                  value="paymentRequestCode"
                  title="Payment Request Number"
                  label="Payment Request Number"
                  onChange={(e) => {
                    formik.setFieldValue('searchType', e.target.value);
                    formik.setFieldValue('searchText', '', false); // Clear TextField
                    formik.setFieldTouched('searchText', false, false);
                  }}
                />
              )}
            </div>
            <div className={styles.searchBar}>
              <TextField
                id="searchText"
                data-testid="searchText"
                className="usa-search__input"
                label={<legend className="usa-label">Search</legend>}
                name="searchText"
                type="search"
                button={
                  <Button
                    data-testid="searchTextSubmit"
                    className={styles.searchButton}
                    type="submit"
                    disabled={!formik.isValid}
                  >
                    Search
                  </Button>
                }
              />
            </div>
          </Form>
        );
      }}
    </Formik>
  );
};

MoveSearchForm.propTypes = {
  onSubmit: PropTypes.func.isRequired,
};

export default MoveSearchForm;
