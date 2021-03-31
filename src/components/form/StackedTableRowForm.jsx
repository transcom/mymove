import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import classNames from 'classnames/bind';

import { EditButton } from './IconButtons';
import { Form } from './Form';
import { ErrorMessage } from './ErrorMessage';
import styles from './StackedTableRowForm.module.scss';

import TextField from 'components/form/fields/TextField';

const cx = classNames.bind(styles);

export const StackedTableRowForm = ({ label, name, validationSchema, initialValues, onSubmit, onReset, ...props }) => {
  const [show, setShow] = React.useState(false);
  const [errors, setErrors] = React.useState({});
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
        <div className={cx('form-buttons')}>
          <Button type="submit">Submit</Button>
          <Button type="reset" secondary>
            Cancel
          </Button>
        </div>
      </Form>
    </Formik>
  ) : (
    <>
      <ErrorMessage className="display-inline" display={!!errorMsg}>
        {errorMsg}
      </ErrorMessage>
      <span>{value || '\u00A0'}</span>
      <EditButton type="button" className="float-right" unstyled onClick={() => setShow(true)} />
    </>
  );
  /* eslint-enable react/jsx-props-no-spreading */
  return (
    <tr className={cx('stacked-table-row')}>
      <th scope="row" className={`${cx('label')} ${classNames({ error: errorMsg })}`}>
        {label}
      </th>
      <td>{content}</td>
    </tr>
  );
};

StackedTableRowForm.propTypes = {
  name: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  onSubmit: PropTypes.func.isRequired,
  onReset: PropTypes.func,
  // following are passed directly to formik
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

StackedTableRowForm.defaultProps = {
  onReset: () => {},
  validationSchema: {},
  initialValues: {},
};

export default StackedTableRowForm;
