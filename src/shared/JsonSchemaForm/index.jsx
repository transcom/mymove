import React from 'react';
import PropTypes from 'prop-types';

import { Field, reduxForm } from 'redux-form';
import './index.css';

const isEmpty = obj =>
  Object.keys(obj).length === 0 && obj.constructor === Object;
const renderGroupOrField = (name, fields, uiSchema) => {
  /*TODO:
   handle enums
   telephone numbers
   textbox vs textarea (e.g for addresses)
   dates look wonky
   styling in accordance with USWDS
   formatting of labels?
   validate group names don't colide with field names
   tests!!!
  */
  const group = uiSchema.groups[name];
  if (group) {
    const keys = group.fields;
    return (
      <fieldset key={name}>
        <legend htmlFor={name}>{group.title}</legend>
        {keys.map(f => renderGroupOrField(f, fields, uiSchema))}
      </fieldset>
    );
  }
  return renderField(name, fields);
};
const renderField = (name, fields) => {
  const field = fields[name];
  if (!field) {
    return;
  }
  return (
    <div key={name}>
      <label htmlFor={name}>
        {field.title || name}
        <Field name={name} component="input" type={field.format} />
      </label>
    </div>
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
