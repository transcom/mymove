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

const validationSchema = Yup.object().shape({
  searchType: Yup.string().required('searchtype error'),
  searchText: Yup.string().when('searchType', {
    is: 'moveCode',
    then: (schema) => schema.length(6, 'Move Code must be exactly 6 characters'),
    otherwise: (schema) =>
      schema.when('searchType', {
        is: 'dodID',
        then: (s) => s.length(10, 'DOD ID must be exactly 10 characters'),
        otherwise: (s) => s.min(1, 'Search must contain at least one character'),
      }),
  }),
});

const MoveSearchForm = ({ onSubmit }) => {
  return (
    <Formik
      initialValues={{ searchType: 'moveCode', searchText: '' }}
      onSubmit={onSubmit}
      validationSchema={validationSchema}
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
                type="radio"
                name="searchType"
                value="moveCode"
                title="Move Code"
                label="Move Code"
              />
              <Field
                as={Radio}
                id="radio-picked-dodid"
                type="radio"
                name="searchType"
                value="dodID"
                title="DOD ID"
                label="DOD ID"
              />
              <Field
                as={Radio}
                id="radio-picked-customername"
                type="radio"
                name="searchType"
                value="customerName"
                title="Customer Name"
                label="Customer Name"
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
                  <Button className={styles.searchButton} type="submit" disabled={!formik.isValid}>
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
