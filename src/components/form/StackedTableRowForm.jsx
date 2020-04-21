import React from 'react';
import PropTypes from 'prop-types';

import { Formik } from 'formik';
import { Button, Label } from '@trussworks/react-uswds';
import classNames from 'classnames';
import { EditButton } from './IconButtons';
import { Form } from './Form';
import { TextInputMinimal } from './fields';

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
        <TextInputMinimal name={name} {...props} />
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
      <th scope="row" className={classNames('default-table-header-class-names', { 'table-header-error': errorMsg })}>
        {show ? (
          <Label error={errorMsg} htmlFor={name}>
            {label}
          </Label>
        ) : (
          <>{label}</>
        )}
      </th>
      <td className="default-table-data-class-names">{content}</td>
    </tr>
  );
};

StackedTableRowForm.propTypes = {
  name: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  onSubmit: PropTypes.func,
  onReset: PropTypes.func,
  // following are passed directly to formik
  // eslint-disable-next-line react/forbid-prop-types
  validationSchema: PropTypes.object.isRequired,
  // eslint-disable-next-line react/forbid-prop-types
  initialValues: PropTypes.object.isRequired,
};

export default StackedTableRowForm;
