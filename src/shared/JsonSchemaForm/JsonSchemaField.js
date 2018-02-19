import React, { Fragment } from 'react';
import { memoize } from 'lodash';

import { Field } from 'redux-form';
import './index.css';

const IS_REQUIRED_KEY = 'x-jsf-is-required';

const createDropDown = (fieldName, field, nameAttr) => {
  return (
    <Field name={nameAttr} component="select">
      <option />
      {field.enum.map(e => (
        <option key={e} value={e}>
          {field['x-display-value'][e]}
        </option>
      ))}
    </Field>
  );
};

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

const createNumberField = (fieldName, field, nameAttr) => {
  let validators = [numberValidator];
  if (field[IS_REQUIRED_KEY]) {
    validators.push(requiredValidator);
  }
  if (field.maximum != null) {
    validators.push(maximumValidator(field.maximum));
  }
  if (field.minimum != null) {
    validators.push(minimumValidator(field.minimum));
  }
  if (field.type === 'integer') {
    validators.push(integerValidator);
  }
  return (
    <Field
      component={renderInputField}
      name={nameAttr}
      parse={parseNumberField}
      type="number"
      validate={validators}
    />
  );
};

const createCheckbox = (fieldName, field, nameAttr) => {
  return (
    <Field id={fieldName} name={nameAttr} component="input" type="checkbox" />
  );
};

const phoneNumberValidator = value => {
  if (value && value.replace(/[^\d]/g, '').length !== 10) {
    return 'Number must have 10 digits.';
  }
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

const createTelephoneField = (fieldName, field, nameAttr) => {
  return (
    <Field
      name={nameAttr}
      component="input"
      type="text"
      placeholder="Phone Number"
      normalize={normalizePhone}
      validate={phoneNumberValidator}
    />
  );
};

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

const renderInputField = ({
  input,
  type,
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

const createTextAreaField = (fieldName, field, nameAttr) => {
  let validators = [];
  if (field[IS_REQUIRED_KEY]) {
    validators.push(requiredValidator);
  }
  if (field.maxLength) {
    validators.push(maxLengthValidator(field.maxLength));
  }
  if (field.minLength) {
    validators.push(minLengthValidator(field.minLength));
  }
  return (
    <Field
      id={fieldName}
      name={nameAttr}
      component={renderInputField}
      componentOverride={'textarea'}
      validate={validators}
    />
  );
};

const createInputField = (fieldName, field, nameAttr) => {
  return <Field name={nameAttr} component="input" type={field.format} />;
};

// Ok, so maybe our switches should be building up a single Field component params? right now this is a bit silly.
// keep at it and it will make sense.

// required is making the submit button not pop, that's a good first start
// need to indicate *if you go through the field* that it's required after.
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

  let fieldComponent;
  if (swaggerField.enum) {
    fieldComponent = createDropDown(fieldName, swaggerField, nameAttr);
  } else if (['integer', 'number'].includes(swaggerField.type)) {
    fieldComponent = createNumberField(fieldName, swaggerField, nameAttr);
  } else if (swaggerField.type === 'string') {
    if (swaggerField.format === 'textarea') {
      fieldComponent = createTextAreaField(fieldName, swaggerField, nameAttr);
    } else if (swaggerField.format === 'telephone') {
      fieldComponent = createTelephoneField(fieldName, swaggerField, nameAttr);
    } else {
      // more cases go here. Datetime, Date, UUID
      fieldComponent = createInputField(fieldName, swaggerField, nameAttr);
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
