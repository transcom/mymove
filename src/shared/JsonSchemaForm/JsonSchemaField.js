import React, { Fragment } from 'react';
import classNames from 'classnames';
import * as normalizer from './reduxFieldNormalizer';
import validator from './validator';
import { Field } from 'redux-form';
import SingleDatePicker from './SingleDatePicker';
import { isNil } from 'lodash';
export const ALWAYS_REQUIRED_KEY = 'x-always-required';

// ---- Parsers -----

const parseNumberField = (value) => {
  // Empty string will fail Swagger validation, so return null
  if (value === '') {
    return null;
  } else if (!value || validator.isNumber(value)) {
    return value;
  } else {
    return parseFloat(value);
  }
};

// ----- Field configuration -----
const createCheckbox = (fieldName, field, nameAttr, isDisabled) => {
  return (
    <Field
      id={fieldName}
      name={nameAttr}
      component="input"
      type="checkbox"
      className="usa-checkbox__input"
      disabled={isDisabled}
    />
  );
};

const configureDropDown = (swaggerField, props) => {
  props.componentNameOverride = 'select';
  props.inputClassOverride = 'usa-select';

  return props;
};

const dropDownChildren = (swaggerField, filteredEnumListOverride, props) => {
  /* eslint-disable security/detect-object-injection */
  return (
    <Fragment>
      <option />
      {(filteredEnumListOverride ? filteredEnumListOverride : swaggerField.enum).map((e) => (
        <option key={e} value={e}>
          {swaggerField['x-display-value'][e]}
        </option>
      ))}
    </Fragment>
  );
  /* eslint-enable security/detect-object-injection */
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

// TODO: This field should be smarter, it should store int-cents in the redux store
// but allow the user to enter in dollars.
// On first pass, that did not seem straightforward.
const configureCentsField = (swaggerField, props) => {
  // Cents field IS a decimal field
  const decimalLength = 2;
  props = configureDecimalField(swaggerField, props, decimalLength, 'Dollar must be only up to 2 decimal places.');
  props.prefixInputClassName = 'dollar-sign';

  if (swaggerField.maximum != null) {
    props.validate.push(validator.maximum(swaggerField.maximum / 100));
  }
  if (swaggerField.minimum != null) {
    props.validate.push(validator.minimum(swaggerField.minimum / 100));
  }

  return props;
};

// This field allows the form field to accept floats and converts values to
// decimal units for db storage (value * (10 ^ decimalLength))
const configureDecimalField = (swaggerField, props, decimalLength, warningMessage) => {
  props.normalize = normalizer.createDecimalNormalizer(decimalLength);
  props.validate.push(validator.patternMatches(swaggerField.pattern, warningMessage));
  props.validate.push(validator.isNumber);
  props.type = 'text';
  return props;
};

const configureTelephoneField = (swaggerField, props) => {
  props.normalize = normalizer.normalizePhone;
  props.validate.push(
    validator.patternMatches(swaggerField.pattern, 'Number must have 10 digits and a valid area code'),
  );
  props.type = 'text';

  return props;
};

const configureZipField = (swaggerField, props, zipPattern) => {
  props.normalize = normalizer.normalizeZip;
  if (zipPattern) {
    if (zipPattern === 'USA') {
      // RA Summary: eslint - security/detect-unsafe-regex - Denial of Service: Regular Expression
      // RA: Locates potentially unsafe regular expressions, which may take a very long time to run, blocking the event loop
      // RA: Per MilMove SSP, predisposing conditions are regex patterns from untrusted sources or unbounded matching.
      // RA: The regex pattern is a constant string set at compile-time and it is bounded to 10 characters (zip code).
      // RA Developer Status: Mitigated
      // RA Validator Status: Mitigated
      // RA Modified Severity: N/A
      // eslint-disable-next-line security/detect-unsafe-regex
      const zipRegex = '^[0-9]{5}(?:-[0-9]{4})?$';
      props.validate.push(validator.patternMatches(zipRegex, 'Zip code must have 5 or 9 digits.'));
    }
  } else if (swaggerField.pattern) {
    props.validate.push(validator.patternMatches(swaggerField.pattern, 'Zip code must have 5 or 9 digits.'));
  }
  props.type = 'text';

  return props;
};

const configureDateField = (swaggerField, props) => {
  props.type = 'date';
  props.customComponent = SingleDatePicker;
  props.normalize = normalizer.normalizeDates;
  return props;
};

const configureRestrictedDateField = (swaggerField, props, minDate, maxDate) => {
  props.type = 'date';
  props.customComponent = SingleDatePicker;
  props.normalize = normalizer.normalizeDates;
  props.validate.push(validator.minDateValidation(minDate, `Date must be no earlier than ${minDate}`));
  props.validate.push(validator.maxDateValidation(maxDate, `Date must be no later than ${maxDate}`));
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
  props.normalize = normalizer.createDigitNormalizer(swaggerField.maxLength);
  props.validate.push(validator.patternMatches(swaggerField.pattern, 'Must be a valid DoD ID # (10 digits long)'));
  props.type = 'text';

  return props;
};

const configureEmailField = (swaggerField, props) => {
  props.validate.push(validator.patternMatches(swaggerField.pattern, 'Must be a valid email address'));
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
  inputClassOverride,
  customComponent,
  meta: { touched, error, warning },
  children,
  className,
  inputProps,
  hideLabel,
  prefixInputClassName,
}) => {
  let component = 'input';
  if (componentNameOverride) {
    component = componentNameOverride;
  }

  if (customComponent) {
    component = customComponent;
  }

  if (componentNameOverride && customComponent) {
    console.error('You should not have specified a componentNameOverride as well as a customComponent. For: ', title);
  }

  const displayError = touched && error;
  const inputClasses = classNames(
    { 'usa-input--error': displayError },
    { 'usa-input-error-long-message': error && error.length > 57 },
    { [`${inputClassOverride}`]: inputClassOverride },
    { 'usa-input': !inputClassOverride },
  );

  const FieldComponent = React.createElement(
    component,
    {
      ...input,
      type: type,
      step: step,
      'aria-describedby': input.name + '-error',
      className: inputClasses,
      ...inputProps,
    },
    children,
  );

  const classes = classNames('usa-form-group', { 'usa-form-group--error': displayError }, className);

  return (
    <div className={classes}>
      {hideLabel || (
        <label className={classNames('usa-label', { 'usa-label--error': displayError })} htmlFor={input.name}>
          {title}
          {!always_required && type !== 'boolean' && !customComponent && <span className="usa-hint"> (optional)</span>}
        </label>
      )}
      {touched && error && (
        <span className="usa-error-message" id={input.name + '-error'} role="alert">
          {error}
        </span>
      )}
      <span className={prefixInputClassName}>{FieldComponent}</span>
    </div>
  );
};

export const SwaggerField = (props) => {
  const {
    fieldName,
    swagger,
    required,
    className,
    disabled,
    component,
    title,
    onChange,
    validate,
    minDate,
    maxDate,
    disabledDays,
    zipPattern,
    filteredEnumListOverride,
    hideLabel,
  } = props;
  let swaggerField;
  if (swagger.properties) {
    // eslint-disable-next-line security/detect-object-injection
    swaggerField = swagger.properties[fieldName];
  }
  if (swaggerField === undefined) {
    return null;
  }

  if (required) {
    // eslint-disable-next-line security/detect-object-injection
    swaggerField[ALWAYS_REQUIRED_KEY] = true;
  } else {
    // eslint-disable-next-line security/detect-object-injection
    swaggerField[ALWAYS_REQUIRED_KEY] = false;
  }

  return createSchemaField(
    fieldName,
    swaggerField,
    undefined,
    className,
    disabled,
    component,
    title,
    onChange,
    validate,
    minDate,
    maxDate,
    disabledDays,
    zipPattern,
    filteredEnumListOverride,
    hideLabel,
  );
};

// This function switches on the type of the field and creates the correct
// Label and Field combination.
const createSchemaField = (
  fieldName,
  swaggerField,
  nameSpace,
  className = '',
  disabled = false,
  component,
  title,
  onChange,
  validate,
  minDate,
  maxDate,
  disabledDays,
  zipPattern,
  filteredEnumListOverride,
  hideLabel,
) => {
  // Early return here, this is an edge case for label placement.
  // USWDS CSS only renders a checkbox if it is followed by its label
  const nameAttr = nameSpace ? `${nameSpace}.${fieldName}` : fieldName;
  if (swaggerField.type === 'boolean' && !component) {
    return (
      <div key={fieldName} className="usa-checkbox">
        {createCheckbox(fieldName, swaggerField, nameAttr, disabled)}
        {hideLabel || (
          <label htmlFor={fieldName} className="usa-checkbox__label">
            {title || swaggerField.title || fieldName}
          </label>
        )}
      </div>
    );
  }

  // configure the basic Field props
  let fieldProps = {};
  fieldProps.name = nameAttr;
  fieldProps.title = title || swaggerField.title || fieldName;
  fieldProps.onChange = onChange;
  fieldProps.component = renderInputField;
  fieldProps.validate = [];
  // eslint-disable-next-line security/detect-object-injection
  fieldProps.always_required = swaggerField[ALWAYS_REQUIRED_KEY];

  let inputProps = {
    disabled: disabled,
  };

  if (validate) {
    fieldProps.validate.push(validate);
  }

  if (fieldProps.always_required) {
    fieldProps.validate.push(validator.isRequired);
  }

  let children = null;
  if (component) {
    fieldProps.customComponent = component;
  } else if (swaggerField.enum) {
    fieldProps = configureDropDown(swaggerField, fieldProps);
    children = dropDownChildren(swaggerField, filteredEnumListOverride);
  } else if (['integer', 'number'].includes(swaggerField.type)) {
    if (swaggerField.format === 'cents') {
      fieldProps = configureCentsField(swaggerField, fieldProps);
      className += ' dollar-input';
    } else if (swaggerField.format === 'basequantity') {
      fieldProps = configureDecimalField(
        swaggerField,
        fieldProps,
        4,
        'Base quantity must be only up to 4 decimal places.',
      );
    } else if (swaggerField.format === 'dimension') {
      fieldProps.name = nameAttr;
      fieldProps = configureDecimalField(swaggerField, fieldProps, 2, 'Dimension must be only up to 2 decimal places.');
    } else {
      fieldProps = configureNumberField(swaggerField, fieldProps);
    }
  } else if (swaggerField.type === 'string') {
    const fieldFormat = swaggerField.format;
    if (fieldFormat === 'date' && (!isNil(minDate) || !isNil(maxDate))) {
      inputProps.disabledDays = disabledDays ? disabledDays : undefined;
      fieldProps = configureRestrictedDateField(swaggerField, fieldProps, minDate, maxDate);
    } else if (fieldFormat === 'date') {
      fieldProps = configureDateField(swaggerField, fieldProps);
    } else if (fieldFormat === 'telephone') {
      fieldProps = configureTelephoneField(swaggerField, fieldProps);
    } else if (fieldFormat === 'zip') {
      fieldProps = configureZipField(swaggerField, fieldProps, zipPattern);
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
        fieldProps.validate.push(validator.patternMatches(swaggerField.pattern, swaggerField.example));
      }
      // The last case is the simple text field / textarea which are the same but the componentNameOverride
      if (swaggerField.format === 'textarea') {
        fieldProps.componentNameOverride = 'textarea';
      }
      fieldProps = configureTextField(swaggerField, fieldProps);
    }
  } else {
    console.error('ERROR: This is an unimplemented type in our JSONSchemaForm implementation');
  }
  return (
    <Field
      key={fieldName}
      className={className}
      inputProps={inputProps}
      {...fieldProps}
      onChange={onChange}
      hideLabel={hideLabel}
    >
      {children}
    </Field>
  );
};

const schemaFieldExports = {
  createSchemaField: createSchemaField,
};

export default schemaFieldExports;
