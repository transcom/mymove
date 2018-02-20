import React, { Fragment } from 'react';
import { memoize } from 'lodash';

import { Field } from 'redux-form';
import './index.css';

const IS_REQUIRED_KEY = 'x-jsf-is-required';

// ---- Validators -----

const requiredValidator = value => (value ? undefined : 'Required');
// Why Memoize? Please see https://github.com/erikras/redux-form/issues/3288
// Since we attach validators inside the render method, without memoization the
// function is re-created on every render which is not handled by react form.
// By memoizing it, it works.
const maxLengthValidator = memoize(maxLength => value => {
  if (value && value.length > maxLength) {
    return `Cannot exceed ${maxLength} characters.`;
  }
});
const minLengthValidator = memoize(minLength => value => {
  if (value && value.length < minLength) {
    return `Must be at least ${minLength} characters long.`;
  }
});

const maximumValidator = memoize(maximum => value => {
  if (value && value > maximum) {
    return `Must be ${maximum} or less`;
  }
});
const minimumValidator = memoize(minimum => value => {
  if (value && value < minimum) {
    return `Must be ${minimum} or more`;
  }
});

const numberValidator = value => {
  if (value) {
    if (isNaN(parseFloat(value))) {
      return 'Must be a number.';
    }
  }
};

const integerValidator = value => {
  if (value) {
    if (!Number.isInteger(value)) {
      return 'Must be an integer';
    }
  }
};

const parseNumberField = value => {
  if (!value || numberValidator(value)) {
    return value;
  } else {
    return parseFloat(value);
  }
};

const phoneNumberValidator = value => {
  if (value && value.replace(/[^\d]/g, '').length !== 10) {
    return 'Number must have 10 digits.';
  }
};

// ----- Field configuration -----
const createCheckbox = (fieldName, field, nameAttr) => {
  return (
    <Field id={fieldName} name={nameAttr} component="input" type="checkbox" />
  );
};

const configureDropDown = (swaggerField, props) => {
  props.component = 'select';

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
    props.validate.push(maximumValidator(swaggerField.maximum));
  }
  if (swaggerField.minimum != null) {
    props.validate.push(minimumValidator(swaggerField.minimum));
  }
  if (swaggerField.type === 'integer') {
    props.validate.push(integerValidator);
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
  props.validate.push(phoneNumberValidator);
  props.type = 'text';

  return props;
};

const configureTextField = (swaggerField, props) => {
  if (swaggerField.maxLength) {
    props.validate.push(maxLengthValidator(swaggerField.maxLength));
  }
  if (swaggerField.minLength) {
    props.validate.push(minLengthValidator(swaggerField.minLength));
  }

  return props;
};

const renderInputField = ({
  input,
  type,
  step,
  componentOverride,
  meta: { touched, error, warning },
}) => {
  let componentName = 'input';
  if (componentOverride) {
    componentName = componentOverride;
  }

  const FieldComponent = React.createElement(componentName, {
    ...input,
    type: type,
    step: step,
  });

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

  let fieldComponent;
  if (swaggerField.enum) {
    fieldProps = configureDropDown(swaggerField, fieldProps);
    children = dropDownChildren(swaggerField);

    fieldComponent = React.createElement(Field, fieldProps, children);
  } else if (['integer', 'number'].includes(swaggerField.type)) {
    fieldProps = configureNumberField(swaggerField, fieldProps);

    fieldComponent = React.createElement(Field, fieldProps);
  } else if (swaggerField.type === 'string') {
    if (swaggerField.format === 'telephone') {
      fieldProps = configureTelephoneField(swaggerField, fieldProps);

      fieldComponent = React.createElement(Field, fieldProps);
      // more cases go here. Datetime, Date, SSN, (UUID)
    } else {
      // The last case is the simple text field / textarea which are the same but the componentOverride
      if (swaggerField.format === 'textarea') {
        fieldProps.componentOverride = 'textarea';
      }
      fieldProps = configureTextField(swaggerField, fieldProps);

      fieldComponent = React.createElement(Field, fieldProps);
    }
  } else {
    console.error(
      'ERROR: This is an unimplemented type in our JSONSchemaForm implmentation',
    );
  }

  return (
    <label key={fieldName}>
      {swaggerField.title || fieldName}
      {fieldComponent}
    </label>
  );
};

export default {
  createSchemaField: createSchemaField,
  IS_REQUIRED_KEY: IS_REQUIRED_KEY,
};
