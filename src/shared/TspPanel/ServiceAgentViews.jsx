import React, { Fragment } from 'react';
import PropTypes from 'prop-types';

import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export const ServiceAgentDisplay = ({ serviceAgentProps, saRole }) => {
  return (
    <div className="editable-panel-3-column">
      <span className="column-subhead">{saRole} agent</span>
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
      <div className="editable-panel-3-column">
        <span className="column-subhead">{saRole} agent</span>
        <SwaggerField fieldName="company" required {...serviceAgentProps} />
        <SwaggerField fieldName="email" required {...serviceAgentProps} />
        <SwaggerField fieldName="phone_number" required {...serviceAgentProps} />
      </div>
    </Fragment>
  );
};

export const TransportationServiceProviderDisplay = ({ tsp }) => {
  const { name, standard_carrier_alpha_code, poc_general_email, poc_general_phone } = tsp;
  const nameWithScac = name ? `${name} (${standard_carrier_alpha_code})` : '';
  return (
    <div className="editable-panel-3-column">
      <span className="column-subhead">TSP</span>
      <PanelField title="Name" value={nameWithScac} />
      <PanelField title="Email" value={poc_general_email} />
      <PanelField title="Phone number" value={poc_general_phone} />
    </div>
  );
};
