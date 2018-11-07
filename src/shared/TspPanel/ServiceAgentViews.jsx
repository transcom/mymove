import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { PanelSwaggerField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const ServiceAgentDisplay = ({ serviceAgentProps, saRole }) => {
  return (
    <div className="editable-panel-column">
      <span className="column-subhead">{saRole}</span>
      <PanelSwaggerField fieldName="company" {...serviceAgentProps} />
      <PanelSwaggerField fieldName="email" {...serviceAgentProps} />
      <PanelSwaggerField fieldName="phone_number" {...serviceAgentProps} />
    </div>
  );
};

ServiceAgentDisplay.propTypes = {
  serviceAgentProps: PropTypes.shape({
    company: PropTypes.string,
    email: PropTypes.string,
    phone_number: PropTypes.string,
  }),
};

export const ServiceAgentEdit = ({ serviceAgentProps, saRole }) => {
  return (
    <Fragment>
      <div className="editable-panel-column">
        <span className="column-subhead">{saRole}</span>
        <SwaggerField fieldName="company" required {...serviceAgentProps} />
        <SwaggerField fieldName="email" required {...serviceAgentProps} />
        <SwaggerField fieldName="phone_number" required {...serviceAgentProps} />
      </div>
    </Fragment>
  );
};
