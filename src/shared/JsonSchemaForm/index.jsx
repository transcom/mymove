import React from 'react';
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
   tests!!!
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
      {field.enum.map(e => (
        <option key={e} value={e}>
          {e}
        </option>
      ))}
    </Field>
  );
};
const createField = (fieldName, field) => {
  //todo: how to determine if multiselect/checkboxes etc
  if (field.enum) return createDropDown(fieldName, field);

  return <Field name={fieldName} component="input" type={field.format} />;
};

const renderField = (fieldName, fields) => {
  const field = fields[fieldName];
  if (!field) {
    return;
  }
  return (
    <label key={fieldName} htmlFor={fieldName}>
      {field.title || fieldName}
      {createField(fieldName, field)}
    </label>
  );
};

const renderSchema = (schema, uiSchema) => {
  if (!isEmpty(schema)) {
    const fields = schema.properties || [];
    return uiSchema.order.map(i => renderGroupOrField(i, fields, uiSchema));
  }
};
const JsonSchemaForm = props => {
  const { handleSubmit, schema, uiSchema } = props;
  const title = schema ? schema.title : '';
  return (
    <form className="default" onSubmit={handleSubmit}>
      <h1>{title}</h1>
      {renderSchema(schema, uiSchema)}
      <button type="submit">Submit</button>
    </form>
  );
};

JsonSchemaForm.propTypes = {
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  handleSubmit: PropTypes.func.isRequired,
};

export const reduxifyForm = name => reduxForm({ form: name })(JsonSchemaForm);
