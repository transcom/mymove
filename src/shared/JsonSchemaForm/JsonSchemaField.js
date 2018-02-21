import React, { Fragment } from 'react';

import validator from './validator';
import { Field } from 'redux-form';
import './index.css';

const IS_REQUIRED_KEY = 'x-jsf-is-required';

// ---- Validator -----

const parseNumberField = value => {
  if (!value || validator.isNumber(value)) {
    return value;
  } else {
    return parseFloat(value);
  }
};

// ----- Field configuration -----
const createCheckbox = (fieldName, field, nameAttr) => {
  return (
    <Field id={fieldName} name={nameAttr} component="input" type="checkbox" />
  );
};

const configureDropDown = (swaggerField, props) => {
  props.componentOverride = 'select';

  return props;
};

const dropDownChildren = (swaggerField, props) => {
  return (
    <Fragment>
      <option />
      {swaggerField.enum.map(e => (
        <option key={e} value={e}>
          {swaggerField['x-display-value'][e]}
        </option>
      ))}
    </Fragment>
  );
};

const configureNumberField = (swaggerField, props) => {
  props.type = 'number';
  props.step = 'any';
  props.parse = parseNumberField;

  if (swaggerField.maximum != null) {
    props.validate.push(validator.maximum(swaggerField.maximum));
  }
  if (swaggerField.minimum != null) {
    props.validate.push(validator.minimum(swaggerField.minimum));
  }
  if (swaggerField.type === 'integer') {
    props.validate.push(validator.isInteger);
  }

  return props;
};

const normalizePhone = (value, previousValue) => {
  if (!value) {
    return value;
  }
  const onlyNums = value.replace(/[^\d]/g, '');
  let normalizedPhone = '';
  for (let i = 0; i < 10; i++) {
    if (i >= onlyNums.length) {
      break;
    }
    if (i === 3 || i === 6) {
      normalizedPhone += '-';
    }
    normalizedPhone += onlyNums[i];
  }
  return normalizedPhone;
};

const configureTelephoneField = (swaggerField, props) => {
  props.normalize = normalizePhone;
  props.validate.push(validator.isPhoneNumber);
  props.type = 'text';

  return props;
};

const configureTextField = (swaggerField, props) => {
  if (swaggerField.maxLength) {
    props.validate.push(validator.maxLength(swaggerField.maxLength));
  }
  if (swaggerField.minLength) {
    props.validate.push(validator.minLength(swaggerField.minLength));
  }

  return props;
};

const renderInputField = ({
  input,
  type,
  step,
  componentOverride,
  meta: { touched, error, warning },
  children,
}) => {
  let componentName = 'input';
  if (componentOverride) {
    componentName = componentOverride;
  }

  const FieldComponent = React.createElement(
    componentName,
    {
      ...input,
      type: type,
      step: step,
    },
    children,
  );

  return (
    <div>
      {FieldComponent}
      {touched &&
        ((error && <span>{error}</span>) ||
          (warning && <span>{warning}</span>))}
    </div>
  );
};
// also should put a star by its name.

// This function switches on the type of the field and creates the correct
// Label and Field combination.
const createSchemaField = (fieldName, swaggerField, nameSpace) => {
  // Early return here, this is an edge case for label placement.
  // USWDS CSS only renders a checkbox if it is followed by its label
  const nameAttr = nameSpace ? `${nameSpace}.${fieldName}` : fieldName;
  if (swaggerField.type === 'boolean') {
    return (
      <Fragment key={fieldName}>
        {createCheckbox(fieldName, swaggerField, nameAttr)}
        <label htmlFor={fieldName}>{swaggerField.title || fieldName}</label>
      </Fragment>
    );
  }

  // Any field can be required
  let validators = [];
  if (swaggerField[IS_REQUIRED_KEY]) {
    validators.push(requiredValidator);
  }

  // configure the basic Field props
  let fieldProps = {};
  fieldProps.name = nameAttr;
  fieldProps.component = renderInputField;
  fieldProps.validate = validators;

  let children = null;

  if (swaggerField.enum) {
    fieldProps = configureDropDown(swaggerField, fieldProps);
    children = dropDownChildren(swaggerField);
  } else if (['integer', 'number'].includes(swaggerField.type)) {
    fieldProps = configureNumberField(swaggerField, fieldProps);
  } else if (swaggerField.type === 'string') {
    if (swaggerField.format === 'telephone') {
      fieldProps = configureTelephoneField(swaggerField, fieldProps);
      // more cases go here. Datetime, Date, SSN, (UUID)
    } else {
      // The last case is the simple text field / textarea which are the same but the componentOverride
      if (swaggerField.format === 'textarea') {
        fieldProps.componentOverride = 'textarea';
      }
      fieldProps = configureTextField(swaggerField, fieldProps);
    }
  } else {
    console.error(
      'ERROR: This is an unimplemented type in our JSONSchemaForm implmentation',
    );
  }

  return (
    <label key={fieldName}>
      {swaggerField.title || fieldName}
      <Field {...fieldProps}>{children}</Field>
    </label>
  );
};

export default {
  createSchemaField: createSchemaField,
  IS_REQUIRED_KEY: IS_REQUIRED_KEY,
};
