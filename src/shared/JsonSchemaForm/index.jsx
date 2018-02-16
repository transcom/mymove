import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { Field, reduxForm } from 'redux-form';
import './index.css';

const IS_REQUIRED_KEY = 'x-jsf-is-required';

const isEmpty = obj =>
  Object.keys(obj).length === 0 && obj.constructor === Object;
const renderGroupOrField = (fieldName, fields, uiSchema, nameSpace) => {
  /*TODO:
   telephone numbers/ pattern validation
   textbox vs textarea (e.g for addresses)
   dates look wonky in chrome
   styling in accordance with USWDS
   validate group names don't colide with field names
  */
  const group = uiSchema.groups && uiSchema.groups[fieldName];
  const isRef =
    fields[fieldName] &&
    fields[fieldName].$$ref &&
    fields[fieldName].properties;
  if (group) {
    const keys = group.fields;
    return (
      <fieldset key={fieldName}>
        <legend htmlFor={fieldName}>{group.title}</legend>
        {keys.map(f => renderGroupOrField(f, fields, uiSchema, nameSpace))}
      </fieldset>
    );
  } else if (isRef) {
    const refName = fields[fieldName].$$ref.split('/').pop();
    const refSchema = uiSchema.definitions[refName];
    return renderSchema(fields[fieldName], refSchema, fieldName);
  }
  return renderField(fieldName, fields, nameSpace);
};

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

const parseNumberField = value => (!value ? null : Number(value));
const createNumberField = (fieldName, field, nameAttr) => (
  <Field
    component="input"
    name={nameAttr}
    parse={parseNumberField}
    type="Number"
  />
);

const createCheckbox = (fieldName, field, nameAttr) => {
  return (
    <Field id={fieldName} name={nameAttr} component="input" type="checkbox" />
  );
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
    />
  );
};

const requiredValidator = value => (value ? undefined : 'Required');

const createTextAreaField = (fieldName, field, nameAttr) => {
  let validators = [];
  if (field[IS_REQUIRED_KEY]) {
    validators.push(requiredValidator);
  }
  return (
    <Field
      id={fieldName}
      name={nameAttr}
      component="textarea"
      validate={validators}
    />
  );
};

const createInputField = (fieldName, field, nameAttr) => {
  return <Field name={nameAttr} component="input" type={field.format} />;
};

// Ok, so maybe our switches should be building up a single Field component params? right now this is a bit silly.
// keep at it and it will make sense.

// This function switches on the type of the field and creates the correct
// Label and Field combination.
const createField = (fieldName, swaggerField, nameSpace) => {
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
  } else if (swaggerField.type === 'integer') {
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

const renderField = (fieldName, fields, nameSpace) => {
  const field = fields[fieldName];
  if (!field) {
    return;
  }
  return createField(fieldName, field, nameSpace);
};

const renderSchema = (schema, uiSchema, nameSpace = '') => {
  if (schema && !isEmpty(schema)) {
    // Mark all the required fields as required.
    if (schema.required) {
      schema.required.forEach(requiredFieldName => {
        console.log(requiredFieldName);
        schema.properties[requiredFieldName][IS_REQUIRED_KEY] = true;
      });
    }

    const fields = schema.properties || [];
    return uiSchema.order.map(i =>
      renderGroupOrField(i, fields, uiSchema, nameSpace),
    );
  }
};
const JsonSchemaForm = props => {
  const { pristine, submitting, invalid } = props;
  const { handleSubmit, schema, uiSchema } = props;
  const title = schema ? schema.title : '';
  return (
    <form className="default" onSubmit={handleSubmit}>
      <h1>{title}</h1>
      {renderSchema(schema, uiSchema)}
      <button type="submit" disabled={pristine || submitting || invalid}>
        Submit
      </button>
    </form>
  );
};

JsonSchemaForm.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
};

export const reduxifyForm = name => reduxForm({ form: name })(JsonSchemaForm);
