import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { PanelSwaggerField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const ServiceAgentDisplay = ({ serviceAgentProps, role }) => {
  return (
    <div className="editable-panel-column">
      <span className="column-subhead">{role}</span>
      <PanelSwaggerField fieldName="company" required {...serviceAgentProps} />
      <PanelSwaggerField fieldName="email" required {...serviceAgentProps} />
      <PanelSwaggerField fieldName="phone_number" required {...serviceAgentProps} />
    </div>
  );
};

ServiceAgentDisplay.defaultProps = {
  serviceAgentProps: {},
};

ServiceAgentDisplay.propTypes = {
  serviceAgentProps: PropTypes.shape({
    company: PropTypes.string,
    email: PropTypes.string,
    phone_number: PropTypes.string,
  }),
};

export const ServiceAgentEdit = ({ serviceAgentProps, role }) => {
  return (
    <Fragment>
      <div className="editable-panel-column">
        <span className="column-subhead">{role}</span>
        <SwaggerField fieldName="company" {...serviceAgentProps} required />
        <SwaggerField fieldName="email" {...serviceAgentProps} required />
        <SwaggerField fieldName="phone_number" {...serviceAgentProps} required />
      </div>
    </Fragment>
  );
};
