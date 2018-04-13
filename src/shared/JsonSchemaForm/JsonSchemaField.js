import React, { Fragment } from 'react';

import validator from './validator';
import { Field } from 'redux-form';

export const ALWAYS_REQUIRED_KEY = 'x-always-required';

// ---- Parsers -----

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
  props.type = 'text';
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

  props.validate.push(validator.isNumber);

  return props;
};

const configureTelephoneField = (swaggerField, props) => {
  props.normalize = validator.normalizePhone;
  props.validate.push(
    validator.patternMatches(
      swaggerField.pattern,
      'Number must have 10 digits.',
    ),
  );
  props.type = 'text';

  return props;
};

const configureSSNField = (swaggerField, props) => {
  props.normalize = validator.normalizeSSN;
  props.validate.push(
    validator.patternMatches(swaggerField.pattern, 'SSN must have 9 digits.'),
  );
  props.type = 'text';

  return props;
};

const configureZipField = (swaggerField, props) => {
  props.normalize = validator.normalizeZip;
  props.validate.push(
    validator.patternMatches(
      swaggerField.pattern,
      'Zip code must have 5 or 9 digits.',
    ),
  );
  props.type = 'text';

  return props;
};

const configureDateField = (swaggerField, props) => {
  props.type = 'date';

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

const configureEmailField = (swaggerField, props) => {
  props.validate.push(
    validator.patternMatches(
      swaggerField.pattern,
      'Must be a valid email address',
    ),
  );
  props.type = 'text';

  return props;
};
const configureEdipiField = (swaggerField, props) => {
  props.validate.push(
    validator.patternMatches(swaggerField.pattern, 'Must be a valid DoD ID #'),
  );
  props.type = 'text';

  return props;
};

const renderInputField = ({
  input,
  type,
  step,
  title,
  always_required,
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
      'aria-describedby': input.name + '-error',
    },
    children,
  );

  const displayError = touched && error;
  return (
    <div className={displayError ? 'usa-input-error' : 'usa-input'}>
      <label
        className={displayError ? 'usa-input-error-label' : 'usa-input-label'}
        htmlFor={input.name}
      >
        {title}
        {!always_required && <span className="label-optional">Optional</span>}
      </label>
      {touched &&
        error && (
          <span
            className="usa-input-error-message"
            id={input.name + '-error'}
            role="alert"
          >
            {error}
          </span>
        )}
      {FieldComponent}
    </div>
  );
};

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
        <label htmlFor={fieldName} className="usa-input-label">
          {swaggerField.title || fieldName}
        </label>
      </Fragment>
    );
  }

  // configure the basic Field props
  let fieldProps = {};
  fieldProps.name = nameAttr;
  fieldProps.title = swaggerField.title || fieldName;
  fieldProps.component = renderInputField;
  fieldProps.validate = [];
  fieldProps.always_required = swaggerField[ALWAYS_REQUIRED_KEY];

  let children = null;
  if (swaggerField.enum) {
    fieldProps = configureDropDown(swaggerField, fieldProps);
    children = dropDownChildren(swaggerField);
  } else if (['integer', 'number'].includes(swaggerField.type)) {
    fieldProps = configureNumberField(swaggerField, fieldProps);
  } else if (swaggerField.type === 'string') {
    const fieldFormat = swaggerField.format || swaggerField['x-format'];
    if (fieldFormat === 'date') {
      fieldProps = configureDateField(swaggerField, fieldProps);
    } else if (fieldFormat === 'telephone') {
      fieldProps = configureTelephoneField(swaggerField, fieldProps);
    } else if (fieldFormat === 'ssn') {
      fieldProps = configureSSNField(swaggerField, fieldProps);
    } else if (fieldFormat === 'zip') {
      fieldProps = configureZipField(swaggerField, fieldProps);
    } else if (fieldFormat === 'email') {
      fieldProps = configureEmailField(swaggerField, fieldProps);
    } else if (fieldFormat === 'edipi') {
      fieldProps = configureEdipiField(swaggerField, fieldProps);
      // more cases go here. Datetime, Date,
      // more cases go here. Datetime, Date,
    } else {
      if (swaggerField.pattern) {
        console.error(
          'This swagger field contains a pattern but does not have a custom "format" property',
          fieldName,
          swaggerField,
        );
        console.error(
          "Since it's not feasable to generate a sensible error message from a regex, please add a new format and matching validator",
        );
        fieldProps.validate.push(
          validator.patternMatches(swaggerField.pattern, swaggerField.example),
        );
      }
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
    <Field key={fieldName} {...fieldProps}>
      {children}
    </Field>
  );
};

export default {
  createSchemaField: createSchemaField,
};
