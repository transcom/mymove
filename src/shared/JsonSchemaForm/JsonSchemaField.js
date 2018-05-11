import React, { Fragment } from 'react';

import validator from './validator';
import { Field } from 'redux-form';
import moment from 'moment';
import SingleDatePicker from './SingleDatePicker';
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
  props.componentNameOverride = 'select';

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
      'Number must have 10 digits and a valid area code.',
    ),
  );
  props.type = 'text';

  return props;
};

const configureSSNField = (swaggerField, props) => {
  props.normalize = validator.normalizeSSN;
  props.validate.push(
    validator.patternMatches(
      '^\\d{3}-\\d{2}-\\d{4}$',
      'SSN must have 9 digits.',
    ),
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

const normalizeDates = value => {
  return value ? moment(value).format('YYYY-MM-DD') : value;
};

const configureDateField = (swaggerField, props) => {
  props.type = 'date';
  props.customComponent = SingleDatePicker;
  props.normalize = normalizeDates;
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

const configureEdipiField = (swaggerField, props) => {
  props.validate.push(
    validator.patternMatches(swaggerField.pattern, 'Must be a valid DoD ID #'),
  );
  props.type = 'text';

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

const renderInputField = ({
  input,
  type,
  step,
  title,
  always_required,
  componentNameOverride,
  customComponent,
  meta: { touched, error, warning },
  children,
}) => {
  let component = 'input';
  if (componentNameOverride) {
    component = componentNameOverride;
  }

  if (customComponent) {
    component = customComponent;
  }

  if (componentNameOverride && customComponent) {
    console.error(
      'You should not have specified a componentNameOverride as well as a customComponent. For: ',
      title,
    );
  }

  const FieldComponent = React.createElement(
    component,
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

export const SwaggerField = props => {
  const { fieldName, swagger, required } = props;

  let swaggerField;
  if (swagger.properties) {
    swaggerField = swagger.properties[fieldName];
  }

  if (swaggerField === undefined) {
    return null;
  }

  if (required) {
    swaggerField[ALWAYS_REQUIRED_KEY] = true;
  }

  return createSchemaField(fieldName, swaggerField, undefined);
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

  if (fieldProps.always_required) {
    fieldProps.validate.push(validator.isRequired);
  }

  let children = null;
  if (swaggerField.enum) {
    fieldProps = configureDropDown(swaggerField, fieldProps);
    children = dropDownChildren(swaggerField);
  } else if (['integer', 'number'].includes(swaggerField.type)) {
    fieldProps = configureNumberField(swaggerField, fieldProps);
  } else if (swaggerField.type === 'string') {
    const fieldFormat = swaggerField.format;
    if (fieldFormat === 'date') {
      fieldProps = configureDateField(swaggerField, fieldProps);
    } else if (fieldFormat === 'telephone') {
      fieldProps = configureTelephoneField(swaggerField, fieldProps);
    } else if (fieldFormat === 'ssn') {
      fieldProps = configureSSNField(swaggerField, fieldProps);
    } else if (fieldFormat === 'zip') {
      fieldProps = configureZipField(swaggerField, fieldProps);
    } else if (fieldFormat === 'edipi') {
      fieldProps = configureEdipiField(swaggerField, fieldProps);
    } else if (fieldFormat === 'x-email') {
      fieldProps = configureEmailField(swaggerField, fieldProps);
    } else {
      if (swaggerField.pattern) {
        console.error(
          'This swagger field contains a pattern but does not have a custom "format" property',
          fieldName,
          swaggerField,
        );
        console.error(
          "Since it's not feasible to generate a sensible error message from a regex, please add a new format and matching validator",
        );
        fieldProps.validate.push(
          validator.patternMatches(swaggerField.pattern, swaggerField.example),
        );
      }
      // The last case is the simple text field / textarea which are the same but the componentNameOverride
      if (swaggerField.format === 'textarea') {
        fieldProps.componentNameOverride = 'textarea';
      }
      fieldProps = configureTextField(swaggerField, fieldProps);
    }
  } else {
    console.error(
      'ERROR: This is an unimplemented type in our JSONSchemaForm implementation',
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
