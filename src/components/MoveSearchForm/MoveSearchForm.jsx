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
    then: Yup.string().length(6, 'Move Code must be exactly 6 characters'),
    otherwise: Yup.string().when('searchType', {
      is: 'dodID',
      then: Yup.string().length(10, 'DOD ID must be exactly 10 characters'),
      otherwise: Yup.string().min(1, 'pls type something'),
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
            id={styles.SearchForm}
            onSubmit={formik.handleSubmit}
            role="search"
          >
            <p>What do you want to search for?</p>
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
            <div className={classnames(styles.searchBar)}>
              <TextField id="searchText" className="usa-search__input" label="Search" name="searchText" type="search" />
              <Button className={classnames(styles.searchButton)} type="submit" disabled={!formik.isValid}>
                Search
              </Button>
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
