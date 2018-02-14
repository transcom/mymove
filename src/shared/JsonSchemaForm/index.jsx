import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { Field, reduxForm } from 'redux-form';
import './index.css';

const isEmpty = obj =>
  Object.keys(obj).length === 0 && obj.constructor === Object;
const renderGroupOrField = (fieldName, fields, uiSchema) => {
  /*TODO:
   telephone numbers/ pattern validation
   textbox vs textarea (e.g for addresses)
   dates look wonky in chrome
   styling in accordance with USWDS
   validate group names don't colide with field names
  */
  const group = uiSchema.groups[fieldName];
  if (group) {
    const keys = group.fields;
    return (
      <fieldset key={fieldName}>
        <legend htmlFor={fieldName}>{group.title}</legend>
        {keys.map(f => renderGroupOrField(f, fields, uiSchema))}
      </fieldset>
    );
  }
  return renderField(fieldName, fields);
};

const createDropDown = (fieldName, field) => {
  return (
    <Field name={fieldName} component="select">
      <option />
      {field.enum.map((e, index) => (
        <option key={e} value={e}>
          {field['x-display-value'][index]}
        </option>
      ))}
    </Field>
  );
};

const parseNumberField = value => (!value ? null : Number(value));
const createNumberField = (fieldName, field) => (
  <Field
    component="input"
    name={fieldName}
    parse={parseNumberField}
    type="Number"
  />
);

const createCheckbox = (fieldName, field) => {
  return (
    <Field id={fieldName} name={fieldName} component="input" type="checkbox" />
  );
};

const createInputField = (fieldName, field) => {
  return <Field name={fieldName} component="input" type={field.format} />;
};

// This function switches on the type of the field and creates the correct
// Label and Field combination.
const createField = (fieldName, swaggerField) => {
  // Early return here, this is an edge case for label placement.
  // USWDS CSS only renders a checkbox if it is followed by its label
  if (swaggerField.type === 'boolean') {
    return (
      <Fragment key={fieldName}>
        {createCheckbox(fieldName, swaggerField)}
        <label htmlFor={fieldName}>{swaggerField.title || fieldName}</label>
      </Fragment>
    );
  }

  let fieldComponent;
  if (swaggerField.enum) {
    fieldComponent = createDropDown(fieldName, swaggerField);
  } else if (swaggerField.type === 'integer') {
    fieldComponent = createNumberField(fieldName, swaggerField);
  } else {
    // more cases go here. Datetime, Date, UUID
    fieldComponent = createInputField(fieldName, swaggerField);
  }

  return (
    <label key={fieldName}>
      {swaggerField.title || fieldName}
      {fieldComponent}
    </label>
  );
};

const renderField = (fieldName, fields) => {
  const field = fields[fieldName];
  if (!field) {
    return;
  }
  return createField(fieldName, field);
};

const renderSchema = (schema, uiSchema) => {
  if (schema && !isEmpty(schema)) {
    const fields = schema.properties || [];
    return uiSchema.order.map(i => renderGroupOrField(i, fields, uiSchema));
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
