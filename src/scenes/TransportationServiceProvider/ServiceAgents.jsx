import { get } from 'lodash';
import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { editablePanelify } from 'shared/EditablePanel';
import { createOrUpdateServiceAgent } from './ducks';

import { PanelSwaggerField } from 'shared/EditablePanel';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const ServiceAgentDisplay = props => {
  let serviceAgent = props.initialValues || {};

  const fieldProps = {
    schema: props.schema,
    values: serviceAgent,
  };
  return (
    <Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField fieldName="company" required {...fieldProps} />
        <PanelSwaggerField fieldName="email" required {...fieldProps} />
        <PanelSwaggerField fieldName="phone_number" required {...fieldProps} />
        <PanelSwaggerField fieldName="email" required {...fieldProps} />
        <PanelSwaggerField fieldName="phone_number" required {...fieldProps} />
      </div>
    </Fragment>
  );
};

const ServiceAgentEdit = props => {
  const schema = props.schema;
  return (
    <Fragment>
      <div className="editable-panel-column">
        <SwaggerField fieldName="company" swagger={schema} required />
        <SwaggerField fieldName="email" swagger={schema} required />
        <SwaggerField fieldName="phone_number" swagger={schema} required />
      </div>
    </Fragment>
  );
};

const formName = 'service_agents_panel';

let ServiceAgentPanel = editablePanelify(ServiceAgentDisplay, ServiceAgentEdit);
ServiceAgentPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(ServiceAgentPanel);

const ServiceAgents = props => {
  const { schema } = props;

  if (!schema) {
    return <LoadingPlaceholder />;
  }

  return (
    <Fragment>
      <ServiceAgentPanel
        form="origin_service_agent"
        title="Origin Service Agent"
        update={props.update}
        schema={props.schema}
        initialValues={props.initialValues.ORIGIN}
        getUpdateArgs={props.getOriginUpdateArgs}
      />

      <ServiceAgentPanel
        form="destination_service_agent"
        title="Destination Service Agent"
        update={props.update}
        schema={props.schema}
        initialValues={props.initialValues.DESTINATION}
        getUpdateArgs={props.getDestinationUpdateArgs}
      />
    </Fragment>
  );
};

function mapStateToProps(state, props) {
  let originFormValues = getFormValues('origin_service_agent')(state);
  let destFormValues = getFormValues('destination_service_agent')(state);
  let serviceAgents = get(state, 'tsp.serviceAgents', []);
  let initialValues = {};
  // This returns the last value. Currently there should only be one record for each role.
  serviceAgents.forEach(agent => (initialValues[agent.role] = agent));

  return {
    // reduxForm
    schema: get(state, 'swagger.spec.definitions.ServiceAgent', null),
    initialValues,

    hasError: false,
    errorMessage: state.tsp.error,
    isUpdating: false,

    // editablePanel
    getOriginUpdateArgs: function() {
      return [
        get(props, 'shipment.id'),
        Object.assign({}, originFormValues, { role: 'ORIGIN' }),
      ];
    },
    getDestinationUpdateArgs: function() {
      return [
        get(props, 'shipment.id'),
        Object.assign({}, destFormValues, { role: 'DESTINATION' }),
      ];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: createOrUpdateServiceAgent,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceAgents);
