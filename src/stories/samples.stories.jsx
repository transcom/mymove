import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { action } from '@storybook/addon-actions';
import { withState } from '@dump247/storybook-state';
import classNames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import { Form, EditButton } from '../components/form';

import TextField from 'components/form/fields/TextField';

const InlineForm = ({ name, label, initialValues, validationSchema, onSubmit, onReset, ...props }) => {
  const [show, setShow] = useState(false);
  const [errors, setErrors] = useState({});
  const errorCallback = (formErrors) => {
    setErrors(formErrors);
  };
  /* eslint-disable security/detect-object-injection */
  const errorMsg = errors[name];
  const value = initialValues[name];
  /* eslint-enable security/detect-object-injection */
  /* eslint-disable react/jsx-props-no-spreading */
  const content = show ? (
    <Formik
      onSubmit={(p) => {
        setShow(false);
        onSubmit(p);
      }}
      onReset={(p) => {
        setShow(false);
        onReset(p);
      }}
      initialValues={initialValues}
      validationSchema={validationSchema}
    >
      <Form errorCallback={errorCallback}>
        <TextField name={name} {...props} />
        <div className="display-flex">
          <Button type="submit">Submit</Button>
          <Button type="reset" secondary>
            Cancel
          </Button>
        </div>
      </Form>
    </Formik>
  ) : (
    <>
      <span>
        {errorMsg}
        {value}
      </span>
      <EditButton type="button" unstyled onClick={() => setShow(true)} />
    </>
  );
  /* eslint-enable react/jsx-props-no-spreading */
  return (
    <tr className="default-table-row-classes">
      <th htmlFor={name} className={classNames('default-table-header-class-names', { 'table-header-error': errorMsg })}>
        {label}
      </th>
      <td className="default-table-data-class-names">{content}</td>
    </tr>
  );
};

InlineForm.propTypes = {
  label: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  onSubmit: PropTypes.func.isRequired,
  onReset: PropTypes.func.isRequired,
  validationSchema: ({ validationSchema }, propName, componentName) => {
    if (!Object.keys(validationSchema).length) {
      return new Error(`Invalid prop ${propName} supplied to ${componentName}. Validation failed.`);
    }
    return null;
  },
  initialValues: ({ initialValues }, propName, componentName) => {
    if (!Object.keys(initialValues).length) {
      return new Error(`Invalid prop ${propName} supplied to ${componentName}. Validation failed.`);
    }
    return null;
  },
};

InlineForm.defaultProps = {
  validationSchema: {},
  initialValues: {},
};

export const personalInfo = () => (
  <div id="samples-orders-container" style={{ padding: '20px' }}>
    <div className="container container--accent--hhg">
      <Formik
        initialValues={{ firstName: '', lastName: '', email: '' }}
        validationSchema={Yup.object({
          firstName: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
          lastName: Yup.string().max(20, 'Must be 20 characters or less').required('Required'),
          email: Yup.string().email('Invalid email address').required('Required'),
        })}
        onSubmit={action('Form Submit')}
        onReset={action('Form Canceled')}
      >
        <Form>
          <TextField name="firstName" label="First Name" type="text" />
          <TextField name="lastName" label="Last Name" type="text" />
          <TextField name="email" label="Email Address" type="email" />
          <div className="display-flex">
            <Button type="submit">Submit</Button>
            <Button type="reset" secondary>
              Cancel
            </Button>
          </div>
        </Form>
      </Formik>
    </div>
  </div>
);

export const inlineFirstName = withState({ firstName: 'James', lastName: '' })(({ store }) => {
  const firstName = (
    <InlineForm
      initialValues={{ firstName: store.state.firstName }}
      validationSchema={Yup.object({
        firstName: Yup.string().max(15, 'Must be 15 characters or less').required('Required'),
      })}
      onSubmit={(formData) => {
        store.set(formData);
        action('First Name Form Submit')(formData);
      }}
      onReset={(formData) => {
        store.set(formData);
        action('First Name Form Canceled')(formData);
      }}
      name="firstName"
      type="text"
      label="First Name"
    />
  );
  const lastName = (
    <InlineForm
      initialValues={{ lastName: store.state.lastName }}
      validationSchema={Yup.object({
        lastName: Yup.string().max(20, 'Must be 20 characters or less').required('Required'),
      })}
      onSubmit={(formData) => {
        store.set(formData);
        action('Last Name Form Submit')(formData);
      }}
      onReset={(formData) => {
        store.set(formData);
        action('Last Name Form Canceled')(formData);
      }}
      name="lastName"
      type="text"
      label="Last Name"
    />
  );

  return (
    <div id="samples-orders-container" style={{ padding: '20px' }}>
      <div className="table--stacked table--stacked-wbuttons">
        <div className="display-flex">
          <div>
            <h4>Releasing Agent Info</h4>
          </div>
        </div>
        <table className="default-table-classes">
          <tbody>
            {firstName}
            {lastName}
          </tbody>
        </table>
      </div>
    </div>
  );
});

export default { title: 'Samples/Form' };
