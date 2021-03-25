import React from 'react';
import PropTypes from 'prop-types';
import { Field, useField } from 'formik';
import { Checkbox } from '@trussworks/react-uswds';

/**
 * This component renders a checkbox
 *
 * It relies on the Formik useField hook to work, so it must ALWAYS be rendered
 * inside of a Formik form context.
 *
 * If you want to use these components outside a Formik form, you can use the
 * ReactUSWDS components directly.
 */

const CheckboxField = ({ name, id, label, ...inputProps }) => {
  const [fieldProps] = useField({ name, type: 'checkbox' });

  /* eslint-disable-next-line react/jsx-props-no-spreading */
  return <Field id={id} as={Checkbox} name={name} label={label} {...fieldProps} {...inputProps} />;
};

CheckboxField.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  label: PropTypes.node.isRequired,
};

export default CheckboxField;
