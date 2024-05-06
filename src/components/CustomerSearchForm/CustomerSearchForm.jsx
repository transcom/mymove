import React from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, Radio } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import PropTypes from 'prop-types';

import styles from './CustomerSearchForm.module.scss';

import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import formStyles from 'styles/form.module.scss';

const baseSchema = Yup.object().shape({
  searchType: Yup.string().required('Select either DOD ID or Customer Name'),
});
const dodIDSchema = baseSchema.concat(
  Yup.object().shape({
    searchText: Yup.string().trim().length(10, 'DOD ID must be exactly 10 digits'),
  }),
);
const customerNameSchema = baseSchema.concat(
  Yup.object().shape({
    searchText: Yup.string().trim().min(1, 'Customer search must contain a value'),
  }),
);

const CustomerSearchForm = ({ onSubmit }) => {
  const getValidationSchema = (values) => {
    switch (values.searchType) {
      case 'dodID':
        return dodIDSchema;
      case 'customerName':
        return customerNameSchema;
      default:
        return Yup.object().shape({
          searchType: Yup.string().required('Search option must be selected'),
          searchText: Yup.string().required('Required'),
        });
    }
  };
  return (
    <Formik
      initialValues={{ searchType: 'dodID', searchText: '' }}
      onSubmit={onSubmit}
      validateOnChange
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
            className={classnames(formStyles.form, styles.CustomerSearchForm)}
            onSubmit={formik.handleSubmit}
            role="search"
          >
            <legend className="usa-label">What do you want to search for?</legend>
            <div role="group" className={formStyles.radioGroup}>
              <Field
                as={Radio}
                id="radio-picked-dodid"
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
            </div>
            <div className={styles.searchBar}>
              <TextField
                id="searchText"
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

CustomerSearchForm.propTypes = {
  onSubmit: PropTypes.func.isRequired,
};

export default CustomerSearchForm;
